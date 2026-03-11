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