package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"github.com/zuxt268/berry/internal/filter"
)

// GSCFetchUseCase GSCデータ取得バッチのユースケース
type GSCFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type gscFetchUseCase struct {
	gscConnRepo        port.GSCConnectionRepository
	gscDailyReportRepo port.GSCDailyReportRepository
	gscDataAdapter     port.GSCDataAdapter
}

func NewGSCFetchUseCase(
	gscConnRepo port.GSCConnectionRepository,
	gscDailyReportRepo port.GSCDailyReportRepository,
	gscDataAdapter port.GSCDataAdapter,
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
	report, err := u.gscDataAdapter.FetchDailyReport(ctx, conn.RefreshToken, conn.SiteURL, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report.GSCConnectionID = conn.ID
	report.ReportDate = targetDate
	report.FetchedAt = now
	report.CreatedAt = now
	report.UpdatedAt = now

	return u.gscDailyReportRepo.Upsert(ctx, report)
}