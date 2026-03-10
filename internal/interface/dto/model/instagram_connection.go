package model

import "time"

type InstagramConnection struct {
	ID                         int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UID                        string     `gorm:"column:uid;uniqueIndex"`
	UserID                     uint64     `gorm:"column:user_id;index"`
	InstagramBusinessAccountID string     `gorm:"column:instagram_business_account_id"`
	FacebookPageID             string     `gorm:"column:facebook_page_id"`
	AccessToken                string     `gorm:"column:access_token"`
	TokenExpiresAt             *time.Time `gorm:"column:token_expires_at"`
	ConnectedAt                time.Time  `gorm:"column:connected_at;default:CURRENT_TIMESTAMP"`
	DisconnectedAt             *time.Time `gorm:"column:disconnected_at"`
	CreatedAt                  time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt                  time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (InstagramConnection) TableName() string {
	return "instagram_connections"
}