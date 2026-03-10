package filter

import (
	"time"

	"gorm.io/gorm"
)

type InstagramConnectionFilter struct {
	ID                         *int64
	UID                        *string
	UserID                     *uint64
	InstagramBusinessAccountID *string
	ActiveOnly                 bool
	TokenExpiringBefore        *time.Time
	Limit                      *int
	Offset                     *int
}

func (f *InstagramConnectionFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UID != nil {
		db = db.Where("uid = ?", *f.UID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	if f.InstagramBusinessAccountID != nil {
		db = db.Where("instagram_business_account_id = ?", *f.InstagramBusinessAccountID)
	}
	if f.ActiveOnly {
		db = db.Where("disconnected_at IS NULL")
	}
	if f.TokenExpiringBefore != nil {
		db = db.Where("token_expires_at IS NOT NULL AND token_expires_at <= ?", *f.TokenExpiringBefore)
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}