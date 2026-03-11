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

// InstagramAuthUseCase Instagram OAuth連携ユースケースのインターフェース
type InstagramAuthUseCase interface {
	// InitiateConnect Instagram OAuth連携フローを開始し、認証URLとstateを返す
	InitiateConnect() (authURL string, state string, err error)
	// HandleCallback OAuthコールバックを処理してInstagram連携を保存
	HandleCallback(ctx context.Context, userID uint64, code string) (*domain.InstagramConnection, error)
	// GetConnections ユーザーのInstagram連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.InstagramConnection, error)
	// Disconnect Instagram連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type instagramAuthUseCase struct {
	instagramOAuthAdapter port.InstagramOAuthAdapter
	instagramConnRepo     port.InstagramConnectionRepository
}

func NewInstagramAuthUseCase(
	instagramOAuthAdapter port.InstagramOAuthAdapter,
	instagramConnRepo port.InstagramConnectionRepository,
) InstagramAuthUseCase {
	return &instagramAuthUseCase{
		instagramOAuthAdapter: instagramOAuthAdapter,
		instagramConnRepo:     instagramConnRepo,
	}
}

// InitiateConnect Instagram OAuth連携フローを開始
func (u *instagramAuthUseCase) InitiateConnect() (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	url := u.instagramOAuthAdapter.GetAuthURL(state)
	return url, state, nil
}

// HandleCallback OAuthコールバックを処理してInstagram連携を保存
func (u *instagramAuthUseCase) HandleCallback(ctx context.Context, userID uint64, code string) (*domain.InstagramConnection, error) {
	// コードをトークン+IGアカウント情報に交換
	result, err := u.instagramOAuthAdapter.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 既存の同一IGアカウントの連携を確認
	existing, _ := u.instagramConnRepo.Find(ctx, &filter.InstagramConnectionFilter{
		UserID:                     &userID,
		InstagramBusinessAccountID: &result.InstagramBusinessAccountID,
	})

	now := time.Now()

	if existing != nil {
		// 既存連携を更新
		existing.AccessToken = result.AccessToken
		existing.TokenExpiresAt = &result.TokenExpiresAt
		existing.FacebookPageID = result.FacebookPageID
		existing.DisconnectedAt = nil
		existing.ConnectedAt = now
		conn, err := u.instagramConnRepo.Update(ctx, existing, &filter.InstagramConnectionFilter{ID: &existing.ID})
		if err != nil {
			return nil, err
		}
		slog.Info("Instagram connection updated", "userID", userID, "igAccountID", result.InstagramBusinessAccountID)
		return conn, nil
	}

	// 新規連携を作成
	conn := &domain.InstagramConnection{
		UID:                        uuid.New().String(),
		UserID:                     userID,
		InstagramBusinessAccountID: result.InstagramBusinessAccountID,
		FacebookPageID:             result.FacebookPageID,
		AccessToken:                result.AccessToken,
		TokenExpiresAt:             &result.TokenExpiresAt,
		ConnectedAt:                now,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	created, err := u.instagramConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("Instagram connection created", "userID", userID, "igAccountID", result.InstagramBusinessAccountID)
	return created, nil
}

// GetConnections ユーザーのInstagram連携一覧を取得
func (u *instagramAuthUseCase) GetConnections(ctx context.Context, userID uint64) ([]*domain.InstagramConnection, error) {
	return u.instagramConnRepo.List(ctx, &filter.InstagramConnectionFilter{UserID: &userID})
}

// Disconnect Instagram連携を解除
func (u *instagramAuthUseCase) Disconnect(ctx context.Context, userID uint64, uid string) error {
	conn, err := u.instagramConnRepo.Find(ctx, &filter.InstagramConnectionFilter{UID: &uid, UserID: &userID})
	if err != nil {
		return err
	}

	now := time.Now()
	conn.DisconnectedAt = &now
	conn.AccessToken = ""

	if _, err := u.instagramConnRepo.Update(ctx, conn, &filter.InstagramConnectionFilter{ID: &conn.ID}); err != nil {
		return err
	}

	slog.Info("Instagram connection disconnected", "userID", userID, "uid", uid)
	return nil
}