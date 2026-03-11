package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type LineReportUseCase interface {
	GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.LineDailyReport, error)
}

type lineReportUseCase struct {
	connectionRepo port.LineConnectionRepository
	reportRepo     port.LineDailyReportRepository
}

func NewLineReportUseCase(
	connectionRepo port.LineConnectionRepository,
	reportRepo port.LineDailyReportRepository,
) LineReportUseCase {
	return &lineReportUseCase{
		connectionRepo: connectionRepo,
		reportRepo:     reportRepo,
	}
}

func (uc *lineReportUseCase) GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.LineDailyReport, error) {
	conn, err := uc.connectionRepo.Find(ctx, &filter.LineConnectionFilter{
		UserID:     &userID,
		ActiveOnly: true,
	})
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return uc.reportRepo.List(ctx, &filter.LineDailyReportFilter{
		LineConnectionID: &conn.ID,
		ReportDateFrom:   &from,
		ReportDateTo:     &to,
	})
}