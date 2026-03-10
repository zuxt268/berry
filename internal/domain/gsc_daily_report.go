package domain

import "time"

type GSCDailyReport struct {
	ID              uint64
	GSCConnectionID int64
	ReportDate      time.Time

	Impressions     int
	Clicks          int
	CTR             float64
	AveragePosition float64

	KeywordBreakdown []KeywordBreakdown
	PageBreakdown    []GSCPageBreakdown

	FetchedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type KeywordBreakdown struct {
	Query           string  `json:"query"`
	Impressions     int     `json:"impressions"`
	Clicks          int     `json:"clicks"`
	CTR             float64 `json:"ctr"`
	AveragePosition float64 `json:"average_position"`
}

type GSCPageBreakdown struct {
	Page            string  `json:"page"`
	Impressions     int     `json:"impressions"`
	Clicks          int     `json:"clicks"`
	CTR             float64 `json:"ctr"`
	AveragePosition float64 `json:"average_position"`
}