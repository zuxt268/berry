package domain

import "time"

type LineDailyReport struct {
	ID               uint64
	LineConnectionID int64
	ReportDate       time.Time
	Followers        int
	TargetedReaches  int
	Blocks           int
	MessageDelivery  *LineMessageDelivery
	Demographic      *LineDemographic
	FetchedAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type LineMessageDelivery struct {
	Status      string `json:"status"`
	Success     int    `json:"success"`
	UniqueClick int    `json:"unique_click"`
	UniqueOpen  int    `json:"unique_open"`
}

type LineDemographic struct {
	Available bool                  `json:"available"`
	Genders   []LineDemographicItem `json:"genders"`
	Ages      []LineDemographicItem `json:"ages"`
	Areas     []LineDemographicItem `json:"areas"`
}

type LineDemographicItem struct {
	Key        string  `json:"key"`
	Percentage float64 `json:"percentage"`
}