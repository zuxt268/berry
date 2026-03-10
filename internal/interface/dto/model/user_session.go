package model

import "time"

type UserSession struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UID          string    `gorm:"column:uid;uniqueIndex"`
	UserID       uint64    `gorm:"column:user_id;index"`
	SessionToken string    `gorm:"column:session_token;uniqueIndex;"`
	IPAddress    string    `gorm:"column:ip_address"`
	UserAgent    string    `gorm:"column:user_agent"`
	ExpiresAt    time.Time `gorm:"column:expires_at;index"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
