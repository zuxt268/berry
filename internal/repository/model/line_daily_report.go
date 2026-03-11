package model

import "time"

type LineDailyReport struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	LineConnectionID int64     `gorm:"column:line_connection_id"`
	ReportDate       time.Time `gorm:"column:report_date;type:date"`
	Followers        int       `gorm:"column:followers"`
	TargetedReaches  int       `gorm:"column:targeted_reaches"`
	Blocks           int       `gorm:"column:blocks"`
	MessageDelivery  []byte    `gorm:"column:message_delivery;type:json"`
	Demographic      []byte    `gorm:"column:demographic;type:json"`
	FetchedAt        time.Time `gorm:"column:fetched_at"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (LineDailyReport) TableName() string {
	return "line_daily_reports"
}