package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type StudentRepository interface {
	GetCount(tx *gorm.DB, totalStudent *int64) error
	FindAll(tx *gorm.DB, students *[]entity.Student) error
	FindById(tx *gorm.DB, student *entity.Student) error
	Create(tx *gorm.DB, student *entity.Student) error
	Update(tx *gorm.DB, student *entity.Student) error
	Delete(tx *gorm.DB, student *entity.Student) error
}

type StudentRepositoryImpl struct {
	Repository[entity.Student]
}

func NewStudentRepository() StudentRepository {
	return &StudentRepositoryImpl{}
}

// GetCount implements StudentRepository.
func (repository *StudentRepositoryImpl) GetCount(tx *gorm.DB, totalStudent *int64) error {
	return tx.Model(entity.Student{}).Count(totalStudent).Error
}

// FindAll implements StudentRepository.
func (repository *StudentRepositoryImpl) FindAll(tx *gorm.DB, students *[]entity.Student) error {
	return tx.Find(students).Error
}

// FindById implements StudentRepository.
func (repository *StudentRepositoryImpl) FindById(tx *gorm.DB, student *entity.Student) error {
	return tx.Where("user_id = ?", student.UserId).First(student).Error
}
