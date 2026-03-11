package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type GA4ReportUseCase interface {
	GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GA4DailyReport, error)
}

type ga4ReportUseCase struct {
	connectionRepo port.GA4ConnectionRepository
	reportRepo     port.GA4DailyReportRepository
}

func NewGA4ReportUseCase(
	connectionRepo port.GA4ConnectionRepository,
	reportRepo port.GA4DailyReportRepository,
) GA4ReportUseCase {
	return &ga4ReportUseCase{
		connectionRepo: connectionRepo,
		reportRepo:     reportRepo,
	}
}

func (uc *ga4ReportUseCase) GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GA4DailyReport, error) {
	conn, err := uc.connectionRepo.Find(ctx, &filter.GA4ConnectionFilter{
		UserID:     &userID,
		ActiveOnly: true,
	})
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return uc.reportRepo.List(ctx, &filter.GA4DailyReportFilter{
		GA4ConnectionID: &conn.ID,
		ReportDateFrom:  &from,
		ReportDateTo:    &to,
	})
}