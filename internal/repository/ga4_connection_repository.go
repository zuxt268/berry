package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type GA4ConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GA4Connection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GA4Connection, error)
	Create(ctx context.Context, conn *domain.GA4Connection) (*domain.GA4Connection, error)
	Update(ctx context.Context, conn *domain.GA4Connection, f filter.Filter) (*domain.GA4Connection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

type ga4ConnectionRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGA4ConnectionRepository(dbDriver infrastructure.DBDriver) GA4ConnectionRepository {
	return &ga4ConnectionRepository{dbDriver: dbDriver}
}

func (r *ga4ConnectionRepository) Find(ctx context.Context, f filter.Filter) (*domain.GA4Connection, error) {
	var m model.GA4Connection
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGA4ConnectionDomain(&m), nil
}

func (r *ga4ConnectionRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GA4Connection, error) {
	var models []*model.GA4Connection
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	conns := make([]*domain.GA4Connection, len(models))
	for i, m := range models {
		conns[i] = toGA4ConnectionDomain(m)
	}
	return conns, nil
}

func (r *ga4ConnectionRepository) Create(ctx context.Context, conn *domain.GA4Connection) (*domain.GA4Connection, error) {
	m := toGA4ConnectionModel(conn)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toGA4ConnectionDomain(m), nil
}

func (r *ga4ConnectionRepository) Update(ctx context.Context, conn *domain.GA4Connection, f filter.Filter) (*domain.GA4Connection, error) {
	m := toGA4ConnectionModel(conn)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toGA4ConnectionDomain(m), nil
}

func (r *ga4ConnectionRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.GA4Connection{}, f)
}

func toGA4ConnectionDomain(m *model.GA4Connection) *domain.GA4Connection {
	return &domain.GA4Connection{
		ID:               m.ID,
		UID:              m.UID,
		UserID:           m.UserID,
		GooglePropertyID: m.GooglePropertyID,
		RefreshToken:     m.RefreshToken,
		ConnectedAt:      m.ConnectedAt,
		DisconnectedAt:   m.DisconnectedAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func toGA4ConnectionModel(c *domain.GA4Connection) *model.GA4Connection {
	return &model.GA4Connection{
		ID:               c.ID,
		UID:              c.UID,
		UserID:           c.UserID,
		GooglePropertyID: c.GooglePropertyID,
		RefreshToken:     c.RefreshToken,
		ConnectedAt:      c.ConnectedAt,
		DisconnectedAt:   c.DisconnectedAt,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
