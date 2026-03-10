package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type GBPConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GBPConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GBPConnection, error)
	Create(ctx context.Context, conn *domain.GBPConnection) (*domain.GBPConnection, error)
	Update(ctx context.Context, conn *domain.GBPConnection, f filter.Filter) (*domain.GBPConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

type gbpConnectionRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGBPConnectionRepository(dbDriver infrastructure.DBDriver) GBPConnectionRepository {
	return &gbpConnectionRepository{dbDriver: dbDriver}
}

func (r *gbpConnectionRepository) Find(ctx context.Context, f filter.Filter) (*domain.GBPConnection, error) {
	var m model.GBPConnection
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGBPConnectionDomain(&m), nil
}

func (r *gbpConnectionRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GBPConnection, error) {
	var models []*model.GBPConnection
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	conns := make([]*domain.GBPConnection, len(models))
	for i, m := range models {
		conns[i] = toGBPConnectionDomain(m)
	}
	return conns, nil
}

func (r *gbpConnectionRepository) Create(ctx context.Context, conn *domain.GBPConnection) (*domain.GBPConnection, error) {
	m := toGBPConnectionModel(conn)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toGBPConnectionDomain(m), nil
}

func (r *gbpConnectionRepository) Update(ctx context.Context, conn *domain.GBPConnection, f filter.Filter) (*domain.GBPConnection, error) {
	m := toGBPConnectionModel(conn)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toGBPConnectionDomain(m), nil
}

func (r *gbpConnectionRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.GBPConnection{}, f)
}

func toGBPConnectionDomain(m *model.GBPConnection) *domain.GBPConnection {
	return &domain.GBPConnection{
		ID:             m.ID,
		UID:            m.UID,
		UserID:         m.UserID,
		LocationID:     m.LocationID,
		AccountID:      m.AccountID,
		RefreshToken:   m.RefreshToken,
		ConnectedAt:    m.ConnectedAt,
		DisconnectedAt: m.DisconnectedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func toGBPConnectionModel(c *domain.GBPConnection) *model.GBPConnection {
	return &model.GBPConnection{
		ID:             c.ID,
		UID:            c.UID,
		UserID:         c.UserID,
		LocationID:     c.LocationID,
		AccountID:      c.AccountID,
		RefreshToken:   c.RefreshToken,
		ConnectedAt:    c.ConnectedAt,
		DisconnectedAt: c.DisconnectedAt,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
