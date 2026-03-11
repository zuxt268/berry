package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"github.com/zuxt268/berry/internal/filter"
)

// GBPFetchUseCase GBPデータ取得バッチのユースケース
type GBPFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type gbpFetchUseCase struct {
	gbpConnRepo        port.GBPConnectionRepository
	gbpDailyReportRepo port.GBPDailyReportRepository
	gbpDataAdapter     port.GBPDataAdapter
}

func NewGBPFetchUseCase(
	gbpConnRepo port.GBPConnectionRepository,
	gbpDailyReportRepo port.GBPDailyReportRepository,
	gbpDataAdapter port.GBPDataAdapter,
) GBPFetchUseCase {
	return &gbpFetchUseCase{
		gbpConnRepo:        gbpConnRepo,
		gbpDailyReportRepo: gbpDailyReportRepo,
		gbpDataAdapter:     gbpDataAdapter,
	}
}

func (u *gbpFetchUseCase) Execute(ctx context.Context, targetDate time.Time) error {
	connections, err := u.gbpConnRepo.List(ctx, &filter.GBPConnectionFilter{ActiveOnly: true})
	if err != nil {
		return err
	}

	if len(connections) == 0 {
		slog.Info("no active GBP connections found")
		return nil
	}

	slog.Info("starting GBP data fetch", "connections", len(connections), "date", targetDate.Format("2006-01-02"))

	var successCount, failCount int

	for _, conn := range connections {
		if err := ctx.Err(); err != nil {
			slog.Warn("context cancelled, stopping batch", "error", err)
			break
		}

		if err := u.fetchAndSave(ctx, conn, targetDate); err != nil {
			failCount++
			slog.Error("failed to fetch GBP data",
				"connectionID", conn.ID,
				"locationID", conn.LocationID,
				"userID", conn.UserID,
				"error", err,
			)
			continue
		}
		successCount++
	}

	slog.Info("GBP data fetch completed",
		"processed", len(connections),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func (u *gbpFetchUseCase) fetchAndSave(ctx context.Context, conn *domain.GBPConnection, targetDate time.Time) error {
	report, err := u.gbpDataAdapter.FetchDailyReport(ctx, conn.RefreshToken, conn.AccountID, conn.LocationID, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report.GBPConnectionID = conn.ID
	report.ReportDate = targetDate
	report.FetchedAt = now
	report.CreatedAt = now
	report.UpdatedAt = now

	return u.gbpDailyReportRepo.Upsert(ctx, report)
}