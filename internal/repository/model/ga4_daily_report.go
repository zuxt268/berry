package model

import "time"

type GA4DailyReport struct {
	ID                 uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	GA4ConnectionID    int64     `gorm:"column:ga4_connection_id"`
	ReportDate         time.Time `gorm:"column:report_date;type:date"`
	Sessions           int       `gorm:"column:sessions"`
	TotalUsers         int       `gorm:"column:total_users"`
	BounceRate         float64   `gorm:"column:bounce_rate"`
	AvgSessionDuration float64   `gorm:"column:avg_session_duration"`
	Conversions        int       `gorm:"column:conversions"`
	ChannelBreakdown   []byte    `gorm:"column:channel_breakdown;type:json"`
	DeviceBreakdown    []byte    `gorm:"column:device_breakdown;type:json"`
	PageBreakdown      []byte    `gorm:"column:page_breakdown;type:json"`
	FetchedAt          time.Time `gorm:"column:fetched_at"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
}

func (GA4DailyReport) TableName() string {
	return "ga4_daily_reports"
}