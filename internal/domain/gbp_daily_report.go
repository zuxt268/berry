package domain

import "time"

type GBPDailyReport struct {
	ID              uint64
	GBPConnectionID int64
	ReportDate      time.Time

	ProfileViews     int
	PhoneCalls       int
	DirectionRequests int
	PhotoViews       int
	ReviewCount      int
	AverageRating    float64

	SearchQueryBreakdown []SearchQueryBreakdown

	FetchedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchQueryBreakdown struct {
	Query       string `json:"query"`
	Impressions int    `json:"impressions"`
}