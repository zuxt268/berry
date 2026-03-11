package filter

import (
	"time"

	"gorm.io/gorm"
)

type LineDailyReportFilter struct {
	ID               *uint64
	LineConnectionID *int64
	ReportDate       *time.Time
	ReportDateFrom   *time.Time
	ReportDateTo     *time.Time
	Limit            *int
	Offset           *int
}

func (f *LineDailyReportFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.LineConnectionID != nil {
		db = db.Where("line_connection_id = ?", *f.LineConnectionID)
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