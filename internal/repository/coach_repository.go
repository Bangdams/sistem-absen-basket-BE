package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type CoachRepository interface {
	Create(tx *gorm.DB, coach *entity.Coach) error
	Update(tx *gorm.DB, coach *entity.Coach) error
	Delete(tx *gorm.DB, coach *entity.Coach) error
	FindAll(tx *gorm.DB, coachs *[]entity.Coach) error
	FindById(tx *gorm.DB, coach *entity.Coach) error
}

type CoachRepositoryImpl struct {
	Repository[entity.Coach]
}

func NewCoachRepository() CoachRepository {
	return &CoachRepositoryImpl{}
}

// FindAll implements CoachRepository.
func (repository *CoachRepositoryImpl) FindAll(tx *gorm.DB, coachs *[]entity.Coach) error {
	return tx.Find(coachs).Error
}

// FindById implements CoachRepository.
func (repository *CoachRepositoryImpl) FindById(tx *gorm.DB, coach *entity.Coach) error {
	return tx.Where("user_id = ?", coach.UserId).First(coach).Error
}
