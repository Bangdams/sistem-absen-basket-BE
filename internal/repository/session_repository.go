package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(tx *gorm.DB, session *entity.Session) error
	Update(tx *gorm.DB, session *entity.Session) error
	Delete(tx *gorm.DB, session *entity.Session) error
	FindAll(tx *gorm.DB, pageSize int, offset int, order string, sessions *[]entity.Session) (int64, error)
	FindById(tx *gorm.DB, session *entity.Session) error
}

type SessionRepositoryImpl struct {
	Repository[entity.Session]
}

func NewSessionRepository() SessionRepository {
	return &SessionRepositoryImpl{}
}

// FindAll implements SessionRepository.
// func (repository *SessionRepositoryImpl) FindAll(tx *gorm.DB, sessions *[]entity.Session) error {
// 	return tx.Preload("Coach", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("user_id", "full_name")
// 	}).Find(sessions).Error
// }

// FindAll implements SessionRepository.
func (repository *SessionRepositoryImpl) FindAll(tx *gorm.DB, pageSize int, offset int, order string, sessions *[]entity.Session) (int64, error) {
	var total int64

	query := tx.Model(&entity.Session{})

	query = query.Preload("Coach", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "full_name")
	})

	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	validOrder := "DESC"
	if order == "ASC" {
		validOrder = "ASC"
	}

	err := query.Order("created_at " + validOrder).
		Limit(pageSize).
		Offset(offset).
		Find(sessions).Error

	return total, err
}

// FindById implements SessionRepository.
func (repository *SessionRepositoryImpl) FindById(tx *gorm.DB, session *entity.Session) error {

	return tx.First(session).Error
}
