package repository

import (
	"context"
	"encoding/json"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type GBPDailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GBPDailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GBPDailyReport, error)
	Upsert(ctx context.Context, report *domain.GBPDailyReport) error
}

type gbpDailyReportRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGBPDailyReportRepository(dbDriver infrastructure.DBDriver) GBPDailyReportRepository {
	return &gbpDailyReportRepository{dbDriver: dbDriver}
}

func (r *gbpDailyReportRepository) Find(ctx context.Context, f filter.Filter) (*domain.GBPDailyReport, error) {
	var m model.GBPDailyReport
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGBPDailyReportDomain(&m), nil
}

func (r *gbpDailyReportRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GBPDailyReport, error) {
	var models []*model.GBPDailyReport
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	reports := make([]*domain.GBPDailyReport, len(models))
	for i, m := range models {
		reports[i] = toGBPDailyReportDomain(m)
	}
	return reports, nil
}

func (r *gbpDailyReportRepository) Upsert(ctx context.Context, report *domain.GBPDailyReport) error {
	m, err := toGBPDailyReportModel(report)
	if err != nil {
		return err
	}
	return r.dbDriver.Upsert(ctx, m,
		[]string{"gbp_connection_id", "report_date"},
		[]string{
			"profile_views", "phone_calls", "direction_requests",
			"photo_views", "review_count", "average_rating",
			"search_query_breakdown", "fetched_at", "updated_at",
		},
	)
}

func toGBPDailyReportDomain(m *model.GBPDailyReport) *domain.GBPDailyReport {
	report := &domain.GBPDailyReport{
		ID:                m.ID,
		GBPConnectionID:   m.GBPConnectionID,
		ReportDate:        m.ReportDate,
		ProfileViews:      m.ProfileViews,
		PhoneCalls:        m.PhoneCalls,
		DirectionRequests: m.DirectionRequests,
		PhotoViews:        m.PhotoViews,
		ReviewCount:       m.ReviewCount,
		AverageRating:     m.AverageRating,
		FetchedAt:         m.FetchedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	if m.SearchQueryBreakdown != nil {
		_ = json.Unmarshal(m.SearchQueryBreakdown, &report.SearchQueryBreakdown)
	}

	return report
}

func toGBPDailyReportModel(r *domain.GBPDailyReport) (*model.GBPDailyReport, error) {
	searchQueryJSON, err := json.Marshal(r.SearchQueryBreakdown)
	if err != nil {
		return nil, err
	}

	return &model.GBPDailyReport{
		ID:                   r.ID,
		GBPConnectionID:      r.GBPConnectionID,
		ReportDate:           r.ReportDate,
		ProfileViews:         r.ProfileViews,
		PhoneCalls:           r.PhoneCalls,
		DirectionRequests:    r.DirectionRequests,
		PhotoViews:           r.PhotoViews,
		ReviewCount:          r.ReviewCount,
		AverageRating:        r.AverageRating,
		SearchQueryBreakdown: searchQueryJSON,
		FetchedAt:            r.FetchedAt,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}, nil
}
