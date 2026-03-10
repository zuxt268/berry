package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/dto/model"
	"github.com/zuxt268/berry/internal/interface/filter"
)

type LineConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.LineConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.LineConnection, error)
	Create(ctx context.Context, conn *domain.LineConnection) (*domain.LineConnection, error)
	Update(ctx context.Context, conn *domain.LineConnection, f filter.Filter) (*domain.LineConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

type lineConnectionRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewLineConnectionRepository(dbDriver infrastructure.DBDriver) LineConnectionRepository {
	return &lineConnectionRepository{dbDriver: dbDriver}
}

func (r *lineConnectionRepository) Find(ctx context.Context, f filter.Filter) (*domain.LineConnection, error) {
	var m model.LineConnection
	if err := r.dbDriver.First(ctx, &m, f); err != nil {
		return nil, err
	}
	return toLineConnectionDomain(&m), nil
}

func (r *lineConnectionRepository) List(ctx context.Context, f filter.Filter) ([]*domain.LineConnection, error) {
	var models []*model.LineConnection
	if err := r.dbDriver.Get(ctx, &models, f); err != nil {
		return nil, err
	}
	conns := make([]*domain.LineConnection, len(models))
	for i, m := range models {
		conns[i] = toLineConnectionDomain(m)
	}
	return conns, nil
}

func (r *lineConnectionRepository) Create(ctx context.Context, conn *domain.LineConnection) (*domain.LineConnection, error) {
	m := toLineConnectionModel(conn)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return nil, err
	}
	return toLineConnectionDomain(m), nil
}

func (r *lineConnectionRepository) Update(ctx context.Context, conn *domain.LineConnection, f filter.Filter) (*domain.LineConnection, error) {
	m := toLineConnectionModel(conn)
	if err := r.dbDriver.Update(ctx, m, f); err != nil {
		return nil, err
	}
	return toLineConnectionDomain(m), nil
}

func (r *lineConnectionRepository) Delete(ctx context.Context, f filter.Filter) error {
	return r.dbDriver.Delete(ctx, &model.LineConnection{}, f)
}

func toLineConnectionDomain(m *model.LineConnection) *domain.LineConnection {
	return &domain.LineConnection{
		ID:                 m.ID,
		UID:                m.UID,
		UserID:             m.UserID,
		ChannelID:          m.ChannelID,
		ChannelSecret:      m.ChannelSecret,
		ChannelAccessToken: m.ChannelAccessToken,
		ChannelName:        m.ChannelName,
		BotUserID:          m.BotUserID,
		ConnectedAt:        m.ConnectedAt,
		DisconnectedAt:     m.DisconnectedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func toLineConnectionModel(c *domain.LineConnection) *model.LineConnection {
	return &model.LineConnection{
		ID:                 c.ID,
		UID:                c.UID,
		UserID:             c.UserID,
		ChannelID:          c.ChannelID,
		ChannelSecret:      c.ChannelSecret,
		ChannelAccessToken: c.ChannelAccessToken,
		ChannelName:        c.ChannelName,
		BotUserID:          c.BotUserID,
		ConnectedAt:        c.ConnectedAt,
		DisconnectedAt:     c.DisconnectedAt,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
