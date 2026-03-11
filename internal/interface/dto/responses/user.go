package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

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
	domain.Paginate
}

func ToUserResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		UID:       u.UID,
		Name:      u.Name,
		Email:     u.Email,
		Status:    int(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUsersResponse(users []*domain.User, total int64) *UsersResponse {
	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		resp[i] = ToUserResponse(u)
	}
	return &UsersResponse{
		Users: resp,
		Paginate: domain.Paginate{
			Total: total,
			Count: int64(len(resp)),
		},
	}
}
