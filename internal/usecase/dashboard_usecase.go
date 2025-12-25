package usecase

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type DashboardUsecase interface {
	Dashboard(ctx context.Context) *model.DashboardResponse
}

type DashboarUsecaseImpl struct {
	SessionRepo repository.SessionRepository
	StudentRepo repository.StudentRepository
	DB          *gorm.DB
	Validate    *validator.Validate
}

func NewDashboardUsecase(sessionRepo repository.SessionRepository, studentRepo repository.StudentRepository, dB *gorm.DB, validate *validator.Validate) DashboardUsecase {
	return &DashboarUsecaseImpl{
		SessionRepo: sessionRepo,
		StudentRepo: studentRepo,
		DB:          dB,
		Validate:    validate,
	}
}

// Dashboard implements DashboardUsecase.
func (dashboardUsecase *DashboarUsecaseImpl) Dashboard(ctx context.Context) *model.DashboardResponse {
	var totalStudent int64

	dashboardUsecase.StudentRepo.GetCount(dashboardUsecase.DB.WithContext(ctx), &totalStudent)

	results := dashboardUsecase.SessionRepo.GetPresentToday(dashboardUsecase.DB.WithContext(ctx))

	data := model.DashboardResponse{
		TotalStudent: int(totalStudent),
	}

	for _, result := range results {
		if result.Status == "Hadir" {
			data.PresentToday += int(result.Total)
		} else {
			data.Absent += int(result.Total)
		}
	}

	attendancePercentage := float64(data.PresentToday) / float64(totalStudent) * 100
	data.AttendancePercentage = attendancePercentage

	return &data
}
