package model

import "time"

type LineConnection struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UID                string     `gorm:"column:uid;uniqueIndex"`
	UserID             uint64     `gorm:"column:user_id;index"`
	ChannelID          string     `gorm:"column:channel_id"`
	ChannelSecret      string     `gorm:"column:channel_secret"`
	ChannelAccessToken string     `gorm:"column:channel_access_token"`
	ChannelName        string     `gorm:"column:channel_name"`
	BotUserID          string     `gorm:"column:bot_user_id"`
	ConnectedAt        time.Time  `gorm:"column:connected_at;default:CURRENT_TIMESTAMP"`
	DisconnectedAt     *time.Time `gorm:"column:disconnected_at"`
	CreatedAt          time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (LineConnection) TableName() string {
	return "line_connections"
}