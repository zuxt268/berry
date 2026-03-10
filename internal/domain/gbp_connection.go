package domain

import "time"

type GBPConnection struct {
	ID             int64
	UID            string
	UserID         uint64
	LocationID     string
	AccountID      string
	RefreshToken   string
	ConnectedAt    time.Time
	DisconnectedAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ConnectGBPRequest struct {
	LocationID string `json:"location_id" validate:"required"`
	AccountID  string `json:"account_id" validate:"required"`
}

type GBPConnectionResponse struct {
	UID            string     `json:"uid"`
	LocationID     string     `json:"location_id"`
	AccountID      string     `json:"account_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}