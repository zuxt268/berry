package domain

import "time"

type LineConnection struct {
	ID                 int64
	UID                string
	UserID             uint64
	ChannelID          string
	ChannelSecret      string
	ChannelAccessToken string
	ChannelName        string
	BotUserID          string
	ConnectedAt        time.Time
	DisconnectedAt     *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}