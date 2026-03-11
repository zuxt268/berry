package usecase

import "github.com/zuxt268/berry/internal/domain"

// User input types

type GetUsersInput struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Status *int    `json:"status"`
	domain.Pagination
}

type CreateUserInput struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status int    `json:"status"`
}

type UpdateUserInput struct {
	UID    string  `json:"uid"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Status *int    `json:"status"`
}

// Operator input types

type GetOperatorsInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	IsActive *bool   `json:"is_active"`
	domain.Pagination
}

type CreateOperatorInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateOperatorInput struct {
	UID      string  `json:"uid"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	IsActive *bool   `json:"is_active"`
}

// LINE input types

type ConnectLineInput struct {
	ChannelID          string `json:"channel_id" validate:"required"`
	ChannelSecret      string `json:"channel_secret" validate:"required"`
	ChannelAccessToken string `json:"channel_access_token" validate:"required"`
	ChannelName        string `json:"channel_name"`
}