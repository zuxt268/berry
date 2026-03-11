package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

// UserAuthUseCase 認証ユースケースのインターフェース
type UserAuthUseCase interface {
	// InitiateLogin stateを生成してOAuth URLとstateを返す
	InitiateLogin() (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してセッションを作成し、ユーザーとセッショントークンを返す
	HandleCallback(ctx context.Context, code, ipAddress, userAgent string) (*domain.User, string, error)
	// GetCurrentUser セッショントークンから現在認証済みのユーザーを取得
	GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, bool, error)
	// Logout セッショントークンに対応するセッションを削除
	Logout(ctx context.Context, sessionToken string) error
}

type userAuthUseCase struct {
	oauthAdapter          port.OAuthAdapter
	userRepository        port.UserRepository
	userSessionRepository port.UserSessionRepository
}

// NewAuthUseCase 新しいAuthUseCaseインスタンスを作成
func NewAuthUseCase(
	oauthAdapter port.OAuthAdapter,
	userRepo port.UserRepository,
	userSessionRepo port.UserSessionRepository,
) UserAuthUseCase {
	return &userAuthUseCase{
		oauthAdapter:          oauthAdapter,
		userRepository:        userRepo,
		userSessionRepository: userSessionRepo,
	}
}

// InitiateLogin stateを生成してOAuth URLとstateを返す
func (u *userAuthUseCase) InitiateLogin() (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.oauthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してセッションを作成
func (u *userAuthUseCase) HandleCallback(ctx context.Context, code, ipAddress, userAgent string) (*domain.User, string, error) {
	oauthResult, err := u.oauthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", err
	}

	// 既存ユーザーをEmailで検索（自動作成しない）
	user, err := u.userRepository.Find(ctx, &filter.UserFilter{Email: &oauthResult.Email})
	if err != nil {
		slog.Warn("customer user not found for login", "email", oauthResult.Email)
		return nil, "", err
	}

	if !user.IsActive() {
		slog.Warn("inactive customer user attempted login", "uid", user.UID, "email", user.Email)
		return nil, "", domain.ErrForbidden
	}

	// 名前を更新
	user.Name = oauthResult.Name

	if _, err := u.userRepository.Update(ctx, user, &filter.UserFilter{
		ID: &user.ID}); err != nil {
		return nil, "", err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", err
	}

	if err := u.userSessionRepository.Delete(ctx, &filter.UserSessionFilter{
		UserID: &user.ID,
	}); err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	slog.Info("customer login successful", "uid", user.UID, "email", user.Email)
	return user, sessionToken, nil
}

// GetCurrentUser セッショントークンから現在認証済みのユーザーを取得
func (u *userAuthUseCase) GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, bool, error) {
	if sessionToken == "" {
		return nil, false, nil
	}

	session, err := u.userSessionRepository.Get(ctx, &filter.UserSessionFilter{SessionToken: &sessionToken})
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

// Logout セッショントークンに対応するセッションを削除
func (u *userAuthUseCase) Logout(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	session, err := u.userSessionRepository.Get(ctx, &filter.UserSessionFilter{SessionToken: &sessionToken})
	if err == nil && session != nil {
		_ = u.userSessionRepository.Delete(ctx, &filter.UserSessionFilter{ID: &session.ID})
		slog.Info("customer logout", "userID", session.UserID)
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