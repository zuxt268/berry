package domain

import "time"

type GA4DailyReport struct {
	ID              uint64
	GA4ConnectionID int64
	ReportDate      time.Time

	Sessions           int
	TotalUsers         int
	BounceRate         float64
	AvgSessionDuration float64
	Conversions        int

	ChannelBreakdown []ChannelBreakdown
	DeviceBreakdown  []DeviceBreakdown
	PageBreakdown    []PageBreakdown

	FetchedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelBreakdown struct {
	Channel  string `json:"channel"`
	Sessions int    `json:"sessions"`
	Users    int    `json:"users"`
}

type DeviceBreakdown struct {
	DeviceCategory string `json:"device_category"`
	Sessions       int    `json:"sessions"`
	Users          int    `json:"users"`
}

type PageBreakdown struct {
	PagePath      string  `json:"page_path"`
	PageViews     int     `json:"page_views"`
	AvgTimeOnPage float64 `json:"avg_time_on_page"`
}