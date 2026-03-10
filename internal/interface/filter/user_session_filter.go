package filter

import (
	"time"

	"gorm.io/gorm"
)

type UserSessionFilter struct {
	ID           *uint64
	UID          *string
	UserID       *uint64
	SessionToken *string
	Expired      bool
	Limit        *int
	Offset       *int
}

func (f *UserSessionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	if f.SessionToken != nil {
		db = db.Where("session_token = ?", *f.SessionToken)
	}
	if f.Expired {
		db = db.Where("expires_at < ?", time.Now())
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}
