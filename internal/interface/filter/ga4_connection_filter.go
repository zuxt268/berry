package filter

import (
	"gorm.io/gorm"
)

type GA4ConnectionFilter struct {
	ID               *int64
	UID              *string
	UserID           *uint64
	GooglePropertyID *string
	ActiveOnly       bool
	Limit            *int
	Offset           *int
}

func (f *GA4ConnectionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UID != nil {
		db = db.Where("uid = ?", *f.UID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	if f.GooglePropertyID != nil {
		db = db.Where("google_property_id = ?", *f.GooglePropertyID)
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