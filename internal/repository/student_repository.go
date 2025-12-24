package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(tx *gorm.DB, student *entity.Student) error
	Update(tx *gorm.DB, student *entity.Student) error
	Delete(tx *gorm.DB, student *entity.Student) error
	FindAll(tx *gorm.DB, students *[]entity.Student) error
	FindById(tx *gorm.DB, student *entity.Student) error
}

type StudentRepositoryImpl struct {
	Repository[entity.Student]
}

func NewStudentRepository() StudentRepository {
	return &StudentRepositoryImpl{}
}

// FindAll implements StudentRepository.
func (repository *StudentRepositoryImpl) FindAll(tx *gorm.DB, students *[]entity.Student) error {
	return tx.Find(students).Error
}

// FindById implements StudentRepository.
func (repository *StudentRepositoryImpl) FindById(tx *gorm.DB, student *entity.Student) error {
	return tx.Where("user_id = ?", student.UserId).First(student).Error
}
