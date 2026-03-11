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

type lineDailyReportRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewLineDailyReportRepository(dbDriver infrastructure.DBDriver) port.LineDailyReportRepository {
	return &lineDailyReportRepository{dbDriver: dbDriver}
}

func (r *lineDailyReportRepository) Find(ctx context.Context, f filter.Filter) (*domain.LineDailyReport, error) {
	var m model.LineDailyReport
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toLineDailyReportDomain(&m), nil
}

func (r *lineDailyReportRepository) List(ctx context.Context, f filter.Filter) ([]*domain.LineDailyReport, error) {
	var models []*model.LineDailyReport
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	reports := make([]*domain.LineDailyReport, len(models))
	for i, m := range models {
		reports[i] = toLineDailyReportDomain(m)
	}
	return reports, nil
}

func (r *lineDailyReportRepository) Upsert(ctx context.Context, report *domain.LineDailyReport) error {
	m, err := toLineDailyReportModel(report)
	if err != nil {
		return err
	}
	return r.dbDriver.Upsert(ctx, m,
		[]string{"line_connection_id", "report_date"},
		[]string{
			"followers", "targeted_reaches", "blocks",
			"message_delivery", "demographic",
			"fetched_at", "updated_at",
		},
	)
}

func toLineDailyReportDomain(m *model.LineDailyReport) *domain.LineDailyReport {
	report := &domain.LineDailyReport{
		ID:               m.ID,
		LineConnectionID: m.LineConnectionID,
		ReportDate:       m.ReportDate,
		Followers:        m.Followers,
		TargetedReaches:  m.TargetedReaches,
		Blocks:           m.Blocks,
		FetchedAt:        m.FetchedAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	if m.MessageDelivery != nil {
		_ = json.Unmarshal(m.MessageDelivery, &report.MessageDelivery)
	}
	if m.Demographic != nil {
		_ = json.Unmarshal(m.Demographic, &report.Demographic)
	}

	return report
}

func toLineDailyReportModel(r *domain.LineDailyReport) (*model.LineDailyReport, error) {
	messageDeliveryJSON, err := json.Marshal(r.MessageDelivery)
	if err != nil {
		return nil, err
	}
	demographicJSON, err := json.Marshal(r.Demographic)
	if err != nil {
		return nil, err
	}

	return &model.LineDailyReport{
		ID:               r.ID,
		LineConnectionID: r.LineConnectionID,
		ReportDate:       r.ReportDate,
		Followers:        r.Followers,
		TargetedReaches:  r.TargetedReaches,
		Blocks:           r.Blocks,
		MessageDelivery:  messageDeliveryJSON,
		Demographic:      demographicJSON,
		FetchedAt:        r.FetchedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}, nil
}
