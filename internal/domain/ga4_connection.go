package domain

import "time"

type GA4Connection struct {
	ID               int64
	UID              string
	UserID           uint64
	GooglePropertyID string
	RefreshToken     string
	ConnectedAt      time.Time
	DisconnectedAt   *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ConnectGA4Request struct {
	GooglePropertyID string `json:"google_property_id" validate:"required"`
}

type GA4ConnectionResponse struct {
	UID              string     `json:"uid"`
	GooglePropertyID string     `json:"google_property_id"`
	ConnectedAt      time.Time  `json:"connected_at"`
	DisconnectedAt   *time.Time `json:"disconnected_at,omitempty"`
}