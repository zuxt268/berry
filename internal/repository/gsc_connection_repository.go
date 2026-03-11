package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/repository/model"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type gscConnectionRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewGSCConnectionRepository(dbDriver infrastructure.DBDriver) port.GSCConnectionRepository {
	return &gscConnectionRepository{dbDriver: dbDriver}
}

func (r *gscConnectionRepository) Find(ctx context.Context, f filter.Filter) (*domain.GSCConnection, error) {
	var m model.GSCConnection
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toGSCConnectionDomain(&m), nil
}

func (r *gscConnectionRepository) List(ctx context.Context, f filter.Filter) ([]*domain.GSCConnection, error) {
	var models []*model.GSCConnection
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	conns := make([]*domain.GSCConnection, len(models))
	for i, m := range models {
		conns[i] = toGSCConnectionDomain(m)
	}
	return conns, nil
}

func (r *gscConnectionRepository) Create(ctx context.Context, conn *domain.GSCConnection) (*domain.GSCConnection, error) {
	m := toGSCConnectionModel(conn)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toGSCConnectionDomain(m), nil
}

func (r *gscConnectionRepository) Update(ctx context.Context, conn *domain.GSCConnection, f filter.Filter) (*domain.GSCConnection, error) {
	m := toGSCConnectionModel(conn)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toGSCConnectionDomain(m), nil
}

func (r *gscConnectionRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.GSCConnection{}, f)
}

func toGSCConnectionDomain(m *model.GSCConnection) *domain.GSCConnection {
	return &domain.GSCConnection{
		ID:             m.ID,
		UID:            m.UID,
		UserID:         m.UserID,
		SiteURL:        m.SiteURL,
		RefreshToken:   m.RefreshToken,
		ConnectedAt:    m.ConnectedAt,
		DisconnectedAt: m.DisconnectedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func toGSCConnectionModel(c *domain.GSCConnection) *model.GSCConnection {
	return &model.GSCConnection{
		ID:             c.ID,
		UID:            c.UID,
		UserID:         c.UserID,
		SiteURL:        c.SiteURL,
		RefreshToken:   c.RefreshToken,
		ConnectedAt:    c.ConnectedAt,
		DisconnectedAt: c.DisconnectedAt,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
