package repository

import (
	"absen-qr-backend/internal/entity"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByIdForUpdate(tx *gorm.DB, user *entity.User) error
	FindAll(tx *gorm.DB, userId uint, role string, users *[]entity.User) error
	FindAllForPagging(tx *gorm.DB, userId uint, role string, pageSize int, offset int, order string, sortBy string, users *[]entity.User) (int64, error)
	FindCoachByName(tx *gorm.DB, user *entity.User) error
	FindStudentByName(tx *gorm.DB, user *entity.User) error
	FindByUsername(tx *gorm.DB, user *entity.User) error
	FindById(tx *gorm.DB, user *entity.User) error
	Login(tx *gorm.DB, user *entity.User) error
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

func (repository *UserRepositoryImpl) FindAllForPagging(tx *gorm.DB, userId uint, role string, pageSize int, offset int, order string, sortBy string, users *[]entity.User) (int64, error) {
	var total int64
	query := tx.Model(&entity.User{})

	sortColumn := "users.created_at"

	switch role {
	case "coach":
		query = query.Preload("Coach")

		if sortBy == "name" {
			query = query.Joins("JOIN coaches ON coaches.user_id = users.id")
			sortColumn = "coaches.full_name"
		}

	case "student":
		query = query.Preload("Student")

		if sortBy == "name" {
			query = query.Joins("JOIN students ON students.user_id = users.id")
			sortColumn = "students.full_name"
		}
	}

	query = query.Not("users.id = ?", userId).Where("users.role = ?", role)

	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	validOrder := "DESC"
	if order == "ASC" {
		validOrder = "ASC"
	}

	err := query.Order(sortColumn + " " + validOrder).
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	return total, err
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
func (repository *UserRepositoryImpl) Login(tx *gorm.DB, user *entity.User) error {
	err := repository.FindByUsername(tx, user)
	if err != nil {
		return err
	}

	switch user.Role {
	case "coach":
		return tx.Preload("Coach", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "full_name")
		}).First(user, user.ID).Error
	case "student":
		return tx.Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "full_name")
		}).First(user, user.ID).Error
	}

	return nil
}
