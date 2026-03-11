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