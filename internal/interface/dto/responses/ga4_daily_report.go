package responses

import "github.com/zuxt268/berry/internal/domain"

type GA4DailyReportResponse struct {
	ReportDate         string                    `json:"report_date"`
	Sessions           int                       `json:"sessions"`
	TotalUsers         int                       `json:"total_users"`
	BounceRate         float64                   `json:"bounce_rate"`
	AvgSessionDuration float64                   `json:"avg_session_duration"`
	Conversions        int                       `json:"conversions"`
	ChannelBreakdown   []domain.ChannelBreakdown `json:"channel_breakdown"`
	DeviceBreakdown    []domain.DeviceBreakdown  `json:"device_breakdown"`
	PageBreakdown      []domain.PageBreakdown    `json:"page_breakdown"`
}

func ToGA4DailyReportResponse(r *domain.GA4DailyReport) GA4DailyReportResponse {
	return GA4DailyReportResponse{
		ReportDate:         r.ReportDate.Format("2006-01-02"),
		Sessions:           r.Sessions,
		TotalUsers:         r.TotalUsers,
		BounceRate:         r.BounceRate,
		AvgSessionDuration: r.AvgSessionDuration,
		Conversions:        r.Conversions,
		ChannelBreakdown:   r.ChannelBreakdown,
		DeviceBreakdown:    r.DeviceBreakdown,
		PageBreakdown:      r.PageBreakdown,
	}
}