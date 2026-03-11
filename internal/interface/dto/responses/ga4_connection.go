package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

type GA4ConnectionResponse struct {
	UID              string     `json:"uid"`
	GooglePropertyID string     `json:"google_property_id"`
	ConnectedAt      time.Time  `json:"connected_at"`
	DisconnectedAt   *time.Time `json:"disconnected_at,omitempty"`
}

func ToGA4ConnectionResponse(c *domain.GA4Connection) *GA4ConnectionResponse {
	return &GA4ConnectionResponse{
		UID:              c.UID,
		GooglePropertyID: c.GooglePropertyID,
		ConnectedAt:      c.ConnectedAt,
		DisconnectedAt:   c.DisconnectedAt,
	}
}