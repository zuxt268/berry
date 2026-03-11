package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

// OperatorAuthUseCase オペレーター認証ユースケースのインターフェース
type OperatorAuthUseCase interface {
	// InitiateLogin stateを生成してOAuth URLとstateを返す
	InitiateLogin() (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してセッションを作成し、オペレーターとセッショントークンを返す
	HandleCallback(ctx context.Context, code, ipAddress, userAgent string) (*domain.Operator, string, error)
	// GetCurrentOperator セッショントークンから現在認証済みのオペレーターを取得
	GetCurrentOperator(ctx context.Context, sessionToken string) (*domain.Operator, bool, error)
	// Logout セッショントークンに対応するセッションを削除
	Logout(ctx context.Context, sessionToken string) error
}

type operatorAuthUseCase struct {
	oauthAdapter        port.OAuthAdapter
	operatorRepo        port.OperatorRepository
	operatorSessionRepo port.OperatorSessionRepository
}

// NewOperatorAuthUseCase 新しいOperatorAuthUseCaseインスタンスを作成
func NewOperatorAuthUseCase(
	oauthAdapter port.OAuthAdapter,
	operatorRepo port.OperatorRepository,
	operatorSessionRepo port.OperatorSessionRepository,
) OperatorAuthUseCase {
	return &operatorAuthUseCase{
		oauthAdapter:        oauthAdapter,
		operatorRepo:        operatorRepo,
		operatorSessionRepo: operatorSessionRepo,
	}
}

// InitiateLogin stateを生成してOAuth URLとstateを返す
func (u *operatorAuthUseCase) InitiateLogin() (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.oauthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してセッションを作成
func (u *operatorAuthUseCase) HandleCallback(ctx context.Context, code, ipAddress, userAgent string) (*domain.Operator, string, error) {
	oauthResult, err := u.oauthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", err
	}

	// 既存のオペレーターをEmailで検索（自動作成しない）
	operator, err := u.operatorRepo.Find(ctx, &filter.OperatorFilter{Email: &oauthResult.Email})
	if err != nil {
		slog.Warn("operator not found for login", "email", oauthResult.Email)
		return nil, "", err
	}

	if !operator.IsActive {
		slog.Warn("inactive operator attempted login", "uid", operator.UID, "email", operator.Email)
		return nil, "", domain.ErrForbidden
	}

	if operator.Name == "" {
		operator.Name = oauthResult.Name
	}

	if _, err := u.operatorRepo.Update(ctx, operator, &filter.OperatorFilter{ID: &operator.ID}); err != nil {
		return nil, "", err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", err
	}

	// 既存セッション削除
	if err := u.operatorSessionRepo.Delete(ctx, &filter.OperatorSessionFilter{OperatorID: &operator.ID}); err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	slog.Info("operator login successful", "uid", operator.UID, "email", operator.Email)
	return operator, sessionToken, nil
}

// GetCurrentOperator セッショントークンから現在認証済みのオペレーターを取得
func (u *operatorAuthUseCase) GetCurrentOperator(ctx context.Context, sessionToken string) (*domain.Operator, bool, error) {
	if sessionToken == "" {
		return nil, false, nil
	}

	session, err := u.operatorSessionRepo.Get(ctx, &filter.OperatorSessionFilter{SessionToken: &sessionToken})
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

// Logout セッショントークンに対応するセッションを削除
func (u *operatorAuthUseCase) Logout(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	session, err := u.operatorSessionRepo.Get(ctx, &filter.OperatorSessionFilter{SessionToken: &sessionToken})
	if err == nil && session != nil {
		_ = u.operatorSessionRepo.Delete(ctx, &filter.OperatorSessionFilter{ID: &session.ID})
		slog.Info("operator logout", "operatorID", session.OperatorID)
	}

	return nil
}