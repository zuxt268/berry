package port

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

// GSCOAuthResult GSC OAuth認証で取得したトークン情報
type GSCOAuthResult struct {
	RefreshToken string
	AccessToken  string
}

// GSCOAuthAdapter GSC OAuth操作のインターフェース
type GSCOAuthAdapter interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*GSCOAuthResult, error)
}

// GSCDataAdapter Search Console APIからレポートデータを取得するアダプター
type GSCDataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, siteURL string, date time.Time) (*domain.GSCDailyReport, error)
}