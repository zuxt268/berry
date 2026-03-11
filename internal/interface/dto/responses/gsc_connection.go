package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

type GSCConnectionResponse struct {
	UID            string     `json:"uid"`
	SiteURL        string     `json:"site_url"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}

func ToGSCConnectionResponse(c *domain.GSCConnection) *GSCConnectionResponse {
	return &GSCConnectionResponse{
		UID:            c.UID,
		SiteURL:        c.SiteURL,
		ConnectedAt:    c.ConnectedAt,
		DisconnectedAt: c.DisconnectedAt,
	}
}