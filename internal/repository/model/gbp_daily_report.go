package model

import "time"

type GBPDailyReport struct {
	ID                    uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	GBPConnectionID       int64     `gorm:"column:gbp_connection_id"`
	ReportDate            time.Time `gorm:"column:report_date;type:date"`
	ProfileViews          int       `gorm:"column:profile_views"`
	PhoneCalls            int       `gorm:"column:phone_calls"`
	DirectionRequests     int       `gorm:"column:direction_requests"`
	PhotoViews            int       `gorm:"column:photo_views"`
	ReviewCount           int       `gorm:"column:review_count"`
	AverageRating         float64   `gorm:"column:average_rating"`
	SearchQueryBreakdown  []byte    `gorm:"column:search_query_breakdown;type:json"`
	FetchedAt             time.Time `gorm:"column:fetched_at"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (GBPDailyReport) TableName() string {
	return "gbp_daily_reports"
}