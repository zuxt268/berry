package responses

import "github.com/zuxt268/berry/internal/domain"

type GSCDailyReportResponse struct {
	ReportDate       string                    `json:"report_date"`
	Impressions      int                       `json:"impressions"`
	Clicks           int                       `json:"clicks"`
	CTR              float64                   `json:"ctr"`
	AveragePosition  float64                   `json:"average_position"`
	KeywordBreakdown []domain.KeywordBreakdown `json:"keyword_breakdown"`
	PageBreakdown    []domain.GSCPageBreakdown `json:"page_breakdown"`
}

func ToGSCDailyReportResponse(r *domain.GSCDailyReport) GSCDailyReportResponse {
	return GSCDailyReportResponse{
		ReportDate:       r.ReportDate.Format("2006-01-02"),
		Impressions:      r.Impressions,
		Clicks:           r.Clicks,
		CTR:              r.CTR,
		AveragePosition:  r.AveragePosition,
		KeywordBreakdown: r.KeywordBreakdown,
		PageBreakdown:    r.PageBreakdown,
	}
}