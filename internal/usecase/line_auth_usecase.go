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

// LineAuthUseCase LINE公式アカウント連携ユースケースのインターフェース
type LineAuthUseCase interface {
	Connect(ctx context.Context, userID uint64, req *ConnectLineInput) (*domain.LineConnection, error)
	GetConnections(ctx context.Context, userID uint64) ([]*domain.LineConnection, error)
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type lineAuthUseCase struct {
	lineTokenAdapter port.LineTokenAdapter
	lineConnRepo     port.LineConnectionRepository
}

func NewLineAuthUseCase(
	lineTokenAdapter port.LineTokenAdapter,
	lineConnRepo port.LineConnectionRepository,
) LineAuthUseCase {
	return &lineAuthUseCase{
		lineTokenAdapter: lineTokenAdapter,
		lineConnRepo:     lineConnRepo,
	}
}

// Connect トークンを検証してLINE連携を保存
func (u *lineAuthUseCase) Connect(ctx context.Context, userID uint64, req *ConnectLineInput) (*domain.LineConnection, error) {
	// トークン有効性を検証し、Bot情報を取得
	botInfo, err := u.lineTokenAdapter.ValidateToken(ctx, req.ChannelAccessToken)
	if err != nil {
		return nil, err
	}

	channelName := req.ChannelName
	if channelName == "" {
		channelName = botInfo.DisplayName
	}

	// 既存の同一チャンネルの連携を確認
	existing, _ := u.lineConnRepo.Find(ctx, &filter.LineConnectionFilter{
		UserID:    &userID,
		ChannelID: &req.ChannelID,
	})

	now := time.Now()

	if existing != nil {
		// 既存連携を更新
		existing.ChannelSecret = req.ChannelSecret
		existing.ChannelAccessToken = req.ChannelAccessToken
		existing.ChannelName = channelName
		existing.BotUserID = botInfo.UserID
		existing.DisconnectedAt = nil
		existing.ConnectedAt = now
		conn, err := u.lineConnRepo.Update(ctx, existing, &filter.LineConnectionFilter{ID: &existing.ID})
		if err != nil {
			return nil, err
		}
		slog.Info("LINE connection updated", "userID", userID, "channelID", req.ChannelID)
		return conn, nil
	}

	// 新規連携を作成
	conn := &domain.LineConnection{
		UID:                uuid.New().String(),
		UserID:             userID,
		ChannelID:          req.ChannelID,
		ChannelSecret:      req.ChannelSecret,
		ChannelAccessToken: req.ChannelAccessToken,
		ChannelName:        channelName,
		BotUserID:          botInfo.UserID,
		ConnectedAt:        now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	created, err := u.lineConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("LINE connection created", "userID", userID, "channelID", req.ChannelID)
	return created, nil
}

// GetConnections ユーザーのLINE連携一覧を取得
func (u *lineAuthUseCase) GetConnections(ctx context.Context, userID uint64) ([]*domain.LineConnection, error) {
	return u.lineConnRepo.List(ctx, &filter.LineConnectionFilter{UserID: &userID})
}

// Disconnect LINE連携を解除
func (u *lineAuthUseCase) Disconnect(ctx context.Context, userID uint64, uid string) error {
	conn, err := u.lineConnRepo.Find(ctx, &filter.LineConnectionFilter{UID: &uid, UserID: &userID})
	if err != nil {
		return err
	}

	now := time.Now()
	conn.DisconnectedAt = &now
	conn.ChannelAccessToken = ""
	conn.ChannelSecret = ""

	if _, err := u.lineConnRepo.Update(ctx, conn, &filter.LineConnectionFilter{ID: &conn.ID}); err != nil {
		return err
	}

	slog.Info("LINE connection disconnected", "userID", userID, "uid", uid)
	return nil
}