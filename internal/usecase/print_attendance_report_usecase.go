package usecase

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/model/converter"
	"absen-qr-backend/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PrintAttendanceReportUsecase interface {
	GetAttendanceReport(ctx context.Context, startDate, endDate string) (*[]model.PrintAttendanceReportResponse, error)
}

type PrintAttendanceReportUsecaseImpl struct {
	SessionRepo repository.SessionRepository
	DB          *gorm.DB
	Validate    *validator.Validate
}

func NewPrintAttendanceReport(sessionRepo repository.SessionRepository, dB *gorm.DB, validate *validator.Validate) PrintAttendanceReportUsecase {
	return &PrintAttendanceReportUsecaseImpl{
		SessionRepo: sessionRepo,
		DB:          dB,
		Validate:    validate,
	}
}

// GetAttendanceReport implements [PrintAttendanceReport].
func (printAttendanceReportUsecase *PrintAttendanceReportUsecaseImpl) GetAttendanceReport(ctx context.Context, startDate, endDate string) (*[]model.PrintAttendanceReportResponse, error) {
	sessions := []entity.Session{}

	err := printAttendanceReportUsecase.SessionRepo.PrintAttendanceReport(printAttendanceReportUsecase.DB.WithContext(ctx), startDate, endDate, &sessions)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.PrintAttendanceReportResponses(&sessions), nil
}
