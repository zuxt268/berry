package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/repository"
)

type OperatorUsecase interface {
	GetByUID(ctx context.Context, uid string) (*domain.Operator, error)
	Gets(ctx context.Context, input domain.GetOperatorsRequest) (*domain.OperatorsResponse, error)
	Create(ctx context.Context, input domain.CreateOperatorRequest) (*domain.OperatorResponse, error)
	Update(ctx context.Context, input domain.UpdateOperatorRequest) (*domain.OperatorResponse, error)
	Delete(ctx context.Context, uid string) error
}

type operatorUsecase struct {
	baseRepository     repository.BaseRepository
	operatorRepository repository.OperatorRepository
}

func NewOperatorUsecase(
	baseRepository repository.BaseRepository,
	operatorRepository repository.OperatorRepository,
) OperatorUsecase {
	return &operatorUsecase{
		baseRepository:     baseRepository,
		operatorRepository: operatorRepository,
	}
}

func (u *operatorUsecase) GetByUID(ctx context.Context, uid string) (*domain.Operator, error) {
	f := &filter.OperatorFilter{UID: &uid}
	operator, err := u.operatorRepository.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	return operator, nil
}

func (u *operatorUsecase) Gets(ctx context.Context, input domain.GetOperatorsRequest) (*domain.OperatorsResponse, error) {
	f := &filter.OperatorFilter{
		Name:     input.Name,
		Email:    input.Email,
		IsActive: input.IsActive,
	}

	operators, err := u.operatorRepository.List(ctx, f)
	if err != nil {
		return nil, err
	}

	total, err := u.operatorRepository.Count(ctx, f)
	if err != nil {
		return nil, err
	}

	return toOperatorsResponse(operators, total), nil
}

func (u *operatorUsecase) Create(ctx context.Context, input domain.CreateOperatorRequest) (*domain.OperatorResponse, error) {
	operator := &domain.Operator{
		UID:      uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		IsActive: true,
	}

	exists, err := u.operatorRepository.Exists(ctx, &filter.OperatorFilter{Email: &operator.Email})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("operator already exists")
	}

	created, err := u.operatorRepository.Create(ctx, operator)
	if err != nil {
		return nil, err
	}

	return toOperatorResponse(created), nil
}

func (u *operatorUsecase) Update(ctx context.Context, input domain.UpdateOperatorRequest) (*domain.OperatorResponse, error) {
	f := &filter.OperatorFilter{UID: &input.UID}

	existing, err := u.operatorRepository.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}

	updated, err := u.operatorRepository.Update(ctx, existing, f)
	if err != nil {
		return nil, err
	}

	return toOperatorResponse(updated), nil
}

func (u *operatorUsecase) Delete(ctx context.Context, uid string) error {
	f := &filter.OperatorFilter{UID: &uid}

	if _, err := u.operatorRepository.Find(ctx, f); err != nil {
		return err
	}

	return u.operatorRepository.Delete(ctx, f)
}

func toOperatorResponse(o *domain.Operator) *domain.OperatorResponse {
	return &domain.OperatorResponse{
		ID:        o.ID,
		UID:       o.UID,
		Name:      o.Name,
		Email:     o.Email,
		IsActive:  o.IsActive,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func toOperatorsResponse(operators []*domain.Operator, total int64) *domain.OperatorsResponse {
	resp := make([]*domain.OperatorResponse, len(operators))
	for i, op := range operators {
		resp[i] = toOperatorResponse(op)
	}
	return &domain.OperatorsResponse{
		Operators: resp,
		Paginate: domain.Paginate{
			Total: total,
			Count: int64(len(resp)),
		},
	}
}
