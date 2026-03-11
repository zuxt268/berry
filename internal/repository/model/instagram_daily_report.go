package model

import "time"

type InstagramDailyReport struct {
	ID                      uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	InstagramConnectionID   int64     `gorm:"column:instagram_connection_id"`
	ReportDate              time.Time `gorm:"column:report_date;type:date"`
	FollowerCount           int       `gorm:"column:follower_count"`
	Impressions             int       `gorm:"column:impressions"`
	Reach                   int       `gorm:"column:reach"`
	ProfileViews            int       `gorm:"column:profile_views"`
	WebsiteClicks           int       `gorm:"column:website_clicks"`
	PostEngagement          []byte    `gorm:"column:post_engagement;type:json"`
	AudienceDemographics    []byte    `gorm:"column:audience_demographics;type:json"`
	StoriesInsights         []byte    `gorm:"column:stories_insights;type:json"`
	FetchedAt               time.Time `gorm:"column:fetched_at"`
	CreatedAt               time.Time `gorm:"column:created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at"`
}

func (InstagramDailyReport) TableName() string {
	return "instagram_daily_reports"
}