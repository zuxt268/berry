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

type ConnectLineRequest struct {
	ChannelID          string `json:"channel_id" validate:"required"`
	ChannelSecret      string `json:"channel_secret" validate:"required"`
	ChannelAccessToken string `json:"channel_access_token" validate:"required"`
	ChannelName        string `json:"channel_name"`
}

type LineConnectionResponse struct {
	UID            string     `json:"uid"`
	ChannelID      string     `json:"channel_id"`
	ChannelName    string     `json:"channel_name"`
	BotUserID      string     `json:"bot_user_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}