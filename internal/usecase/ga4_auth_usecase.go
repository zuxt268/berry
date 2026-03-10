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

// GA4AuthUseCase GA4 OAuth連携ユースケースのインターフェース
type GA4AuthUseCase interface {
	// InitiateConnect GA4 OAuth連携フローを開始
	InitiateConnect(r *http.Request, w http.ResponseWriter, propertyID string) (string, error)
	// HandleCallback OAuthコールバックを処理してGA4連携を保存
	HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.GA4Connection, error)
	// GetConnections ユーザーのGA4連携一覧を取得
	GetConnections(ctx context.Context, userID uint64) ([]*domain.GA4Connection, error)
	// Disconnect GA4連携を解除
	Disconnect(ctx context.Context, userID uint64, uid string) error
}

type ga4AuthUseCase struct {
	ga4OAuthAdapter adapter.GA4OAuthAdapter
	ga4ConnRepo     repository.GA4ConnectionRepository
}

// NewGA4AuthUseCase 新しいGA4AuthUseCaseインスタンスを作成
func NewGA4AuthUseCase(
	ga4OAuthAdapter adapter.GA4OAuthAdapter,
	ga4ConnRepo repository.GA4ConnectionRepository,
) GA4AuthUseCase {
	return &ga4AuthUseCase{
		ga4OAuthAdapter: ga4OAuthAdapter,
		ga4ConnRepo:     ga4ConnRepo,
	}
}

// InitiateConnect GA4 OAuth連携フローを開始
func (u *ga4AuthUseCase) InitiateConnect(r *http.Request, w http.ResponseWriter, propertyID string) (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	// stateをクッキーに保存
	stateCookie := &http.Cookie{
		Name:     "ga4_oauthstate",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, stateCookie)

	// property_idをクッキーに保存
	propertyCookie := &http.Cookie{
		Name:     "ga4_property_id",
		Value:    propertyID,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, propertyCookie)

	url := u.ga4OAuthAdapter.GetAuthURL(state)
	return url, nil
}

// HandleCallback OAuthコールバックを処理してGA4連携を保存
func (u *ga4AuthUseCase) HandleCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, userID uint64, code, state string) (*domain.GA4Connection, error) {
	// state検証
	stateCookie, err := r.Cookie("ga4_oauthstate")
	if err != nil || state != stateCookie.Value {
		return nil, domain.ErrInvalidToken
	}

	// property_id取得
	propertyCookie, err := r.Cookie("ga4_property_id")
	if err != nil || propertyCookie.Value == "" {
		return nil, domain.ErrInvalidArgument
	}
	propertyID := propertyCookie.Value

	// クッキー削除
	clearCookie(w, "ga4_oauthstate")
	clearCookie(w, "ga4_property_id")

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

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
}
