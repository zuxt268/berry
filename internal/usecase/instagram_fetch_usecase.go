package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"github.com/zuxt268/berry/internal/filter"
)

// InstagramFetchUseCase Instagramデータ取得バッチのユースケース
type InstagramFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type instagramFetchUseCase struct {
	instagramConnRepo        port.InstagramConnectionRepository
	instagramDailyReportRepo port.InstagramDailyReportRepository
	instagramDataAdapter     port.InstagramDataAdapter
}

func NewInstagramFetchUseCase(
	instagramConnRepo port.InstagramConnectionRepository,
	instagramDailyReportRepo port.InstagramDailyReportRepository,
	instagramDataAdapter port.InstagramDataAdapter,
) InstagramFetchUseCase {
	return &instagramFetchUseCase{
		instagramConnRepo:        instagramConnRepo,
		instagramDailyReportRepo: instagramDailyReportRepo,
		instagramDataAdapter:     instagramDataAdapter,
	}
}

func (u *instagramFetchUseCase) Execute(ctx context.Context, targetDate time.Time) error {
	connections, err := u.instagramConnRepo.List(ctx, &filter.InstagramConnectionFilter{ActiveOnly: true})
	if err != nil {
		return err
	}

	if len(connections) == 0 {
		slog.Info("no active Instagram connections found")
		return nil
	}

	slog.Info("starting Instagram data fetch", "connections", len(connections), "date", targetDate.Format("2006-01-02"))

	var successCount, failCount int

	for _, conn := range connections {
		if err := ctx.Err(); err != nil {
			slog.Warn("context cancelled, stopping batch", "error", err)
			break
		}

		if err := u.fetchAndSave(ctx, conn, targetDate); err != nil {
			failCount++
			slog.Error("failed to fetch Instagram data",
				"connectionID", conn.ID,
				"igAccountID", conn.InstagramBusinessAccountID,
				"userID", conn.UserID,
				"error", err,
			)
			continue
		}
		successCount++
	}

	slog.Info("Instagram data fetch completed",
		"processed", len(connections),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func (u *instagramFetchUseCase) fetchAndSave(ctx context.Context, conn *domain.InstagramConnection, targetDate time.Time) error {
	report, err := u.instagramDataAdapter.FetchDailyReport(ctx, conn.AccessToken, conn.InstagramBusinessAccountID, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report.InstagramConnectionID = conn.ID
	report.ReportDate = targetDate
	report.FetchedAt = now
	report.CreatedAt = now
	report.UpdatedAt = now

	return u.instagramDailyReportRepo.Upsert(ctx, report)
}