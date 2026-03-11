package model

import "time"

type User struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	UID       string    `gorm:"column:uid"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Status    int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
