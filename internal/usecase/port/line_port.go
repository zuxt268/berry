package port

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

// LineBotInfo Bot Info APIのレスポンス
type LineBotInfo struct {
	UserID      string
	DisplayName string
	BasicID     string
}

// LineTokenAdapter LINEチャンネルアクセストークンの検証
type LineTokenAdapter interface {
	ValidateToken(ctx context.Context, channelAccessToken string) (*LineBotInfo, error)
}

// LineDataAdapter LINE Messaging APIから統計データを取得するアダプター
type LineDataAdapter interface {
	FetchDailyReport(ctx context.Context, channelAccessToken string, date time.Time) (*domain.LineDailyReport, error)
}