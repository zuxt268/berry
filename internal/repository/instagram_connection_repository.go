package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type InstagramConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.InstagramConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.InstagramConnection, error)
	Create(ctx context.Context, conn *domain.InstagramConnection) (*domain.InstagramConnection, error)
	Update(ctx context.Context, conn *domain.InstagramConnection, f filter.Filter) (*domain.InstagramConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

type instagramConnectionRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewInstagramConnectionRepository(dbDriver infrastructure.DBDriver) InstagramConnectionRepository {
	return &instagramConnectionRepository{dbDriver: dbDriver}
}

func (r *instagramConnectionRepository) Find(ctx context.Context, f filter.Filter) (*domain.InstagramConnection, error) {
	var m model.InstagramConnection
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toInstagramConnectionDomain(&m), nil
}

func (r *instagramConnectionRepository) List(ctx context.Context, f filter.Filter) ([]*domain.InstagramConnection, error) {
	var models []*model.InstagramConnection
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	conns := make([]*domain.InstagramConnection, len(models))
	for i, m := range models {
		conns[i] = toInstagramConnectionDomain(m)
	}
	return conns, nil
}

func (r *instagramConnectionRepository) Create(ctx context.Context, conn *domain.InstagramConnection) (*domain.InstagramConnection, error) {
	m := toInstagramConnectionModel(conn)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toInstagramConnectionDomain(m), nil
}

func (r *instagramConnectionRepository) Update(ctx context.Context, conn *domain.InstagramConnection, f filter.Filter) (*domain.InstagramConnection, error) {
	m := toInstagramConnectionModel(conn)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toInstagramConnectionDomain(m), nil
}

func (r *instagramConnectionRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.InstagramConnection{}, f)
}

func toInstagramConnectionDomain(m *model.InstagramConnection) *domain.InstagramConnection {
	return &domain.InstagramConnection{
		ID:                         m.ID,
		UID:                        m.UID,
		UserID:                     m.UserID,
		InstagramBusinessAccountID: m.InstagramBusinessAccountID,
		FacebookPageID:             m.FacebookPageID,
		AccessToken:                m.AccessToken,
		TokenExpiresAt:             m.TokenExpiresAt,
		ConnectedAt:                m.ConnectedAt,
		DisconnectedAt:             m.DisconnectedAt,
		CreatedAt:                  m.CreatedAt,
		UpdatedAt:                  m.UpdatedAt,
	}
}

func toInstagramConnectionModel(c *domain.InstagramConnection) *model.InstagramConnection {
	return &model.InstagramConnection{
		ID:                         c.ID,
		UID:                        c.UID,
		UserID:                     c.UserID,
		InstagramBusinessAccountID: c.InstagramBusinessAccountID,
		FacebookPageID:             c.FacebookPageID,
		AccessToken:                c.AccessToken,
		TokenExpiresAt:             c.TokenExpiresAt,
		ConnectedAt:                c.ConnectedAt,
		DisconnectedAt:             c.DisconnectedAt,
		CreatedAt:                  c.CreatedAt,
		UpdatedAt:                  c.UpdatedAt,
	}
}
