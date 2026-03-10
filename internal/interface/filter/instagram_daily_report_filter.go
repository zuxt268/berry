package filter

import (
	"time"

	"gorm.io/gorm"
)

type InstagramDailyReportFilter struct {
	ID                    *uint64
	InstagramConnectionID *int64
	ReportDate            *time.Time
	ReportDateFrom        *time.Time
	ReportDateTo          *time.Time
	Limit                 *int
	Offset                *int
}

func (f *InstagramDailyReportFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.InstagramConnectionID != nil {
		db = db.Where("instagram_connection_id = ?", *f.InstagramConnectionID)
	}
	if f.ReportDate != nil {
		db = db.Where("report_date = ?", *f.ReportDate)
	}
	if f.ReportDateFrom != nil {
		db = db.Where("report_date >= ?", *f.ReportDateFrom)
	}
	if f.ReportDateTo != nil {
		db = db.Where("report_date <= ?", *f.ReportDateTo)
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}