package model

import "time"

type Operator struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UID       string    `gorm:"column:uid;uniqueIndex"`
	Email     string    `gorm:"column:email;uniqueIndex"`
	Name      string    `gorm:"column:name"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (Operator) TableName() string {
	return "operators"
}