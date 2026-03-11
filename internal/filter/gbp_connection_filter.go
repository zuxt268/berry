package filter

import (
	"gorm.io/gorm"
)

type GBPConnectionFilter struct {
	ID         *int64
	UID        *string
	UserID     *uint64
	LocationID *string
	AccountID  *string
	ActiveOnly bool
	Limit      *int
	Offset     *int
}

func (f *GBPConnectionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UID != nil {
		db = db.Where("uid = ?", *f.UID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	if f.LocationID != nil {
		db = db.Where("location_id = ?", *f.LocationID)
	}
	if f.AccountID != nil {
		db = db.Where("account_id = ?", *f.AccountID)
	}
	if f.ActiveOnly {
		db = db.Where("disconnected_at IS NULL")
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}