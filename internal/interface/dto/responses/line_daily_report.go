package responses

import "github.com/zuxt268/berry/internal/domain"

type LineDailyReportResponse struct {
	ReportDate      string                      `json:"report_date"`
	Followers       int                         `json:"followers"`
	TargetedReaches int                         `json:"targeted_reaches"`
	Blocks          int                         `json:"blocks"`
	MessageDelivery *domain.LineMessageDelivery `json:"message_delivery,omitempty"`
	Demographic     *domain.LineDemographic     `json:"demographic,omitempty"`
}

func ToLineDailyReportResponse(r *domain.LineDailyReport) LineDailyReportResponse {
	return LineDailyReportResponse{
		ReportDate:      r.ReportDate.Format("2006-01-02"),
		Followers:       r.Followers,
		TargetedReaches: r.TargetedReaches,
		Blocks:          r.Blocks,
		MessageDelivery: r.MessageDelivery,
		Demographic:     r.Demographic,
	}
}