package repository

import (
	"absen-qr-backend/internal/entity"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByIdForUpdate(tx *gorm.DB, user *entity.User) error
	FindAll(tx *gorm.DB, userId uint, role string, users *[]entity.User) error
	FindCoachByName(tx *gorm.DB, user *entity.User) error
	FindStudentByName(tx *gorm.DB, user *entity.User) error
	FindByUsername(tx *gorm.DB, user *entity.User) error
	FindById(tx *gorm.DB, user *entity.User) error
	Login(tx *gorm.DB, user *entity.User, keyword string) error
	Create(tx *gorm.DB, user *entity.User) error
	Update(tx *gorm.DB, user *entity.User) error
	Delete(tx *gorm.DB, user *entity.User) error
}

type UserRepositoryImpl struct {
	Repository[entity.User]
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// FindByIdForUpdate implements UserRepository.
func (repository *UserRepositoryImpl) FindByIdForUpdate(tx *gorm.DB, user *entity.User) error {
	query := tx.Model(&entity.User{})

	switch user.Role {
	case "coach":
		query = query.Preload("Coach")
	case "student":
		query = query.Preload("Student")
	}

	return query.First(&user).Error
}

// FindAll implements UserRepository.
func (repository *UserRepositoryImpl) FindAll(tx *gorm.DB, userId uint, role string, users *[]entity.User) error {
	query := tx.Model(&entity.User{})

	switch role {
	case "coach":
		query = query.Preload("Coach")
	case "student":
		log.Println("ada")
		query = query.Preload("Student")
	}

	return query.Not("id = ?", userId).Where("role = ?", role).Find(users).Error
}

// FindByName implements UserRepository.
func (repository *UserRepositoryImpl) FindCoachByName(tx *gorm.DB, user *entity.User) error {
	return tx.First(user, "full_name=?", user.Coach.FullName).Error
}

// FindByName implements UserRepository.
func (repository *UserRepositoryImpl) FindStudentByName(tx *gorm.DB, user *entity.User) error {
	return tx.First(user, "full_name=?", user.Student.FullName).Error
}

// FindByUsername implements UserRepository.
func (repository *UserRepositoryImpl) FindByUsername(tx *gorm.DB, user *entity.User) error {
	return tx.First(user, "username=?", user.Username).Error
}

// FindById implements UserRepository.
func (repository *UserRepositoryImpl) FindById(tx *gorm.DB, user *entity.User) error {
	return tx.First(user).Error
}

// Login implements UserRepository.
func (repository *UserRepositoryImpl) Login(tx *gorm.DB, user *entity.User, keyword string) error {
	return tx.Where("username = ?", keyword).First(user).Error
}
