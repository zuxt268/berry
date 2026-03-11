package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/repository/model"
	"github.com/zuxt268/berry/internal/usecase/port"
)

func NewOperatorSessionRepository(
	dbDriver infrastructure.DBDriver,
) port.OperatorSessionRepository {
	return &operatorSessionRepository{
		dbDriver: dbDriver,
	}
}

type operatorSessionRepository struct {
	dbDriver infrastructure.DBDriver
}

func (r *operatorSessionRepository) Create(ctx context.Context, session *domain.OperatorSession) error {
	m := toOperatorSessionModel(session)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return err
	}
	session.ID = m.ID
	return nil
}

func (r *operatorSessionRepository) Get(ctx context.Context, f *filter.OperatorSessionFilter) (*domain.OperatorSession, error) {
	var m model.OperatorSession
	err := r.dbDriver.First(ctx, &m, f)
	if err != nil {
		return nil, err
	}
	return toOperatorSessionEntity(&m), nil
}

func (r *operatorSessionRepository) FindAll(ctx context.Context, f *filter.OperatorSessionFilter) ([]*domain.OperatorSession, error) {
	var models []*model.OperatorSession
	err := r.dbDriver.Get(ctx, &models, f)
	if err != nil {
		return nil, err
	}
	sessions := make([]*domain.OperatorSession, len(models))
	for i, m := range models {
		sessions[i] = toOperatorSessionEntity(m)
	}
	return sessions, nil
}

func (r *operatorSessionRepository) Exists(ctx context.Context, f *filter.OperatorSessionFilter) (bool, error) {
	var models []*model.OperatorSession
	err := r.dbDriver.Get(ctx, &models, f)
	if err != nil {
		return false, err
	}
	return len(models) > 0, nil
}

func (r *operatorSessionRepository) Delete(ctx context.Context, f *filter.OperatorSessionFilter) error {
	return r.dbDriver.Delete(ctx, &model.OperatorSession{}, f)
}

func toOperatorSessionEntity(m *model.OperatorSession) *domain.OperatorSession {
	return &domain.OperatorSession{
		ID:           m.ID,
		UID:          m.UID,
		OperatorID:   m.OperatorID,
		SessionToken: m.SessionToken,
		IPAddress:    m.IPAddress,
		UserAgent:    m.UserAgent,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
	}
}

func toOperatorSessionModel(e *domain.OperatorSession) *model.OperatorSession {
	return &model.OperatorSession{
		ID:           e.ID,
		UID:          e.UID,
		OperatorID:   e.OperatorID,
		SessionToken: e.SessionToken,
		IPAddress:    e.IPAddress,
		UserAgent:    e.UserAgent,
		ExpiresAt:    e.ExpiresAt,
		CreatedAt:    e.CreatedAt,
	}
}
