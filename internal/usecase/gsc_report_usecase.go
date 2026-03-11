package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type GSCReportUseCase interface {
	GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GSCDailyReport, error)
}

type gscReportUseCase struct {
	connectionRepo port.GSCConnectionRepository
	reportRepo     port.GSCDailyReportRepository
}

func NewGSCReportUseCase(
	connectionRepo port.GSCConnectionRepository,
	reportRepo port.GSCDailyReportRepository,
) GSCReportUseCase {
	return &gscReportUseCase{
		connectionRepo: connectionRepo,
		reportRepo:     reportRepo,
	}
}

func (uc *gscReportUseCase) GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.GSCDailyReport, error) {
	conn, err := uc.connectionRepo.Find(ctx, &filter.GSCConnectionFilter{
		UserID:     &userID,
		ActiveOnly: true,
	})
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return uc.reportRepo.List(ctx, &filter.GSCDailyReportFilter{
		GSCConnectionID: &conn.ID,
		ReportDateFrom:  &from,
		ReportDateTo:    &to,
	})
}