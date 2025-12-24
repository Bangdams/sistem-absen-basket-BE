package entity

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null"`
	CreatedAt time.Time
	Student   Student `gorm:"foreignKey:UserId;references:ID"`
	Coach     Coach   `gorm:"foreignKey:UserId;references:ID"`
}
