package filter

import (
	"time"

	"gorm.io/gorm"
)

type OperatorSessionFilter struct {
	ID           *int64
	UID          *string
	OperatorID   *int64
	SessionToken *string
	Expired      bool
	Limit        *int
	Offset       *int
}

func (f *OperatorSessionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.OperatorID != nil {
		db = db.Where("operator_id = ?", *f.OperatorID)
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