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

// GBPFetchUseCase GBPデータ取得バッチのユースケース
type GBPFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type gbpFetchUseCase struct {
	gbpConnRepo        repository.GBPConnectionRepository
	gbpDailyReportRepo repository.GBPDailyReportRepository
	gbpDataAdapter     adapter.GBPDataAdapter
}

func NewGBPFetchUseCase(
	gbpConnRepo repository.GBPConnectionRepository,
	gbpDailyReportRepo repository.GBPDailyReportRepository,
	gbpDataAdapter adapter.GBPDataAdapter,
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
	data, err := u.gbpDataAdapter.FetchDailyReport(ctx, conn.RefreshToken, conn.AccountID, conn.LocationID, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report := &domain.GBPDailyReport{
		GBPConnectionID:      conn.ID,
		ReportDate:           targetDate,
		ProfileViews:         data.ProfileViews,
		PhoneCalls:           data.PhoneCalls,
		DirectionRequests:    data.DirectionRequests,
		PhotoViews:           data.PhotoViews,
		ReviewCount:          data.ReviewCount,
		AverageRating:        data.AverageRating,
		SearchQueryBreakdown: data.SearchQueryBreakdown,
		FetchedAt:            now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	return u.gbpDailyReportRepo.Upsert(ctx, report)
}
