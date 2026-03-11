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