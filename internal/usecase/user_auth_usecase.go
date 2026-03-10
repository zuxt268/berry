package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/repository"
)

// UserAuthUseCase 認証ユースケースのインターフェース
type UserAuthUseCase interface {
	// InitiateLogin stateを生成してOAuth URLを返す
	InitiateLogin(w http.ResponseWriter) (string, error)
	// HandleCallback OAuthコールバックを処理してセッションを作成
	HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, code, state string) (*domain.User, error)
	// GetCurrentUser 現在認証済みのユーザーを取得
	GetCurrentUser(ctx context.Context, r *http.Request) (*domain.User, bool, error)
	// Logout セッションを削除してユーザーをログアウト
	Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) error
	// VerifyState OAuthのstateパラメータを検証
	VerifyState(r *http.Request, state string) error
}

type userAuthUseCase struct {
	oauthAdapter          adapter.OAuthAdapter
	sessionAdapter        adapter.SessionAdapter
	userRepository        repository.UserRepository
	userSessionRepository repository.UserSessionRepository
}

// NewAuthUseCase 新しいAuthUseCaseインスタンスを作成
func NewAuthUseCase(
	oauthAdapter adapter.OAuthAdapter,
	sessionAdapter adapter.SessionAdapter,
	userRepo repository.UserRepository,
	userSessionRepo repository.UserSessionRepository,
) UserAuthUseCase {
	return &userAuthUseCase{
		oauthAdapter:          oauthAdapter,
		sessionAdapter:        sessionAdapter,
		userRepository:        userRepo,
		userSessionRepository: userSessionRepo,
	}
}

// InitiateLogin stateを生成してOAuth URLを返す
func (u *userAuthUseCase) InitiateLogin(w http.ResponseWriter) (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	setStateCookie(w, state)

	url := u.oauthAdapter.GetAuthURL(state)
	return url, nil
}

// HandleCallback OAuthコールバックを処理してセッションを作成
func (u *userAuthUseCase) HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, code, state string) (*domain.User, error) {
	if err := u.VerifyState(r, state); err != nil {
		return nil, err
	}

	oauthResult, err := u.oauthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存ユーザーをEmailで検索（自動作成しない）
	user, err := u.userRepository.Find(ctx, &filter.UserFilter{Email: &oauthResult.Email})
	if err != nil {
		slog.Warn("customer user not found for login", "email", oauthResult.Email)
		return nil, err
	}

	if !user.IsActive() {
		slog.Warn("inactive customer user attempted login", "uid", user.UID, "email", user.Email)
		return nil, err
	}

	// 名前を更新、GoogleSubが未設定なら保存
	user.Name = oauthResult.Name

	if _, err := u.userRepository.Update(ctx, user, &filter.UserFilter{
		ID: &user.ID}); err != nil {
		return nil, err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	if err := u.userSessionRepository.Delete(ctx, &filter.UserSessionFilter{
		UserID: &user.ID,
	}); err != nil {
		return nil, err
	}

	session := &domain.UserSession{
		UID:          uuid.NewString(),
		UserID:       user.ID,
		SessionToken: sessionToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := u.userSessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}

	if err := u.sessionAdapter.SaveSessionToken(r, w, sessionToken); err != nil {
		return nil, err
	}

	slog.Info("customer login successful", "uid", user.UID, "email", user.Email)
	return user, nil
}

// GetCurrentUser 現在認証済みのユーザーを取得
func (u *userAuthUseCase) GetCurrentUser(ctx context.Context, r *http.Request) (*domain.User, bool, error) {
	token, ok, err := u.sessionAdapter.GetSessionToken(r)
	if err != nil || !ok {
		return nil, false, err
	}

	session, err := u.userSessionRepository.Get(ctx, &filter.UserSessionFilter{SessionToken: &token})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if session.IsExpired() {
		return nil, false, nil
	}

	customerUser, err := u.userRepository.Find(ctx, &filter.UserFilter{ID: &session.UserID})
	if err != nil {
		return nil, false, err
	}

	return customerUser, true, nil
}

// Logout セッションを削除してユーザーをログアウト
func (u *userAuthUseCase) Logout(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	token, ok, err := u.sessionAdapter.GetSessionToken(r)
	if err != nil || !ok {
		return err
	}

	session, err := u.userSessionRepository.Get(ctx, &filter.UserSessionFilter{SessionToken: &token})
	if err == nil && session != nil {
		_ = u.userSessionRepository.Delete(ctx, &filter.UserSessionFilter{ID: &session.ID})
		slog.Info("customer logout", "userID", session.UserID)
	}

	return u.sessionAdapter.DeleteSessionToken(r, w)
}

// VerifyState OAuthのstateパラメータを検証
func (u *userAuthUseCase) VerifyState(r *http.Request, state string) error {
	cookie, err := r.Cookie("oauthstate")
	if err != nil {
		return err
	}
	if state != cookie.Value {
		return http.ErrNoCookie
	}
	return nil
}

// ヘルパー関数

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func setStateCookie(w http.ResponseWriter, state string) {
	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}
