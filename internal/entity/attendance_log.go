package entity

import (
	"time"
)

type AttendanceLog struct {
	ID        uint `gorm:"primaryKey"`
	SessionId uint `gorm:"not null"`
	StudentId uint `gorm:"not null"`
	ScannedAt *time.Time
	Status    string  `gorm:"not null"`
	Session   Session `gorm:"foreignKey:SessionId;references:ID"`
	Student   Student `gorm:"foreignKey:StudentId;references:UserId"`
}
