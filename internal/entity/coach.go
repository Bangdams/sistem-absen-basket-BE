package entity

type Coach struct {
	UserId   uint      `gorm:"primaryKey"`
	Nip      string    `gorm:"not null;unique"`
	FullName string    `gorm:"not null"`
	User     *User     `gorm:"foreignKey:UserId;references:ID"`
	Session  []Session `gorm:"foreignKey:CoachId;references:UserId"`
}
