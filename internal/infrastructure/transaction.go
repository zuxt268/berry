package infrastructure

import "gorm.io/gorm"

type Transaction interface {
	Commit() error
	Rollback()
	GetTx() *gorm.DB
}

type Tx struct {
	tx *gorm.DB
}

func (t Tx) Commit() error {
	return t.tx.Commit().Error
}

func (t Tx) Rollback() {
	t.tx.Rollback()
}

func (t Tx) GetTx() *gorm.DB {
	return t.tx
}
