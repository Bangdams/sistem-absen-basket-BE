package entity

type Student struct {
	UserId        uint            `gorm:"primaryKey"`
	Nis           string          `gorm:"not null;unique"`
	FullName      string          `gorm:"not null"`
	Address       string          `gorm:"not null"`
	PhoneNumber   string          `gorm:"not null;unique"`
	User          *User           `gorm:"foreignKey:UserId;references:ID"`
	AttendanceLog []AttendanceLog `gorm:"foreignKey:StudentId;references:UserId"`
}
