package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/repository"
)

// GSCFetchUseCase GSCデータ取得バッチのユースケース
type GSCFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type gscFetchUseCase struct {
	gscConnRepo        repository.GSCConnectionRepository
	gscDailyReportRepo repository.GSCDailyReportRepository
	gscDataAdapter     adapter.GSCDataAdapter
}

func NewGSCFetchUseCase(
	gscConnRepo repository.GSCConnectionRepository,
	gscDailyReportRepo repository.GSCDailyReportRepository,
	gscDataAdapter adapter.GSCDataAdapter,
) GSCFetchUseCase {
	return &gscFetchUseCase{
		gscConnRepo:        gscConnRepo,
		gscDailyReportRepo: gscDailyReportRepo,
		gscDataAdapter:     gscDataAdapter,
	}
}

func (u *gscFetchUseCase) Execute(ctx context.Context, targetDate time.Time) error {
	connections, err := u.gscConnRepo.List(ctx, &filter.GSCConnectionFilter{ActiveOnly: true})
	if err != nil {
		return err
	}

	if len(connections) == 0 {
		slog.Info("no active GSC connections found")
		return nil
	}

	slog.Info("starting GSC data fetch", "connections", len(connections), "date", targetDate.Format("2006-01-02"))

	var successCount, failCount int

	for _, conn := range connections {
		if err := ctx.Err(); err != nil {
			slog.Warn("context cancelled, stopping batch", "error", err)
			break
		}

		if err := u.fetchAndSave(ctx, conn, targetDate); err != nil {
			failCount++
			slog.Error("failed to fetch GSC data",
				"connectionID", conn.ID,
				"siteURL", conn.SiteURL,
				"userID", conn.UserID,
				"error", err,
			)
			continue
		}
		successCount++
	}

	slog.Info("GSC data fetch completed",
		"processed", len(connections),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func (u *gscFetchUseCase) fetchAndSave(ctx context.Context, conn *domain.GSCConnection, targetDate time.Time) error {
	data, err := u.gscDataAdapter.FetchDailyReport(ctx, conn.RefreshToken, conn.SiteURL, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report := &domain.GSCDailyReport{
		GSCConnectionID:  conn.ID,
		ReportDate:       targetDate,
		Impressions:      data.Impressions,
		Clicks:           data.Clicks,
		CTR:              data.CTR,
		AveragePosition:  data.AveragePosition,
		KeywordBreakdown: data.KeywordBreakdown,
		PageBreakdown:    data.PageBreakdown,
		FetchedAt:        now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return u.gscDailyReportRepo.Upsert(ctx, report)
}
