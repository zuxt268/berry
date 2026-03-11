package domain

import "time"

type UserStatus uint8

const (
	UserStatusInactive UserStatus = 0
	UserStatusActive   UserStatus = 1
)

type User struct {
	ID        uint64
	UID       string
	Name      string
	Email     string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}