package filter

import (
	"gorm.io/gorm"
)

type UserFilter struct {
	ID     *uint64
	UID    *string
	Name   *string
	Email  *string
	Status *int
	Limit  *int
	Offset *int
}

func (f *UserFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UID != nil {
		db = db.Where("uid = ?", *f.UID)
	}
	if f.Name != nil {
		db = db.Where("name = ?", *f.Name)
	}
	if f.Email != nil {
		db = db.Where("email = ?", *f.Email)
	}
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}
