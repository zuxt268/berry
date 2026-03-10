package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/lib"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type DBDriver interface {
	BeginTransaction() Transaction
	First(ctx context.Context, model interface{}, filter filter.Filter) error
	Create(ctx context.Context, model interface{}, insertIgnore bool) error
	CreateBatch(ctx context.Context, model interface{}, insertIgnore bool) error
	Get(ctx context.Context, model interface{}, filter filter.Filter) error
	Count(ctx context.Context, model interface{}, filter filter.Filter) (int64, error)
	Update(ctx context.Context, model interface{}, filter filter.Filter) error
	Delete(ctx context.Context, model interface{}, filter filter.Filter) error
	FirstForUpdate(ctx context.Context, model interface{}, filter filter.Filter) error
	GetForUpdate(ctx context.Context, model interface{}, filter filter.Filter) error
	RawSQL(ctx context.Context, query string, args []interface{}, model interface{}) error
	Upsert(ctx context.Context, model interface{}, conflictColumns []string, updateColumns []string) error
}

func NewDBDriver(readClient *gorm.DB, writeClient *gorm.DB) DBDriver {
	return &dbDriver{
		readClient:  readClient,
		writeClient: writeClient,
	}
}

type dbDriver struct {
	readClient  *gorm.DB
	writeClient *gorm.DB
}

func (d *dbDriver) First(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getReadDBFromContext(ctx)

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	err := db.First(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: record not found", domain.ErrNotFound)
		}
		return err
	}

	return nil
}

func (d *dbDriver) Create(ctx context.Context, model interface{}, insertIgnore bool) error {
	db := d.getDBFromContext(ctx)
	if insertIgnore {
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(model).Error
	}
	return db.Create(model).Error
}

const batchSize = 100

func (d *dbDriver) CreateBatch(ctx context.Context, model interface{}, insertIgnore bool) error {
	db := d.getDBFromContext(ctx)
	if insertIgnore {
		return db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(model, batchSize).Error
	}
	return d.getDBFromContext(ctx).CreateInBatches(model, batchSize).Error
}

func (d *dbDriver) Get(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getReadDBFromContext(ctx)

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	return db.Find(model).Error
}

func (d *dbDriver) Count(ctx context.Context, model interface{}, f filter.Filter) (int64, error) {
	db := d.getReadDBFromContext(ctx).Model(model)

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	var count int64
	err := db.Count(&count).Error
	return count, err
}

func (d *dbDriver) Update(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getDBFromContext(ctx).Model(model).Select("*").Omit("created_at")

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	return db.Updates(model).Error
}

func (d *dbDriver) Delete(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getDBFromContext(ctx)

	if lib.IsNil(f) {
		return domain.ErrFilterRequired
	}

	return f.Apply(db).Delete(model).Error
}

func (d *dbDriver) BeginTransaction() Transaction {
	tx := d.writeClient.Begin()
	return &Tx{
		tx: tx,
	}
}

func (d *dbDriver) RawSQL(ctx context.Context, query string, args []interface{}, model interface{}) error {
	db := d.getReadDBFromContext(ctx)
	result := db.Raw(query, args...).Scan(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *dbDriver) FirstForUpdate(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getDBFromContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	err := db.First(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: record not found", domain.ErrNotFound)
		}
		return err
	}

	return nil
}

func (d *dbDriver) GetForUpdate(ctx context.Context, model interface{}, f filter.Filter) error {
	db := d.getDBFromContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"})

	if !lib.IsNil(f) {
		db = f.Apply(db)
	}

	return db.Find(model).Error
}

func (d *dbDriver) Upsert(ctx context.Context, model interface{}, conflictColumns []string, updateColumns []string) error {
	columns := make([]clause.Column, len(conflictColumns))
	for i, col := range conflictColumns {
		columns[i] = clause.Column{Name: col}
	}

	db := d.getDBFromContext(ctx)
	return db.Clauses(clause.OnConflict{
		Columns:   columns,
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(model).Error
}

// getTransactionFromContext contextからトランザクションを取得する
func getTransactionFromContext(ctx context.Context) Transaction {
	if tx, ok := ctx.Value(domain.TxKey{}).(Transaction); ok {
		return tx
	}
	return nil
}

// getDBFromContext contextの内容に応じて適切なDB接続を返す
func (d *dbDriver) getDBFromContext(ctx context.Context) *gorm.DB {
	if tx := getTransactionFromContext(ctx); tx != nil {
		return tx.GetTx().WithContext(ctx)
	}
	return d.writeClient.WithContext(ctx)
}

// getReadDBFromContext 読み取り用のDB接続を返す
func (d *dbDriver) getReadDBFromContext(ctx context.Context) *gorm.DB {
	if tx := getTransactionFromContext(ctx); tx != nil {
		return tx.GetTx().WithContext(ctx)
	}
	return d.readClient.WithContext(ctx)
}
