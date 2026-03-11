package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type InstagramReportUseCase interface {
	GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.InstagramDailyReport, error)
}

type instagramReportUseCase struct {
	connectionRepo port.InstagramConnectionRepository
	reportRepo     port.InstagramDailyReportRepository
}

func NewInstagramReportUseCase(
	connectionRepo port.InstagramConnectionRepository,
	reportRepo port.InstagramDailyReportRepository,
) InstagramReportUseCase {
	return &instagramReportUseCase{
		connectionRepo: connectionRepo,
		reportRepo:     reportRepo,
	}
}

func (uc *instagramReportUseCase) GetReports(ctx context.Context, userID uint64, from, to time.Time) ([]*domain.InstagramDailyReport, error) {
	conn, err := uc.connectionRepo.Find(ctx, &filter.InstagramConnectionFilter{
		UserID:     &userID,
		ActiveOnly: true,
	})
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return uc.reportRepo.List(ctx, &filter.InstagramDailyReportFilter{
		InstagramConnectionID: &conn.ID,
		ReportDateFrom:        &from,
		ReportDateTo:          &to,
	})
}