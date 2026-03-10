package domain

import "time"

type InstagramConnection struct {
	ID                         int64
	UID                        string
	UserID                     uint64
	InstagramBusinessAccountID string
	FacebookPageID             string
	AccessToken                string
	TokenExpiresAt             *time.Time
	ConnectedAt                time.Time
	DisconnectedAt             *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type ConnectInstagramRequest struct {
	InstagramBusinessAccountID string `json:"instagram_business_account_id" validate:"required"`
}

type InstagramConnectionResponse struct {
	UID                        string     `json:"uid"`
	InstagramBusinessAccountID string     `json:"instagram_business_account_id"`
	FacebookPageID             string     `json:"facebook_page_id"`
	TokenExpiresAt             *time.Time `json:"token_expires_at,omitempty"`
	ConnectedAt                time.Time  `json:"connected_at"`
	DisconnectedAt             *time.Time `json:"disconnected_at,omitempty"`
}