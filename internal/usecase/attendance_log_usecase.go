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
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type AttendanceLogUsecase interface {
	FindAllBySessionId(ctx context.Context, sessionId uint) (*[]model.AttendanceLogResponse, error)
	FindAllByStudent(ctx context.Context, studentId uint) (*[]model.StudentAttendanceLogResponse, error)
	FindById(ctx context.Context, attendanceLogId uint) (*model.AttendanceLogResponse, error)
	Update(ctx context.Context, request *model.AttendanceLogRequest, studentId uint) (*model.AttendanceLogResponse, error)
	ManualUpdate(ctx context.Context, request *model.UpdateAttendanceLogRequest) (*model.AttendanceLogResponse, error)
}

type AttendanceLogUsecaseImpl struct {
	AttendanceLogRepository repository.AttendanceLogRepository
	DB                      *gorm.DB
	Validate                *validator.Validate
}

func NewAttendanceLogUsecase(attendanceLogRepository repository.AttendanceLogRepository, DB *gorm.DB, validate *validator.Validate) AttendanceLogUsecase {
	return &AttendanceLogUsecaseImpl{
		AttendanceLogRepository: attendanceLogRepository,
		DB:                      DB,
		Validate:                validate,
	}
}

// FindAllBySessionId implements AttendanceLogUsecase.
func (attendanceLogUsecase *AttendanceLogUsecaseImpl) FindAllBySessionId(ctx context.Context, sessionId uint) (*[]model.AttendanceLogResponse, error) {
	attendanceLogs := &[]entity.AttendanceLog{}

	err := attendanceLogUsecase.AttendanceLogRepository.FindAllBySessionId(attendanceLogUsecase.DB.WithContext(ctx), attendanceLogs, sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error findAllBySessionId : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error findAllBySessionId : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AttendanceLogToResponses(attendanceLogs), nil
}

// FindAllByStudent implements AttendanceLogUsecase.
func (attendanceLogUsecase *AttendanceLogUsecaseImpl) FindAllByStudent(ctx context.Context, studentId uint) (*[]model.StudentAttendanceLogResponse, error) {
	attendanceLogs := &[]entity.AttendanceLog{}

	err := attendanceLogUsecase.AttendanceLogRepository.FindAllByStudent(attendanceLogUsecase.DB.WithContext(ctx), attendanceLogs, studentId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error FindAllByStudent : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error FindAllByStudent : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StudentAttendanceLogToResponses(attendanceLogs), nil
}

// FindById implements AttendanceLogUsecase.
func (attendanceLogUsecase *AttendanceLogUsecaseImpl) FindById(ctx context.Context, attendanceLogId uint) (*model.AttendanceLogResponse, error) {
	attendanceLog := &entity.AttendanceLog{
		ID: attendanceLogId,
	}

	err := attendanceLogUsecase.AttendanceLogRepository.FindById(attendanceLogUsecase.DB.WithContext(ctx), attendanceLog)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Session data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error findbyid attendance : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error findbyid attendance : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AttendanceLogToResponse(attendanceLog), nil
}

// Update implements AttendanceLogUsecase.
func (attendanceLogUsecase *AttendanceLogUsecaseImpl) Update(ctx context.Context, request *model.AttendanceLogRequest, studentId uint) (*model.AttendanceLogResponse, error) {
	tx := attendanceLogUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := attendanceLogUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update attendanceLog : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	token, err := util.ParseToken(request.QrToken, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			errorResponse.Message = "Your QR code has expired"
			errorResponse.Details = []string{"Please request a new QR code to continue."}
			jsonString, _ := json.Marshal(errorResponse)

			return nil, fiber.NewError(fiber.StatusUnauthorized, string(jsonString))
		}

		return nil, fiber.ErrUnauthorized
	}

	sessionId := token["session_id"].(float64)

	now := time.Now()

	attendanceLog := &entity.AttendanceLog{
		SessionId: uint(sessionId),
		StudentId: studentId,
		ScannedAt: &now,
		Status:    "Hadir",
	}

	alreadyAttended, err := attendanceLogUsecase.AttendanceLogRepository.IsUserPresent(tx, attendanceLog.SessionId, attendanceLog.StudentId)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if alreadyAttended {
		return nil, fiber.ErrConflict
	}

	err = attendanceLogUsecase.AttendanceLogRepository.Update(tx, attendanceLog)
	if err != nil {
		log.Println("update attendanceLog err:", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)

		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase attendance")

	return converter.AttendanceLogToResponse(attendanceLog), nil
}

// ManualUpdate implements AttendanceLogUsecase.
func (attendanceLogUsecase *AttendanceLogUsecaseImpl) ManualUpdate(ctx context.Context, request *model.UpdateAttendanceLogRequest) (*model.AttendanceLogResponse, error) {
	tx := attendanceLogUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	// Change to Title
	c := cases.Lower(language.Indonesian)
	request.Status = c.String(request.Status)

	err := attendanceLogUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update attendanceLog : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	now := time.Now()

	attendanceLog := &entity.AttendanceLog{
		SessionId: request.SessionId,
		StudentId: request.StudentId,
		ScannedAt: &now,
		Status:    request.Status,
	}

	err = attendanceLogUsecase.AttendanceLogRepository.ExistsBySessionAndStudent(tx, attendanceLog.SessionId, attendanceLog.StudentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("data was not found")
			return nil, fiber.ErrNotFound
		}
		return nil, fiber.ErrInternalServerError
	}

	err = attendanceLogUsecase.AttendanceLogRepository.Update(tx, attendanceLog)
	if err != nil {
		log.Println("update attendanceLog err:", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)

		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase attendance")

	return converter.AttendanceLogToResponse(attendanceLog), nil
}
