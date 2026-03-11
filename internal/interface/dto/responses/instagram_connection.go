package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

type InstagramConnectionResponse struct {
	UID                        string     `json:"uid"`
	InstagramBusinessAccountID string     `json:"instagram_business_account_id"`
	FacebookPageID             string     `json:"facebook_page_id"`
	TokenExpiresAt             *time.Time `json:"token_expires_at,omitempty"`
	ConnectedAt                time.Time  `json:"connected_at"`
	DisconnectedAt             *time.Time `json:"disconnected_at,omitempty"`
}

func ToInstagramConnectionResponse(c *domain.InstagramConnection) *InstagramConnectionResponse {
	return &InstagramConnectionResponse{
		UID:                        c.UID,
		InstagramBusinessAccountID: c.InstagramBusinessAccountID,
		FacebookPageID:             c.FacebookPageID,
		TokenExpiresAt:             c.TokenExpiresAt,
		ConnectedAt:                c.ConnectedAt,
		DisconnectedAt:             c.DisconnectedAt,
	}
}