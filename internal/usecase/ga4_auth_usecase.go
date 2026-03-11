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

// GA4AuthUseCase GA4 OAuth連携ユースケースのインターフェース
type GA4AuthUseCase interface {
	// InitiateConnect GA4 OAuth連携フローを開始し、認証URLとstateを返す
	InitiateConnect(propertyID string) (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してGA4連携を保存
	HandleCallback(ctx context.Context, userID uint64, code, propertyID string) (*domain.GA4Connection, error)
	// GetConnections ユーザーのGA4連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.GA4Connection, error)
	// Disconnect GA4連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type ga4AuthUseCase struct {
	ga4OAuthAdapter port.GA4OAuthAdapter
	ga4ConnRepo     port.GA4ConnectionRepository
}

// NewGA4AuthUseCase 新しいGA4AuthUseCaseインスタンスを作成
func NewGA4AuthUseCase(
	ga4OAuthAdapter port.GA4OAuthAdapter,
	ga4ConnRepo port.GA4ConnectionRepository,
) GA4AuthUseCase {
	return &ga4AuthUseCase{
		ga4OAuthAdapter: ga4OAuthAdapter,
		ga4ConnRepo:     ga4ConnRepo,
	}
}

// InitiateConnect GA4 OAuth連携フローを開始
func (u *ga4AuthUseCase) InitiateConnect(propertyID string) (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.ga4OAuthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してGA4連携を保存
func (u *ga4AuthUseCase) HandleCallback(ctx context.Context, userID uint64, code, propertyID string) (*domain.GA4Connection, error) {
	// コードをトークンに交換
	result, err := u.ga4OAuthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存の同一プロパティの連携を確認
	existing, _ := u.ga4ConnRepo.Find(ctx, &filter.GA4ConnectionFilter{
		UserID:           &userID,
		GooglePropertyID: &propertyID,
	})

	now := time.Now()

	if existing != nil {
		// 既存連携を更新
		existing.RefreshToken = result.RefreshToken
		existing.DisconnectedAt = nil
		existing.ConnectedAt = now
		conn, err := u.ga4ConnRepo.Update(ctx, existing, &filter.GA4ConnectionFilter{ID: &existing.ID})
		if err != nil {
			return nil, err
		}
		slog.Info("GA4 connection updated", "userID", userID, "propertyID", propertyID)
		return conn, nil
	}

	// 新規連携を作成
	conn := &domain.GA4Connection{
		UID:              uuid.New().String(),
		UserID:           userID,
		GooglePropertyID: propertyID,
		RefreshToken:     result.RefreshToken,
		ConnectedAt:      now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	created, err := u.ga4ConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("GA4 connection created", "userID", userID, "propertyID", propertyID)
	return created, nil
}

// GetConnections ユーザーのGA4連携一覧を取得
func (u *ga4AuthUseCase) GetConnections(ctx context.Context, userID uint64) ([]*domain.GA4Connection, error) {
	return u.ga4ConnRepo.List(ctx, &filter.GA4ConnectionFilter{UserID: &userID})
}

// Disconnect GA4連携を解除
func (u *ga4AuthUseCase) Disconnect(ctx context.Context, userID uint64, uid string) error {
	conn, err := u.ga4ConnRepo.Find(ctx, &filter.GA4ConnectionFilter{UID: &uid, UserID: &userID})
	if err != nil {
		return err
	}

	now := time.Now()
	conn.DisconnectedAt = &now
	conn.RefreshToken = ""

	if _, err := u.ga4ConnRepo.Update(ctx, conn, &filter.GA4ConnectionFilter{ID: &conn.ID}); err != nil {
		return err
	}

	slog.Info("GA4 connection disconnected", "userID", userID, "uid", uid)
	return nil
}