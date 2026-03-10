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

type GetUsersRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Status *int    `json:"status"`
	Pagination
}

type CreateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status int    `json:"status"`
}

type UpdateUserRequest struct {
	UID    string  `json:"uid"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Status *int    `json:"status"`
}

type UserResponse struct {
	ID        uint64    `json:"id"`
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse struct {
	Users []*UserResponse `json:"users"`
	Paginate
}
