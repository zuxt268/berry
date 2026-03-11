package filter

import "gorm.io/gorm"

type Filter interface {
	// Apply Where句などの条件を適用
	Apply(db *gorm.DB) *gorm.DB
}
