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

type instagramDailyReportRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewInstagramDailyReportRepository(dbDriver infrastructure.DBDriver) port.InstagramDailyReportRepository {
	return &instagramDailyReportRepository{dbDriver: dbDriver}
}

func (r *instagramDailyReportRepository) Find(ctx context.Context, f filter.Filter) (*domain.InstagramDailyReport, error) {
	var m model.InstagramDailyReport
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toInstagramDailyReportDomain(&m), nil
}

func (r *instagramDailyReportRepository) List(ctx context.Context, f filter.Filter) ([]*domain.InstagramDailyReport, error) {
	var models []*model.InstagramDailyReport
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	reports := make([]*domain.InstagramDailyReport, len(models))
	for i, m := range models {
		reports[i] = toInstagramDailyReportDomain(m)
	}
	return reports, nil
}

func (r *instagramDailyReportRepository) Upsert(ctx context.Context, report *domain.InstagramDailyReport) error {
	m, err := toInstagramDailyReportModel(report)
	if err != nil {
		return err
	}
	return r.dbDriver.Upsert(ctx, m,
		[]string{"instagram_connection_id", "report_date"},
		[]string{
			"follower_count", "impressions", "reach", "profile_views",
			"website_clicks", "post_engagement", "audience_demographics",
			"stories_insights", "fetched_at", "updated_at",
		},
	)
}

func toInstagramDailyReportDomain(m *model.InstagramDailyReport) *domain.InstagramDailyReport {
	report := &domain.InstagramDailyReport{
		ID:                    m.ID,
		InstagramConnectionID: m.InstagramConnectionID,
		ReportDate:            m.ReportDate,
		FollowerCount:         m.FollowerCount,
		Impressions:           m.Impressions,
		Reach:                 m.Reach,
		ProfileViews:          m.ProfileViews,
		WebsiteClicks:         m.WebsiteClicks,
		FetchedAt:             m.FetchedAt,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}

	if m.PostEngagement != nil {
		_ = json.Unmarshal(m.PostEngagement, &report.PostEngagement)
	}
	if m.AudienceDemographics != nil {
		_ = json.Unmarshal(m.AudienceDemographics, &report.AudienceDemographics)
	}
	if m.StoriesInsights != nil {
		_ = json.Unmarshal(m.StoriesInsights, &report.StoriesInsights)
	}

	return report
}

func toInstagramDailyReportModel(r *domain.InstagramDailyReport) (*model.InstagramDailyReport, error) {
	postEngagementJSON, err := json.Marshal(r.PostEngagement)
	if err != nil {
		return nil, err
	}
	audienceJSON, err := json.Marshal(r.AudienceDemographics)
	if err != nil {
		return nil, err
	}
	storiesJSON, err := json.Marshal(r.StoriesInsights)
	if err != nil {
		return nil, err
	}

	return &model.InstagramDailyReport{
		ID:                    r.ID,
		InstagramConnectionID: r.InstagramConnectionID,
		ReportDate:            r.ReportDate,
		FollowerCount:         r.FollowerCount,
		Impressions:           r.Impressions,
		Reach:                 r.Reach,
		ProfileViews:          r.ProfileViews,
		WebsiteClicks:         r.WebsiteClicks,
		PostEngagement:        postEngagementJSON,
		AudienceDemographics:  audienceJSON,
		StoriesInsights:       storiesJSON,
		FetchedAt:             r.FetchedAt,
		CreatedAt:             r.CreatedAt,
		UpdatedAt:             r.UpdatedAt,
	}, nil
}
