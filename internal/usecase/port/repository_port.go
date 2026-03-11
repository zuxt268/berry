package port

import (
	"context"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/filter"
)

// BaseRepository トランザクション管理のインターフェース
type BaseRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.User, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.User, error)
	Count(ctx context.Context, f filter.Filter) (int64, error)
	Exists(ctx context.Context, f filter.Filter) (bool, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User, f filter.Filter) (*domain.User, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// UserSessionRepository ユーザーセッションリポジトリのインターフェース
type UserSessionRepository interface {
	Create(ctx context.Context, session *domain.UserSession) error
	Get(ctx context.Context, f *filter.UserSessionFilter) (*domain.UserSession, error)
	FindAll(ctx context.Context, f *filter.UserSessionFilter) ([]*domain.UserSession, error)
	Exists(ctx context.Context, f *filter.UserSessionFilter) (bool, error)
	Delete(ctx context.Context, f *filter.UserSessionFilter) error
}

// OperatorRepository オペレーターリポジトリのインターフェース
type OperatorRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.Operator, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.Operator, error)
	Count(ctx context.Context, f filter.Filter) (int64, error)
	Exists(ctx context.Context, f filter.Filter) (bool, error)
	Create(ctx context.Context, operator *domain.Operator) (*domain.Operator, error)
	Update(ctx context.Context, operator *domain.Operator, f filter.Filter) (*domain.Operator, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// OperatorSessionRepository オペレーターセッションリポジトリのインターフェース
type OperatorSessionRepository interface {
	Create(ctx context.Context, session *domain.OperatorSession) error
	Get(ctx context.Context, f *filter.OperatorSessionFilter) (*domain.OperatorSession, error)
	FindAll(ctx context.Context, f *filter.OperatorSessionFilter) ([]*domain.OperatorSession, error)
	Exists(ctx context.Context, f *filter.OperatorSessionFilter) (bool, error)
	Delete(ctx context.Context, f *filter.OperatorSessionFilter) error
}

// GA4ConnectionRepository GA4接続リポジトリのインターフェース
type GA4ConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GA4Connection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GA4Connection, error)
	Create(ctx context.Context, conn *domain.GA4Connection) (*domain.GA4Connection, error)
	Update(ctx context.Context, conn *domain.GA4Connection, f filter.Filter) (*domain.GA4Connection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// GA4DailyReportRepository GA4日次レポートリポジトリのインターフェース
type GA4DailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GA4DailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GA4DailyReport, error)
	Upsert(ctx context.Context, report *domain.GA4DailyReport) error
}

// GBPConnectionRepository GBP接続リポジトリのインターフェース
type GBPConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GBPConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GBPConnection, error)
	Create(ctx context.Context, conn *domain.GBPConnection) (*domain.GBPConnection, error)
	Update(ctx context.Context, conn *domain.GBPConnection, f filter.Filter) (*domain.GBPConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// GBPDailyReportRepository GBP日次レポートリポジトリのインターフェース
type GBPDailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GBPDailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GBPDailyReport, error)
	Upsert(ctx context.Context, report *domain.GBPDailyReport) error
}

// GSCConnectionRepository GSC接続リポジトリのインターフェース
type GSCConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GSCConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GSCConnection, error)
	Create(ctx context.Context, conn *domain.GSCConnection) (*domain.GSCConnection, error)
	Update(ctx context.Context, conn *domain.GSCConnection, f filter.Filter) (*domain.GSCConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// GSCDailyReportRepository GSC日次レポートリポジトリのインターフェース
type GSCDailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.GSCDailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.GSCDailyReport, error)
	Upsert(ctx context.Context, report *domain.GSCDailyReport) error
}

// InstagramConnectionRepository Instagram接続リポジトリのインターフェース
type InstagramConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.InstagramConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.InstagramConnection, error)
	Create(ctx context.Context, conn *domain.InstagramConnection) (*domain.InstagramConnection, error)
	Update(ctx context.Context, conn *domain.InstagramConnection, f filter.Filter) (*domain.InstagramConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// InstagramDailyReportRepository Instagram日次レポートリポジトリのインターフェース
type InstagramDailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.InstagramDailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.InstagramDailyReport, error)
	Upsert(ctx context.Context, report *domain.InstagramDailyReport) error
}

// LineConnectionRepository LINE接続リポジトリのインターフェース
type LineConnectionRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.LineConnection, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.LineConnection, error)
	Create(ctx context.Context, conn *domain.LineConnection) (*domain.LineConnection, error)
	Update(ctx context.Context, conn *domain.LineConnection, f filter.Filter) (*domain.LineConnection, error)
	Delete(ctx context.Context, f filter.Filter) error
}

// LineDailyReportRepository LINE日次レポートリポジトリのインターフェース
type LineDailyReportRepository interface {
	Find(ctx context.Context, f filter.Filter) (*domain.LineDailyReport, error)
	List(ctx context.Context, f filter.Filter) ([]*domain.LineDailyReport, error)
	Upsert(ctx context.Context, report *domain.LineDailyReport) error
}