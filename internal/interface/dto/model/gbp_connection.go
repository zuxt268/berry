package model

import "time"

type GBPConnection struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UID            string     `gorm:"column:uid;uniqueIndex"`
	UserID         uint64     `gorm:"column:user_id;index"`
	LocationID     string     `gorm:"column:location_id"`
	AccountID      string     `gorm:"column:account_id"`
	RefreshToken   string     `gorm:"column:refresh_token"`
	ConnectedAt    time.Time  `gorm:"column:connected_at;default:CURRENT_TIMESTAMP"`
	DisconnectedAt *time.Time `gorm:"column:disconnected_at"`
	CreatedAt      time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (GBPConnection) TableName() string {
	return "gbp_connections"
}