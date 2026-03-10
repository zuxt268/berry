package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type UserRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.User, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.User, error)
	Count(ctx context.Context, f filter.Filter) (int64, error)
	Exists(ctx context.Context, f filter.Filter) (bool, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User, f filter.Filter) (*domain.User, error)
	Delete(ctx context.Context, f filter.Filter) error
}

type userRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewUserRepository(
	dbDriver infrastructure.DBDriver,
) UserRepository {
	return &userRepository{
		dbDriver: dbDriver,
	}
}

func (r *userRepository) Find(ctx context.Context, f filter.Filter) (*domain.User, error) {
	var m model.User
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toUserDomain(&m), nil
}

func (r *userRepository) List(ctx context.Context, f filter.Filter) ([]*domain.User, error) {
	var models []*model.User
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = toUserDomain(m)
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context, f filter.Filter) (int64, error) {
	count, err := r.dbDriver.Count(ctx, &model.User{}, f)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) Exists(ctx context.Context, f filter.Filter) (bool, error) {
	var models []*model.User
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return false, err
	}
	return len(models) > 0, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	m := toUserModel(user)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toUserDomain(m), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User, f filter.Filter) (*domain.User, error) {
	m := toUserModel(user)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toUserDomain(m), nil
}

func (r *userRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.User{}, f)
}

func toUserDomain(m *model.User) *domain.User {
	return &domain.User{
		ID:        m.ID,
		UID:       m.UID,
		Name:      m.Name,
		Email:     m.Email,
		Status:    domain.UserStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toUserModel(u *domain.User) *model.User {
	return &model.User{
		ID:        u.ID,
		UID:       u.UID,
		Name:      u.Name,
		Email:     u.Email,
		Status:    int(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
