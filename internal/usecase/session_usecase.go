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

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// cek hela
type SessionUsecase interface {
	FindById(ctx context.Context, sessionId uint) (*model.SessionResponse, error)
	FindAll(ctx context.Context, order string, page int, limit int) (*[]model.SessionResponse, *int, *int, *int, error)
	Create(ctx context.Context, request *model.SessionRequest) (*model.SessionResponse, error)
	Delete(ctx context.Context, sessionId uint) error
	Update(ctx context.Context, request *model.UpdateSessionRequest) (*model.SessionResponse, error)
}

type SessionUsecaseImpl struct {
	SessionRepo repository.SessionRepository
	UserRepo    repository.UserRepository
	DB          *gorm.DB
	Validate    *validator.Validate
}

func NewSessionUsecase(sessionRepo repository.SessionRepository, userRepo repository.UserRepository, DB *gorm.DB, validate *validator.Validate) SessionUsecase {
	return &SessionUsecaseImpl{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		DB:          DB,
		Validate:    validate,
	}
}

// FindById implements SessionUsecase.
func (sessionUsecase *SessionUsecaseImpl) FindById(ctx context.Context, sessionId uint) (*model.SessionResponse, error) {
	session := &entity.Session{
		ID: sessionId,
	}

	err := sessionUsecase.SessionRepo.FindById(sessionUsecase.DB.WithContext(ctx), session)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Session data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error findbyid session : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error findbyid session : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SessionToResponse(session), nil
}

// FindAll implements SessionUsecase.
func (sessionUsecase *SessionUsecaseImpl) FindAll(ctx context.Context, order string, page int, limit int) (*[]model.SessionResponse, *int, *int, *int, error) {
	sessions := &[]entity.Session{}

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

	totalRecords, err := sessionUsecase.SessionRepo.FindAll(sessionUsecase.DB.WithContext(ctx), limit, offset, order, sessions)
	if err != nil {
		log.Println("failed when find all repo session : ", err)
		return nil, nil, nil, nil, fiber.ErrInternalServerError
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	taotalRecodeInt := int(totalRecords)

	log.Println("success find all from usecase sessions")

	return converter.SessionToResponses(sessions), &page, &taotalRecodeInt, &totalPages, nil
}

// Create implements SessionUsecase.
func (sessionUsecase *SessionUsecaseImpl) Create(ctx context.Context, request *model.SessionRequest) (*model.SessionResponse, error) {
	tx := sessionUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := sessionUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create session : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	session := &entity.Session{
		CoachId: request.CoachId,
		Title:   request.Title,
	}

	session.StartedAt = request.StartedAt
	session.ExpiresAt = request.ExpiresAt

	users := &[]entity.User{}

	sessionUsecase.UserRepo.FindAll(tx, request.CoachId, "student", users)

	attendanceLogs := &[]entity.AttendanceLog{}

	for _, user := range *users {
		attendanceLog := &entity.AttendanceLog{
			StudentId: user.ID,
			Status:    "Alpa",
		}

		*attendanceLogs = append(*attendanceLogs, *attendanceLog)
	}

	session.AttendanceLog = *attendanceLogs

	startedAt, err := util.ParseTimeToday(request.StartedAt)
	if err != nil {
		return nil, err
	}

	expiresAt, err := util.ParseTimeToday(request.ExpiresAt)
	if err != nil {
		return nil, err
	}

	randomChar := uuid.New()
	fileName := randomChar.String() + ".png"

	session.QrToken = fileName

	err = sessionUsecase.SessionRepo.Create(tx, session)
	if err != nil {
		log.Println("create sessions err:", err)
		return nil, fiber.ErrInternalServerError
	}

	tokenJwt, err := util.GenerateTokenAttendance(startedAt, expiresAt, session.ID)
	if err != nil {
		return nil, err
	}

	// create qr image
	if err := util.QrCodeGen(tokenJwt, fileName); err != nil {
		log.Println("generate qrCode err in session usecase:", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)

		// delete image
		if err := util.DeleteQrImage(fileName); err != nil {
			log.Println("Delete qrCode image err in session usecase:", err)
			return nil, fiber.ErrInternalServerError
		}

		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase sessions")

	return converter.SessionToResponse(session), nil
}

// Delete implements SessionUsecase.
func (sessionUsecase *SessionUsecaseImpl) Delete(ctx context.Context, sessionId uint) error {
	tx := sessionUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	session := new(entity.Session)
	session.ID = sessionId

	err := sessionUsecase.SessionRepo.FindById(tx, session)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Session data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete session : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete session : ", err)
		return fiber.ErrInternalServerError
	}

	err = sessionUsecase.SessionRepo.Delete(tx, session)
	if err != nil {
		log.Println("failed when delete repo session : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	// delete image
	if err := util.DeleteQrImage(session.QrToken); err != nil {
		log.Println("Delete qrCode image err in session usecase:", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase session")

	return nil
}

// Update implements SessionUsecase.
func (sessionUsecase *SessionUsecaseImpl) Update(ctx context.Context, request *model.UpdateSessionRequest) (*model.SessionResponse, error) {
	tx := sessionUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := sessionUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update session : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	session := &entity.Session{
		ID: request.ID,
	}

	err = sessionUsecase.SessionRepo.FindById(tx, session)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Session data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete session : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error delete session : ", err)
		return nil, fiber.ErrInternalServerError
	}

	startedAt, err := util.ParseTimeToday(request.StartedAt)
	if err != nil {
		return nil, err
	}

	expiresAt, err := util.ParseTimeToday(request.ExpiresAt)
	if err != nil {
		return nil, err
	}

	// Compare session and request times (HH:MM) to detect changes
	status := util.CompareSessionTime(session.StartedAt, request.StartedAt, session.ExpiresAt, request.ExpiresAt)

	session.Title = request.Title
	session.StartedAt = request.StartedAt
	session.ExpiresAt = request.ExpiresAt

	err = sessionUsecase.SessionRepo.Update(tx, session)
	if err != nil {
		log.Println("failed when update repo session : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	tokenJwt, err := util.GenerateTokenAttendance(startedAt, expiresAt, session.ID)
	if err != nil {
		return nil, err
	}

	if status {
		// delete qr image
		if err := util.DeleteQrImage(session.QrToken); err != nil {
			log.Println("Delete qrCode image err in session usecase:", err)
			return nil, fiber.ErrInternalServerError
		}

		// create qr image
		if err := util.QrCodeGen(tokenJwt, session.QrToken); err != nil {
			log.Println("generate qrCode err in session usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	log.Println("success create from usecase session")

	return converter.SessionToResponse(session), nil
}
