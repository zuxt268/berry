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

// GBPAuthUseCase GBP OAuth連携ユースケースのインターフェース
type GBPAuthUseCase interface {
	// InitiateConnect GBP OAuth連携フローを開始し、認証URLとstateを返す
	InitiateConnect(locationID, accountID string) (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してGBP連携を保存
	HandleCallback(ctx context.Context, userID uint64, code, locationID, accountID string) (*domain.GBPConnection, error)
	// GetConnections ユーザーのGBP連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.GBPConnection, error)
	// Disconnect GBP連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type gbpAuthUseCase struct {
	gbpOAuthAdapter port.GBPOAuthAdapter
	gbpConnRepo     port.GBPConnectionRepository
}

// NewGBPAuthUseCase 新しいGBPAuthUseCaseインスタンスを作成
func NewGBPAuthUseCase(
	gbpOAuthAdapter port.GBPOAuthAdapter,
	gbpConnRepo port.GBPConnectionRepository,
) GBPAuthUseCase {
	return &gbpAuthUseCase{
		gbpOAuthAdapter: gbpOAuthAdapter,
		gbpConnRepo:     gbpConnRepo,
	}
}

// InitiateConnect GBP OAuth連携フローを開始
func (u *gbpAuthUseCase) InitiateConnect(locationID, accountID string) (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.gbpOAuthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してGBP連携を保存
func (u *gbpAuthUseCase) HandleCallback(ctx context.Context, userID uint64, code, locationID, accountID string) (*domain.GBPConnection, error) {
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