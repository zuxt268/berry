package responses

import (
	"time"

	"github.com/zuxt268/berry/internal/domain"
)

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
	domain.Paginate
}

func ToOperatorResponse(o *domain.Operator) *OperatorResponse {
	return &OperatorResponse{
		ID:        o.ID,
		UID:       o.UID,
		Name:      o.Name,
		Email:     o.Email,
		IsActive:  o.IsActive,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func ToOperatorsResponse(operators []*domain.Operator, total int64) *OperatorsResponse {
	resp := make([]*OperatorResponse, len(operators))
	for i, op := range operators {
		resp[i] = ToOperatorResponse(op)
	}
	return &OperatorsResponse{
		Operators: resp,
		Paginate: domain.Paginate{
			Total: total,
			Count: int64(len(resp)),
		},
	}
}