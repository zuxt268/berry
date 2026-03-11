package responses

import "github.com/zuxt268/berry/internal/domain"

type GBPDailyReportResponse struct {
	ReportDate           string                        `json:"report_date"`
	ProfileViews         int                           `json:"profile_views"`
	PhoneCalls           int                           `json:"phone_calls"`
	DirectionRequests    int                           `json:"direction_requests"`
	PhotoViews           int                           `json:"photo_views"`
	ReviewCount          int                           `json:"review_count"`
	AverageRating        float64                       `json:"average_rating"`
	SearchQueryBreakdown []domain.SearchQueryBreakdown `json:"search_query_breakdown"`
}

func ToGBPDailyReportResponse(r *domain.GBPDailyReport) GBPDailyReportResponse {
	return GBPDailyReportResponse{
		ReportDate:           r.ReportDate.Format("2006-01-02"),
		ProfileViews:         r.ProfileViews,
		PhoneCalls:           r.PhoneCalls,
		DirectionRequests:    r.DirectionRequests,
		PhotoViews:           r.PhotoViews,
		ReviewCount:          r.ReviewCount,
		AverageRating:        r.AverageRating,
		SearchQueryBreakdown: r.SearchQueryBreakdown,
	}
}