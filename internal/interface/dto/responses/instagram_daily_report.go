package responses

import "github.com/zuxt268/berry/internal/domain"

type InstagramDailyReportResponse struct {
	ReportDate           string                       `json:"report_date"`
	FollowerCount        int                          `json:"follower_count"`
	Impressions          int                          `json:"impressions"`
	Reach                int                          `json:"reach"`
	ProfileViews         int                          `json:"profile_views"`
	WebsiteClicks        int                          `json:"website_clicks"`
	PostEngagement       []domain.PostEngagement      `json:"post_engagement"`
	AudienceDemographics *domain.AudienceDemographics `json:"audience_demographics,omitempty"`
	StoriesInsights      []domain.StoriesInsight      `json:"stories_insights"`
}

func ToInstagramDailyReportResponse(r *domain.InstagramDailyReport) InstagramDailyReportResponse {
	return InstagramDailyReportResponse{
		ReportDate:           r.ReportDate.Format("2006-01-02"),
		FollowerCount:        r.FollowerCount,
		Impressions:          r.Impressions,
		Reach:                r.Reach,
		ProfileViews:         r.ProfileViews,
		WebsiteClicks:        r.WebsiteClicks,
		PostEngagement:       r.PostEngagement,
		AudienceDemographics: r.AudienceDemographics,
		StoriesInsights:      r.StoriesInsights,
	}
}