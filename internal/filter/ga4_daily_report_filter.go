package filter

import (
	"time"

	"gorm.io/gorm"
)

type GA4DailyReportFilter struct {
	ID              *uint64
	GA4ConnectionID *int64
	ReportDate      *time.Time
	ReportDateFrom  *time.Time
	ReportDateTo    *time.Time
	Limit           *int
	Offset          *int
}

func (f *GA4DailyReportFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.GA4ConnectionID != nil {
		db = db.Where("ga4_connection_id = ?", *f.GA4ConnectionID)
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