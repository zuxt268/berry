package domain

import "time"

type GSCConnection struct {
	ID             int64
	UID            string
	UserID         uint64
	SiteURL        string
	RefreshToken   string
	ConnectedAt    time.Time
	DisconnectedAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ConnectGSCRequest struct {
	SiteURL string `json:"site_url" validate:"required"`
}

type GSCConnectionResponse struct {
	UID            string     `json:"uid"`
	SiteURL        string     `json:"site_url"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}