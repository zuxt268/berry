package usecase

import (
	"context"
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

// GBPAuthUseCase GBP OAuth連携ユースケースのインターフェース
type GBPAuthUseCase interface {
	// InitiateConnect GBP OAuth連携フローを開始
	InitiateConnect(r *http.Request, w http.ResponseWriter, locationID, accountID string) (string, error)
	// HandleCallback OAuthコールバックを処理してGBP連携を保存
	HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.GBPConnection, error)
	// GetConnections ユーザーのGBP連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.GBPConnection, error)
	// Disconnect GBP連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type gbpAuthUseCase struct {
	gbpOAuthAdapter adapter.GBPOAuthAdapter
	gbpConnRepo     repository.GBPConnectionRepository
}

// NewGBPAuthUseCase 新しいGBPAuthUseCaseインスタンスを作成
func NewGBPAuthUseCase(
	gbpOAuthAdapter adapter.GBPOAuthAdapter,
	gbpConnRepo repository.GBPConnectionRepository,
) GBPAuthUseCase {
	return &gbpAuthUseCase{
		gbpOAuthAdapter: gbpOAuthAdapter,
		gbpConnRepo:     gbpConnRepo,
	}
}

// InitiateConnect GBP OAuth連携フローを開始
func (u *gbpAuthUseCase) InitiateConnect(r *http.Request, w http.ResponseWriter, locationID, accountID string) (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	// stateをクッキーに保存
	stateCookie := &http.Cookie{
		Name:     "gbp_oauthstate",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, stateCookie)

	// location_idをクッキーに保存
	locationCookie := &http.Cookie{
		Name:     "gbp_location_id",
		Value:    locationID,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, locationCookie)

	// account_idをクッキーに保存
	accountCookie := &http.Cookie{
		Name:     "gbp_account_id",
		Value:    accountID,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, accountCookie)

	url := u.gbpOAuthAdapter.GetAuthURL(state)
	return url, nil
}

// HandleCallback OAuthコールバックを処理してGBP連携を保存
func (u *gbpAuthUseCase) HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.GBPConnection, error) {
	// state検証
	stateCookie, err := r.Cookie("gbp_oauthstate")
	if err != nil || state != stateCookie.Value {
		return nil, domain.ErrInvalidToken
	}

	// location_id取得
	locationCookie, err := r.Cookie("gbp_location_id")
	if err != nil || locationCookie.Value == "" {
		return nil, domain.ErrInvalidArgument
	}
	locationID := locationCookie.Value

	// account_id取得
	accountCookie, err := r.Cookie("gbp_account_id")
	if err != nil || accountCookie.Value == "" {
		return nil, domain.ErrInvalidArgument
	}
	accountID := accountCookie.Value

	// クッキー削除
	clearCookie(w, "gbp_oauthstate")
	clearCookie(w, "gbp_location_id")
	clearCookie(w, "gbp_account_id")

	// コードをトークンに交換
	result, err := u.gbpOAuthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存の同一ロケーションの連携を確認
	existing, _ := u.gbpConnRepo.Find(ctx, &filter.GBPConnectionFilter{
		UserID:     &userID,
		LocationID: &locationID,
	})

	now := time.Now()

	if existing != nil {
		// 既存連携を更新
		existing.RefreshToken = result.RefreshToken
		existing.DisconnectedAt = nil
		existing.ConnectedAt = now
		conn, err := u.gbpConnRepo.Update(ctx, existing, &filter.GBPConnectionFilter{ID: &existing.ID})
		if err != nil {
			return nil, err
		}
		slog.Info("GBP connection updated", "userID", userID, "locationID", locationID)
		return conn, nil
	}

	// 新規連携を作成
	conn := &domain.GBPConnection{
		UID:          uuid.New().String(),
		UserID:       userID,
		LocationID:   locationID,
		AccountID:    accountID,
		RefreshToken: result.RefreshToken,
		ConnectedAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := u.gbpConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("GBP connection created", "userID", userID, "locationID", locationID)
	return created, nil
}

// GetConnections ユーザーのGBP連携一覧を取得
func (u *gbpAuthUseCase) GetConnections(ctx context.Context, userID uint64) ([]*domain.GBPConnection, error) {
	return u.gbpConnRepo.List(ctx, &filter.GBPConnectionFilter{UserID: &userID})
}

// Disconnect GBP連携を解除
func (u *gbpAuthUseCase) Disconnect(ctx context.Context, userID uint64, uid string) error {
	conn, err := u.gbpConnRepo.Find(ctx, &filter.GBPConnectionFilter{UID: &uid, UserID: &userID})
	if err != nil {
		return err
	}

	now := time.Now()
	conn.DisconnectedAt = &now
	conn.RefreshToken = ""

	if _, err := u.gbpConnRepo.Update(ctx, conn, &filter.GBPConnectionFilter{ID: &conn.ID}); err != nil {
		return err
	}

	slog.Info("GBP connection disconnected", "userID", userID, "uid", uid)
	return nil
}
