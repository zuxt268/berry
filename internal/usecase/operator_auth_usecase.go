package usecase

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/repository"
)

// OperatorAuthUseCase オペレーター認証ユースケースのインターフェース
type OperatorAuthUseCase interface {
	// InitiateLogin stateを生成してOAuth URLを返す
	InitiateLogin(w http.ResponseWriter) (string, error)
	// HandleCallback OAuthコールバックを処理してセッションを作成
	HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, code, state string) (*domain.Operator, error)
	// GetCurrentOperator 現在認証済みのオペレーターを取得
	GetCurrentOperator(ctx context.Context, r *http.Request) (*domain.Operator, bool, error)
	// Logout セッションを削除してオペレーターをログアウト
	Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) error
	// VerifyState OAuthのstateパラメータを検証
	VerifyState(r *http.Request, state string) error
}

type operatorAuthUseCase struct {
	oauthAdapter        adapter.OAuthAdapter
	sessionAdapter      adapter.SessionAdapter
	operatorRepo        repository.OperatorRepository
	operatorSessionRepo repository.OperatorSessionRepository
}

// NewOperatorAuthUseCase 新しいOperatorAuthUseCaseインスタンスを作成
func NewOperatorAuthUseCase(
	oauthAdapter adapter.OAuthAdapter,
	sessionAdapter adapter.SessionAdapter,
	operatorRepo repository.OperatorRepository,
	operatorSessionRepo repository.OperatorSessionRepository,
) OperatorAuthUseCase {
	return &operatorAuthUseCase{
		oauthAdapter:        oauthAdapter,
		sessionAdapter:      sessionAdapter,
		operatorRepo:        operatorRepo,
		operatorSessionRepo: operatorSessionRepo,
	}
}

// InitiateLogin stateを生成してOAuth URLを返す
func (u *operatorAuthUseCase) InitiateLogin(w http.ResponseWriter) (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     "operator_oauthstate",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	url := u.oauthAdapter.GetAuthURL(state)
	return url, nil
}

// HandleCallback OAuthコールバックを処理してセッションを作成
func (u *operatorAuthUseCase) HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, code, state string) (*domain.Operator, error) {
	if err := u.VerifyState(r, state); err != nil {
		return nil, err
	}

	oauthResult, err := u.oauthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存のオペレーターをEmailで検索（自動作成しない）
	operator, err := u.operatorRepo.Find(ctx, &filter.OperatorFilter{Email: &oauthResult.Email})
	if err != nil {
		slog.Warn("operator not found for login", "email", oauthResult.Email)
		return nil, err
	}

	if !operator.IsActive {
		slog.Warn("inactive operator attempted login", "uid", operator.UID, "email", operator.Email)
		return nil, err
	}

	if operator.Name == "" {
		operator.Name = oauthResult.Name
	}

	if _, err := u.operatorRepo.Update(ctx, operator, &filter.OperatorFilter{ID: &operator.ID}); err != nil {
		return nil, err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// 既存セッション削除
	if err := u.operatorSessionRepo.Delete(ctx, &filter.OperatorSessionFilter{OperatorID: &operator.ID}); err != nil {
		return nil, err
	}

	session := &domain.OperatorSession{
		OperatorID:   operator.ID,
		SessionToken: sessionToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	if err := u.operatorSessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	if err := u.sessionAdapter.SaveSessionToken(r, w, sessionToken); err != nil {
		return nil, err
	}

	slog.Info("operator login successful", "uid", operator.UID, "email", operator.Email)
	return operator, nil
}

// GetCurrentOperator 現在認証済みのオペレーターを取得
func (u *operatorAuthUseCase) GetCurrentOperator(ctx context.Context, r *http.Request) (*domain.Operator, bool, error) {
	token, ok, err := u.sessionAdapter.GetSessionToken(r)
	if err != nil || !ok {
		return nil, false, err
	}

	session, err := u.operatorSessionRepo.Get(ctx, &filter.OperatorSessionFilter{SessionToken: &token})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if session.IsExpired() {
		return nil, false, nil
	}

	operator, err := u.operatorRepo.Find(ctx, &filter.OperatorFilter{ID: &session.OperatorID})
	if err != nil {
		return nil, false, err
	}

	return operator, true, nil
}

// Logout セッションを削除してオペレーターをログアウト
func (u *operatorAuthUseCase) Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	token, ok, err := u.sessionAdapter.GetSessionToken(r)
	if err != nil || !ok {
		return err
	}

	session, err := u.operatorSessionRepo.Get(ctx, &filter.OperatorSessionFilter{SessionToken: &token})
	if err == nil && session != nil {
		_ = u.operatorSessionRepo.Delete(ctx, &filter.OperatorSessionFilter{ID: &session.ID})
		slog.Info("operator logout", "operatorID", session.OperatorID)
	}

	return u.sessionAdapter.DeleteSessionToken(r, w)
}

// VerifyState OAuthのstateパラメータを検証
func (u *operatorAuthUseCase) VerifyState(r *http.Request, state string) error {
	cookie, err := r.Cookie("operator_oauthstate")
	if err != nil {
		return err
	}
	if state != cookie.Value {
		return http.ErrNoCookie
	}
	return nil
}
