package domain

import "time"

type Operator struct {
	ID        int64
	UID       string
	Email     string
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}