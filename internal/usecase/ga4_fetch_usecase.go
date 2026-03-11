package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

// GA4FetchUseCase GA4データ取得バッチのユースケース
type GA4FetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type ga4FetchUseCase struct {
	ga4ConnRepo        port.GA4ConnectionRepository
	ga4DailyReportRepo port.GA4DailyReportRepository
	ga4DataAdapter     port.GA4DataAdapter
}

func NewGA4FetchUseCase(
	ga4ConnRepo port.GA4ConnectionRepository,
	ga4DailyReportRepo port.GA4DailyReportRepository,
	ga4DataAdapter port.GA4DataAdapter,
) GA4FetchUseCase {
	return &ga4FetchUseCase{
		ga4ConnRepo:        ga4ConnRepo,
		ga4DailyReportRepo: ga4DailyReportRepo,
		ga4DataAdapter:     ga4DataAdapter,
	}
}

func (u *ga4FetchUseCase) Execute(ctx context.Context, targetDate time.Time) error {
	connections, err := u.ga4ConnRepo.List(ctx, &filter.GA4ConnectionFilter{ActiveOnly: true})
	if err != nil {
		return err
	}

	if len(connections) == 0 {
		slog.Info("no active GA4 connections found")
		return nil
	}

	slog.Info("starting GA4 data fetch", "connections", len(connections), "date", targetDate.Format("2006-01-02"))

	var successCount, failCount int

	for _, conn := range connections {
		if err := ctx.Err(); err != nil {
			slog.Warn("context cancelled, stopping batch", "error", err)
			break
		}

		if err := u.fetchAndSave(ctx, conn, targetDate); err != nil {
			failCount++
			slog.Error("failed to fetch GA4 data",
				"connectionID", conn.ID,
				"propertyID", conn.GooglePropertyID,
				"userID", conn.UserID,
				"error", err,
			)
			continue
		}
		successCount++
	}

	slog.Info("GA4 data fetch completed",
		"processed", len(connections),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func (u *ga4FetchUseCase) fetchAndSave(ctx context.Context, conn *domain.GA4Connection, targetDate time.Time) error {
	report, err := u.ga4DataAdapter.FetchDailyReport(ctx, conn.RefreshToken, conn.GooglePropertyID, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report.GA4ConnectionID = conn.ID
	report.ReportDate = targetDate
	report.FetchedAt = now
	report.CreatedAt = now
	report.UpdatedAt = now

	return u.ga4DailyReportRepo.Upsert(ctx, report)
}
