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