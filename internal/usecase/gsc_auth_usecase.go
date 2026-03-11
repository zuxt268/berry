package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

// GSCAuthUseCase GSC OAuth連携ユースケースのインターフェース
type GSCAuthUseCase interface {
	// InitiateConnect GSC OAuth連携フローを開始し、認証URLとstateを返す
	InitiateConnect(siteURL string) (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してGSC連携を保存
	HandleCallback(ctx context.Context, userID uint64, code, siteURL string) (*domain.GSCConnection, error)
	// GetConnections ユーザーのGSC連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.GSCConnection, error)
	// Disconnect GSC連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type gscAuthUseCase struct {
	gscOAuthAdapter port.GSCOAuthAdapter
	gscConnRepo     port.GSCConnectionRepository
}

// NewGSCAuthUseCase 新しいGSCAuthUseCaseインスタンスを作成
func NewGSCAuthUseCase(
	gscOAuthAdapter port.GSCOAuthAdapter,
	gscConnRepo port.GSCConnectionRepository,
) GSCAuthUseCase {
	return &gscAuthUseCase{
		gscOAuthAdapter: gscOAuthAdapter,
		gscConnRepo:     gscConnRepo,
	}
}

// InitiateConnect GSC OAuth連携フローを開始
func (u *gscAuthUseCase) InitiateConnect(siteURL string) (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.gscOAuthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してGSC連携を保存
func (u *gscAuthUseCase) HandleCallback(ctx context.Context, userID uint64, code, siteURL string) (*domain.GSCConnection, error) {
	// コードをトークンに交換
	result, err := u.gscOAuthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存の同一サイトの連携を確認
	existing, _ := u.gscConnRepo.Find(ctx, &filter.GSCConnectionFilter{
		UserID:  &userID,
		SiteURL: &siteURL,
	})

	now := time.Now()

	if existing != nil {
		// 既存連携を更新
		existing.RefreshToken = result.RefreshToken
		existing.DisconnectedAt = nil
		existing.ConnectedAt = now
		conn, err := u.gscConnRepo.Update(ctx, existing, &filter.GSCConnectionFilter{ID: &existing.ID})
		if err != nil {
			return nil, err
		}
		slog.Info("GSC connection updated", "userID", userID, "siteURL", siteURL)
		return conn, nil
	}

	// 新規連携を作成
	conn := &domain.GSCConnection{
		UID:          uuid.New().String(),
		UserID:       userID,
		SiteURL:      siteURL,
		RefreshToken: result.RefreshToken,
		ConnectedAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := u.gscConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("GSC connection created", "userID", userID, "siteURL", siteURL)
	return created, nil
}

// GetConnections ユーザーのGSC連携一覧を取得
func (u *gscAuthUseCase) GetConnections(ctx context.Context, userID uint64) ([]*domain.GSCConnection, error) {
	return u.gscConnRepo.List(ctx, &filter.GSCConnectionFilter{UserID: &userID})
}

// Disconnect GSC連携を解除
func (u *gscAuthUseCase) Disconnect(ctx context.Context, userID uint64, uid string) error {
	conn, err := u.gscConnRepo.Find(ctx, &filter.GSCConnectionFilter{UID: &uid, UserID: &userID})
	if err != nil {
		return err
	}

	now := time.Now()
	conn.DisconnectedAt = &now
	conn.RefreshToken = ""

	if _, err := u.gscConnRepo.Update(ctx, conn, &filter.GSCConnectionFilter{ID: &conn.ID}); err != nil {
		return err
	}

	slog.Info("GSC connection disconnected", "userID", userID, "uid", uid)
	return nil
}