package usecase

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/model/converter"
	"absen-qr-backend/internal/repository"
	"absen-qr-backend/internal/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	FindByIdForUpdate(ctx context.Context, userId uint, role string) (*model.UserResponse, error)
	Create(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, userId uint) error
	FindAll(ctx context.Context, userId uint, role string, order string, page int, limit int) (*[]model.UserResponse, *int, *int, *int, error)
	Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, *string, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
}

type UserUsecaseImpl struct {
	UserRepo repository.UserRepository
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserUsecase(userRepo repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		UserRepo: userRepo,
		DB:       DB,
		Validate: validate,
	}
}

// FindById implements UserUsecase.
func (userUsecase *UserUsecaseImpl) FindByIdForUpdate(ctx context.Context, userId uint, role string) (*model.UserResponse, error) {
	user := &entity.User{
		ID:   userId,
		Role: role,
	}

	err := userUsecase.UserRepo.FindByIdForUpdate(userUsecase.DB.WithContext(ctx), user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "User data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error findbyid user : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error findbyid user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	fmt.Println(user.Coach)

	return converter.UserToResponse(user), nil
}

// Create implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Create(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create user : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate password")
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Username: request.Username,
		Password: string(password),
		Role:     request.Role,
	}

	switch user.Role {
	case "coach":
		user.Coach = entity.Coach{
			Nip:      request.CoachRequest.Nip,
			FullName: request.CoachRequest.FullName,
		}
	case "student":
		user.Student = entity.Student{
			Nis:         request.StudentRequest.Nis,
			FullName:    request.StudentRequest.FullName,
			Address:     request.StudentRequest.Address,
			PhoneNumber: request.StudentRequest.PhoneNumber,
		}
	}

	if err := userUsecase.UserRepo.FindByUsername(tx, user); err == nil {
		errorResponse.Message = "Duplicate entry"
		errorResponse.Details = []string{"Username already exists in the database."}

		jsonString, _ := json.Marshal(errorResponse)

		return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
	}

	err = userUsecase.UserRepo.Create(tx, user)
	if err != nil {
		log.Println("failed when create repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase user")

	return converter.UserToResponse(user), nil
}

// Delete implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Delete(ctx context.Context, userId uint) error {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{}
	user.ID = userId

	err := userUsecase.UserRepo.FindById(tx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "User data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete user : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete user : ", err)
		return fiber.ErrInternalServerError
	}

	err = userUsecase.UserRepo.Delete(tx, user)
	if err != nil {
		log.Println("failed when delete repo user : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase user")

	return nil
}

// FindAll implements UserUsecase.
func (userUsecase *UserUsecaseImpl) FindAll(ctx context.Context, userId uint, role string, order string, page int, limit int) (*[]model.UserResponse, *int, *int, *int, error) {
	var users = &[]entity.User{}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	if order == "" {
		order = "DESC"
	}

	totalRecords, err := userUsecase.UserRepo.FindAllForPagging(userUsecase.DB.WithContext(ctx), userId, role, limit, offset, order, users)
	if err != nil {
		log.Println("failed when find all repo user : ", err)
		return nil, nil, nil, nil, fiber.ErrInternalServerError
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	taotalRecodeInt := int(totalRecords)

	log.Println("success find all from usecase user")

	return converter.UserToResponses(users), &page, &taotalRecodeInt, &totalPages, nil
}

// Login implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, *string, error) {
	user := &entity.User{}

	if err := userUsecase.UserRepo.Login(userUsecase.DB.WithContext(ctx), user, request.Username); err != nil {
		log.Println("Login failed, username not found:", request.Username)
		return nil, nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println("Login failed, invalid password for:", request.Username)
		return nil, nil, fiber.ErrUnauthorized
	}

	token, err := util.GenerateTokenLogin(user)
	if err != nil {
		log.Printf("Failed to generate token for user %s: %v\n", request.Username, err)
		return nil, nil, fiber.ErrInternalServerError
	}

	log.Println("Success login:", request.Username)

	return converter.LoginUserToResponse(user), &token, nil
}

// Update implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update user : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	user := &entity.User{
		ID: request.ID,
	}

	err = userUsecase.UserRepo.FindByIdForUpdate(tx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "User data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			log.Println("Data not found")

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by user id : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("failed to generate password")
			return nil, fiber.ErrInternalServerError
		}

		user.Password = string(password)
	}

	user.Username = request.Username

	switch user.Role {
	case "coach":
		user.Coach = entity.Coach{
			Nip:      request.CoachRequest.Nip,
			FullName: request.CoachRequest.FullName,
		}
	case "student":
		user.Student = entity.Student{
			Nis:         request.StudentRequest.Nis,
			FullName:    request.StudentRequest.FullName,
			Address:     request.StudentRequest.Address,
			PhoneNumber: request.StudentRequest.PhoneNumber,
		}
	}

	err = userUsecase.UserRepo.Update(tx, user)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when update repo user : ", err)

		var errorField string
		parts := strings.Split(mysqlErr.Message, "'")
		if len(parts) > 2 {
			errorField = parts[1]
		}

		if mysqlErr.Number == 1062 {
			errorResponse.Message = "Duplicate entry"
			errorResponse.Details = []string{errorField + " already exists in the database."}

			jsonString, _ := json.Marshal(errorResponse)

			return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
		}

		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase user")

	return converter.UserToResponse(user), nil
}
