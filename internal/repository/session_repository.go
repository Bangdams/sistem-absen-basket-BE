package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(tx *gorm.DB, session *entity.Session) error
	Update(tx *gorm.DB, session *entity.Session) error
	Delete(tx *gorm.DB, session *entity.Session) error
	FindAll(tx *gorm.DB, sessions *[]entity.Session) error
	FindById(tx *gorm.DB, session *entity.Session) error
}

type SessionRepositoryImpl struct {
	Repository[entity.Session]
}

func NewSessionRepository() SessionRepository {
	return &SessionRepositoryImpl{}
}

// FindAll implements SessionRepository.
func (repository *SessionRepositoryImpl) FindAll(tx *gorm.DB, sessions *[]entity.Session) error {
	return tx.Preload("Coach", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "full_name")
	}).Find(sessions).Error
}

// FindById implements SessionRepository.
func (repository *SessionRepositoryImpl) FindById(tx *gorm.DB, session *entity.Session) error {

	return tx.First(session).Error
}
