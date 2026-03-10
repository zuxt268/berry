package filter

import (
	"gorm.io/gorm"
)

type OperatorFilter struct {
	ID       *int64
	UID      *string
	Email    *string
	Name     *string
	IsActive *bool
	Limit    *int
	Offset   *int
}

func (f *OperatorFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UID != nil {
		db = db.Where("uid = ?", *f.UID)
	}
	if f.Email != nil {
		db = db.Where("email = ?", *f.Email)
	}
	if f.Name != nil {
		db = db.Where("name = ?", *f.Name)
	}
	if f.IsActive != nil {
		db = db.Where("is_active = ?", *f.IsActive)
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}