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

// InstagramAuthUseCase Instagram OAuth連携ユースケースのインターフェース
type InstagramAuthUseCase interface {
	InitiateConnect(r *http.Request, w http.ResponseWriter) (string, error)
	HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.InstagramConnection, error)
	GetConnections(ctx context.Context, userID uint64) ([]*domain.InstagramConnection, error)
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type instagramAuthUseCase struct {
	instagramOAuthAdapter adapter.InstagramOAuthAdapter
	instagramConnRepo     repository.InstagramConnectionRepository
}

func NewInstagramAuthUseCase(
	instagramOAuthAdapter adapter.InstagramOAuthAdapter,
	instagramConnRepo repository.InstagramConnectionRepository,
) InstagramAuthUseCase {
	return &instagramAuthUseCase{
		instagramOAuthAdapter: instagramOAuthAdapter,
		instagramConnRepo:     instagramConnRepo,
	}
}

// InitiateConnect Instagram OAuth連携フローを開始
func (u *instagramAuthUseCase) InitiateConnect(r *http.Request, w http.ResponseWriter) (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	stateCookie := &http.Cookie{
		Name:     "instagram_oauthstate",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, stateCookie)

	url := u.instagramOAuthAdapter.GetAuthURL(state)
	return url, nil
}

// HandleCallback OAuthコールバックを処理してInstagram連携を保存
func (u *instagramAuthUseCase) HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.InstagramConnection, error) {
	// state検証
	stateCookie, err := r.Cookie("instagram_oauthstate")
	if err != nil || state != stateCookie.Value {
		return nil, domain.ErrInvalidToken
	}

	// クッキー削除
	clearCookie(w, "instagram_oauthstate")

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
