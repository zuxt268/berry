package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/repository/model"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type operatorRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewOperatorRepository(
	dbDriver infrastructure.DBDriver,
) port.OperatorRepository {
	return &operatorRepository{
		dbDriver: dbDriver,
	}
}

func (r *operatorRepository) Find(ctx context.Context, f filter.Filter) (*domain.Operator, error) {
	var m model.Operator
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toOperatorDomain(&m), nil
}

func (r *operatorRepository) List(ctx context.Context, f filter.Filter) ([]*domain.Operator, error) {
	var models []*model.Operator
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	operators := make([]*domain.Operator, len(models))
	for i, m := range models {
		operators[i] = toOperatorDomain(m)
	}
	return operators, nil
}

func (r *operatorRepository) Count(ctx context.Context, f filter.Filter) (int64, error) {
	count, err := r.dbDriver.Count(ctx, &model.Operator{}, f)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *operatorRepository) Exists(ctx context.Context, f filter.Filter) (bool, error) {
	var models []*model.Operator
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return false, err
	}
	return len(models) > 0, nil
}

func (r *operatorRepository) Create(ctx context.Context, operator *domain.Operator) (*domain.Operator, error) {
	m := toOperatorModel(operator)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toOperatorDomain(m), nil
}

func (r *operatorRepository) Update(ctx context.Context, operator *domain.Operator, f filter.Filter) (*domain.Operator, error) {
	m := toOperatorModel(operator)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toOperatorDomain(m), nil
}

func (r *operatorRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.Operator{}, f)
}

func toOperatorDomain(m *model.Operator) *domain.Operator {
	return &domain.Operator{
		ID:        m.ID,
		UID:       m.UID,
		Email:     m.Email,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toOperatorModel(o *domain.Operator) *model.Operator {
	return &model.Operator{
		ID:        o.ID,
		UID:       o.UID,
		Email:     o.Email,
		Name:      o.Name,
		IsActive:  o.IsActive,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
