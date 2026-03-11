package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

type GBPConnectionResponse struct {
	UID            string     `json:"uid"`
	LocationID     string     `json:"location_id"`
	AccountID      string     `json:"account_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}

func ToGBPConnectionResponse(c *domain.GBPConnection) *GBPConnectionResponse {
	return &GBPConnectionResponse{
		UID:            c.UID,
		LocationID:     c.LocationID,
		AccountID:      c.AccountID,
		ConnectedAt:    c.ConnectedAt,
		DisconnectedAt: c.DisconnectedAt,
	}
}