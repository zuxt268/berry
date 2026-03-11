package port

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

// GA4OAuthResult GA4 OAuth認証で取得したトークン情報
type GA4OAuthResult struct {
	RefreshToken string
	AccessToken  string
}

// GA4OAuthAdapter GA4 OAuth操作のインターフェース
type GA4OAuthAdapter interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*GA4OAuthResult, error)
}

// GA4DataAdapter GA4 Data APIからレポートデータを取得するアダプター
type GA4DataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, propertyID string, date time.Time) (*domain.GA4DailyReport, error)
}