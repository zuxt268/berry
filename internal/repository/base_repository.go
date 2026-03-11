package repository

import (
	"context"
	"log/slog"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/usecase/port"
)

func NewBaseRepository(
	dbDriver infrastructure.DBDriver,
) port.BaseRepository {
	return &baseRepository{
		dbDriver: dbDriver,
	}
}

type baseRepository struct {
	dbDriver infrastructure.DBDriver
}

// WithTransaction トランザクション内で関数を実行する
func (r *baseRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {

	// ネストされたWithTransactionには対応しない
	if _, ok := ctx.Value(domain.TxKey{}).(infrastructure.Transaction); ok {
		return domain.ErrTransactionAlreadyExists
	}

	tx := r.dbDriver.BeginTransaction()

	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("transaction panic recovered", "panic", rec)
			tx.Rollback()
			err = domain.ErrTransactionPanic
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				tx.Rollback()
				err = domain.ErrTransactionCommit
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, domain.TxKey{}, tx)
	return fn(ctxWithTx)
}
