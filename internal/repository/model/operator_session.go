package model

import "time"

type OperatorSession struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UID          string    `gorm:"column:uid;uniqueIndex"`
	OperatorID   int64     `gorm:"column:operator_id;index"`
	SessionToken string    `gorm:"column:session_token;uniqueIndex"`
	IPAddress    string    `gorm:"column:ip_address"`
	UserAgent    string    `gorm:"column:user_agent"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (OperatorSession) TableName() string {
	return "operator_sessions"
}