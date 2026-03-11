package model

import "time"

type GSCDailyReport struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	GSCConnectionID  int64     `gorm:"column:gsc_connection_id"`
	ReportDate       time.Time `gorm:"column:report_date;type:date"`
	Impressions      int       `gorm:"column:impressions"`
	Clicks           int       `gorm:"column:clicks"`
	CTR              float64   `gorm:"column:ctr"`
	AveragePosition  float64   `gorm:"column:average_position"`
	KeywordBreakdown []byte    `gorm:"column:keyword_breakdown;type:json"`
	PageBreakdown    []byte    `gorm:"column:page_breakdown;type:json"`
	FetchedAt        time.Time `gorm:"column:fetched_at"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (GSCDailyReport) TableName() string {
	return "gsc_daily_reports"
}