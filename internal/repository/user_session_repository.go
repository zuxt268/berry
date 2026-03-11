package repository

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/repository/model"
	"github.com/zuxt268/berry/internal/usecase/port"
)

func NewUserSessionRepository(
	dbDriver infrastructure.DBDriver,
) port.UserSessionRepository {
	return &customerUserSessionRepository{
		dbDriver: dbDriver,
	}
}

type customerUserSessionRepository struct {
	dbDriver infrastructure.DBDriver
}

func (r *customerUserSessionRepository) Create(ctx context.Context, session *domain.UserSession) error {
	m := toUserSessionModel(session)
	if err := r.dbDriver.Create(ctx, m, false); err != nil {
		return err
	}
	session.ID = m.ID
	return nil
}

func (r *customerUserSessionRepository) Get(ctx context.Context, f *filter.UserSessionFilter) (*domain.UserSession, error) {
	var m model.UserSession
	err := r.dbDriver.First(ctx, &m, f)
	if err != nil {
		return nil, err
	}
	return toUserSessionEntity(&m), nil
}

func (r *customerUserSessionRepository) FindAll(ctx context.Context, f *filter.UserSessionFilter) ([]*domain.UserSession, error) {
	var models []*model.UserSession
	err := r.dbDriver.Get(ctx, &models, f)
	if err != nil {
		return nil, err
	}
	sessions := make([]*domain.UserSession, len(models))
	for i, m := range models {
		sessions[i] = toUserSessionEntity(m)
	}
	return sessions, nil
}

func (r *customerUserSessionRepository) Exists(ctx context.Context, f *filter.UserSessionFilter) (bool, error) {
	var models []*model.UserSession
	err := r.dbDriver.Get(ctx, &models, f)
	if err != nil {
		return false, err
	}
	return len(models) > 0, nil
}

func (r *customerUserSessionRepository) Delete(ctx context.Context, f *filter.UserSessionFilter) error {
	return r.dbDriver.Delete(ctx, &model.UserSession{}, f)
}

func toUserSessionEntity(m *model.UserSession) *domain.UserSession {
	return &domain.UserSession{
		ID:           m.ID,
		UID:          m.UID,
		UserID:       m.UserID,
		SessionToken: m.SessionToken,
		IPAddress:    m.IPAddress,
		UserAgent:    m.UserAgent,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func toUserSessionModel(e *domain.UserSession) *model.UserSession {
	return &model.UserSession{
		ID:           e.ID,
		UID:          e.UID,
		UserID:       e.UserID,
		SessionToken: e.SessionToken,
		IPAddress:    e.IPAddress,
		UserAgent:    e.UserAgent,
		ExpiresAt:    e.ExpiresAt,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
