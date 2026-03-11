package port

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

// GBPOAuthResult GBP OAuth認証で取得したトークン情報
type GBPOAuthResult struct {
	RefreshToken string
	AccessToken  string
}

// GBPOAuthAdapter GBP OAuth操作のインターフェース
type GBPOAuthAdapter interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*GBPOAuthResult, error)
}

// GBPDataAdapter Google Business Profile APIからレポートデータを取得するアダプター
type GBPDataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, accountID string, locationID string, date time.Time) (*domain.GBPDailyReport, error)
}