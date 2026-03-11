package port

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

// InstagramOAuthResult Instagram OAuth認証で取得したトークン情報
type InstagramOAuthResult struct {
	AccessToken                string
	TokenExpiresAt             time.Time
	InstagramBusinessAccountID string
	FacebookPageID             string
}

// InstagramOAuthAdapter Instagram OAuth操作のインターフェース
type InstagramOAuthAdapter interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*InstagramOAuthResult, error)
	RefreshLongLivedToken(ctx context.Context, currentToken string) (string, *time.Time, error)
}

// InstagramDataAdapter Instagram Graph APIからレポートデータを取得するアダプター
type InstagramDataAdapter interface {
	FetchDailyReport(ctx context.Context, accessToken string, igAccountID string, date time.Time) (*domain.InstagramDailyReport, error)
}