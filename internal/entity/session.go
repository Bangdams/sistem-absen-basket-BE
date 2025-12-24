package entity

import (
	"time"
)

type Session struct {
	ID            uint   `gorm:"primaryKey"`
	CoachId       uint   `gorm:"not null"`
	Title         string `gorm:"not null"`
	QrToken       string `gorm:"not null;unique"`
	StartedAt     string `gorm:"type:TIME;not null" validate:"required,datetime=15:04"`
	ExpiresAt     string `gorm:"type:TIME;not null" validate:"required,datetime=15:04"`
	CreatedAt     time.Time
	Coach         Coach           `gorm:"foreignKey:CoachId;references:UserId"`
	AttendanceLog []AttendanceLog `gorm:"foreignKey:SessionId;references:ID"`
}
