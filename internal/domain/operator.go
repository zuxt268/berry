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

type GetOperatorsRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	IsActive *bool   `json:"is_active"`
	Pagination
}

type CreateOperatorRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateOperatorRequest struct {
	UID      string  `json:"uid"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	IsActive *bool   `json:"is_active"`
}

type OperatorResponse struct {
	ID        int64     `json:"id"`
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OperatorsResponse struct {
	Operators []*OperatorResponse `json:"operators"`
	Paginate
}
