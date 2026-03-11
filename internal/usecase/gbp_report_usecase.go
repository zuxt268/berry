package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type GBPReportUseCase interface {
	GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GBPDailyReport, error)
}

type gbpReportUseCase struct {
	connectionRepo port.GBPConnectionRepository
	reportRepo     port.GBPDailyReportRepository
}

func NewGBPReportUseCase(
	connectionRepo port.GBPConnectionRepository,
	reportRepo port.GBPDailyReportRepository,
) GBPReportUseCase {
	return &gbpReportUseCase{
		connectionRepo: connectionRepo,
		reportRepo:     reportRepo,
	}
}

func (uc *gbpReportUseCase) GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GBPDailyReport, error) {
	conn, err := uc.connectionRepo.Find(ctx, &filter.GBPConnectionFilter{
		UserID:     &userID,
		ActiveOnly: true,
	})
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return uc.reportRepo.List(ctx, &filter.GBPDailyReportFilter{
		GBPConnectionID: &conn.ID,
		ReportDateFrom:  &from,
		ReportDateTo:    &to,
	})
}