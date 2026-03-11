package repository

import (
	"context"
	"encoding/json"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/repository/model"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type gscDailyReportRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGSCDailyReportRepository(dbDriver infrastructure.DBDriver) port.GSCDailyReportRepository {
	return &gscDailyReportRepository{dbDriver: dbDriver}
}

func (r *gscDailyReportRepository) Find(ctx context.Context, f filter.Filter) (*domain.GSCDailyReport, error) {
	var m model.GSCDailyReport
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGSCDailyReportDomain(&m), nil
}

func (r *gscDailyReportRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GSCDailyReport, error) {
	var models []*model.GSCDailyReport
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	reports := make([]*domain.GSCDailyReport, len(models))
	for i, m := range models {
		reports[i] = toGSCDailyReportDomain(m)
	}
	return reports, nil
}

func (r *gscDailyReportRepository) Upsert(ctx context.Context, report *domain.GSCDailyReport) error {
	m, err := toGSCDailyReportModel(report)
	if err != nil {
		return err
	}
	return r.dbDriver.Upsert(ctx, m,
		[]string{"gsc_connection_id", "report_date"},
		[]string{
			"impressions", "clicks", "ctr", "average_position",
			"keyword_breakdown", "page_breakdown",
			"fetched_at", "updated_at",
		},
	)
}

func toGSCDailyReportDomain(m *model.GSCDailyReport) *domain.GSCDailyReport {
	report := &domain.GSCDailyReport{
		ID:              m.ID,
		GSCConnectionID: m.GSCConnectionID,
		ReportDate:      m.ReportDate,
		Impressions:     m.Impressions,
		Clicks:          m.Clicks,
		CTR:             m.CTR,
		AveragePosition: m.AveragePosition,
		FetchedAt:       m.FetchedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	if m.KeywordBreakdown != nil {
		_ = json.Unmarshal(m.KeywordBreakdown, &report.KeywordBreakdown)
	}
	if m.PageBreakdown != nil {
		_ = json.Unmarshal(m.PageBreakdown, &report.PageBreakdown)
	}

	return report
}

func toGSCDailyReportModel(r *domain.GSCDailyReport) (*model.GSCDailyReport, error) {
	keywordJSON, err := json.Marshal(r.KeywordBreakdown)
	if err != nil {
		return nil, err
	}
	pageJSON, err := json.Marshal(r.PageBreakdown)
	if err != nil {
		return nil, err
	}

	return &model.GSCDailyReport{
		ID:               r.ID,
		GSCConnectionID:  r.GSCConnectionID,
		ReportDate:       r.ReportDate,
		Impressions:      r.Impressions,
		Clicks:           r.Clicks,
		CTR:              r.CTR,
		AveragePosition:  r.AveragePosition,
		KeywordBreakdown: keywordJSON,
		PageBreakdown:    pageJSON,
		FetchedAt:        r.FetchedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}, nil
}
