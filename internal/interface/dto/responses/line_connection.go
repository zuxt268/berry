package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

type LineConnectionResponse struct {
	UID            string     `json:"uid"`
	ChannelID      string     `json:"channel_id"`
	ChannelName    string     `json:"channel_name"`
	BotUserID      string     `json:"bot_user_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}

func ToLineConnectionResponse(c *domain.LineConnection) *LineConnectionResponse {
	return &LineConnectionResponse{
		UID:            c.UID,
		ChannelID:      c.ChannelID,
		ChannelName:    c.ChannelName,
		BotUserID:      c.BotUserID,
		ConnectedAt:    c.ConnectedAt,
		DisconnectedAt: c.DisconnectedAt,
	}
}