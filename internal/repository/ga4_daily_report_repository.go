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

type ga4DailyReportRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGA4DailyReportRepository(dbDriver infrastructure.DBDriver) port.GA4DailyReportRepository {
	return &ga4DailyReportRepository{dbDriver: dbDriver}
}

func (r *ga4DailyReportRepository) Find(ctx context.Context, f filter.Filter) (*domain.GA4DailyReport, error) {
	var m model.GA4DailyReport
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGA4DailyReportDomain(&m), nil
}

func (r *ga4DailyReportRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GA4DailyReport, error) {
	var models []*model.GA4DailyReport
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	reports := make([]*domain.GA4DailyReport, len(models))
	for i, m := range models {
		reports[i] = toGA4DailyReportDomain(m)
	}
	return reports, nil
}

func (r *ga4DailyReportRepository) Upsert(ctx context.Context, report *domain.GA4DailyReport) error {
	m, err := toGA4DailyReportModel(report)
	if err != nil {
		return err
	}
	return r.dbDriver.Upsert(ctx, m,
		[]string{"ga4_connection_id", "report_date"},
		[]string{
			"sessions", "total_users", "bounce_rate", "avg_session_duration",
			"conversions", "channel_breakdown", "device_breakdown",
			"page_breakdown", "fetched_at", "updated_at",
		},
	)
}

func toGA4DailyReportDomain(m *model.GA4DailyReport) *domain.GA4DailyReport {
	report := &domain.GA4DailyReport{
		ID:                 m.ID,
		GA4ConnectionID:    m.GA4ConnectionID,
		ReportDate:         m.ReportDate,
		Sessions:           m.Sessions,
		TotalUsers:         m.TotalUsers,
		BounceRate:         m.BounceRate,
		AvgSessionDuration: m.AvgSessionDuration,
		Conversions:        m.Conversions,
		FetchedAt:          m.FetchedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}

	if m.ChannelBreakdown != nil {
		_ = json.Unmarshal(m.ChannelBreakdown, &report.ChannelBreakdown)
	}
	if m.DeviceBreakdown != nil {
		_ = json.Unmarshal(m.DeviceBreakdown, &report.DeviceBreakdown)
	}
	if m.PageBreakdown != nil {
		_ = json.Unmarshal(m.PageBreakdown, &report.PageBreakdown)
	}

	return report
}

func toGA4DailyReportModel(r *domain.GA4DailyReport) (*model.GA4DailyReport, error) {
	channelJSON, err := json.Marshal(r.ChannelBreakdown)
	if err != nil {
		return nil, err
	}
	deviceJSON, err := json.Marshal(r.DeviceBreakdown)
	if err != nil {
		return nil, err
	}
	pageJSON, err := json.Marshal(r.PageBreakdown)
	if err != nil {
		return nil, err
	}

	return &model.GA4DailyReport{
		ID:                 r.ID,
		GA4ConnectionID:    r.GA4ConnectionID,
		ReportDate:         r.ReportDate,
		Sessions:           r.Sessions,
		TotalUsers:         r.TotalUsers,
		BounceRate:         r.BounceRate,
		AvgSessionDuration: r.AvgSessionDuration,
		Conversions:        r.Conversions,
		ChannelBreakdown:   channelJSON,
		DeviceBreakdown:    deviceJSON,
		PageBreakdown:      pageJSON,
		FetchedAt:          r.FetchedAt,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}, nil
}
