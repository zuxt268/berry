package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"github.com/zuxt268/berry/internal/filter"
)

// LineFetchUseCase LINEデータ取得バッチのユースケース
type LineFetchUseCase interface {
	Execute(ctx context.Context, targetDate time.Time) error
}

type lineFetchUseCase struct {
	lineConnRepo        port.LineConnectionRepository
	lineDailyReportRepo port.LineDailyReportRepository
	lineDataAdapter     port.LineDataAdapter
}

func NewLineFetchUseCase(
	lineConnRepo port.LineConnectionRepository,
	lineDailyReportRepo port.LineDailyReportRepository,
	lineDataAdapter port.LineDataAdapter,
) LineFetchUseCase {
	return &lineFetchUseCase{
		lineConnRepo:        lineConnRepo,
		lineDailyReportRepo: lineDailyReportRepo,
		lineDataAdapter:     lineDataAdapter,
	}
}

func (u *lineFetchUseCase) Execute(ctx context.Context, targetDate time.Time) error {
	connections, err := u.lineConnRepo.List(ctx, &filter.LineConnectionFilter{ActiveOnly: true})
	if err != nil {
		return err
	}

	if len(connections) == 0 {
		slog.Info("no active LINE connections found")
		return nil
	}

	slog.Info("starting LINE data fetch", "connections", len(connections), "date", targetDate.Format("2006-01-02"))

	var successCount, failCount int

	for _, conn := range connections {
		if err := ctx.Err(); err != nil {
			slog.Warn("context cancelled, stopping batch", "error", err)
			break
		}

		if err := u.fetchAndSave(ctx, conn, targetDate); err != nil {
			failCount++
			slog.Error("failed to fetch LINE data",
				"connectionID", conn.ID,
				"channelID", conn.ChannelID,
				"userID", conn.UserID,
				"error", err,
			)
			continue
		}
		successCount++
	}

	slog.Info("LINE data fetch completed",
		"processed", len(connections),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func (u *lineFetchUseCase) fetchAndSave(ctx context.Context, conn *domain.LineConnection, targetDate time.Time) error {
	report, err := u.lineDataAdapter.FetchDailyReport(ctx, conn.ChannelAccessToken, targetDate)
	if err != nil {
		return err
	}

	now := time.Now()
	report.LineConnectionID = conn.ID
	report.ReportDate = targetDate
	report.FetchedAt = now
	report.CreatedAt = now
	report.UpdatedAt = now

	return u.lineDailyReportRepo.Upsert(ctx, report)
}