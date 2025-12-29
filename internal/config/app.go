package config

import (
	"absen-qr-backend/internal/delivery/http"
	"absen-qr-backend/internal/delivery/http/route"
	"absen-qr-backend/internal/repository"
	"absen-qr-backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// repo
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()
	attendanceRepo := repository.NewAttendanceLogRepository()
	studentRepo := repository.NewStudentRepository()

	// usecase
	userUsecase := usecase.NewUserUsecase(userRepo, config.DB, config.Validate)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo, config.DB, config.Validate)
	attendanceUsecase := usecase.NewAttendanceLogUsecase(attendanceRepo, config.DB, config.Validate)
	dashboardUsecase := usecase.NewDashboardUsecase(sessionRepo, studentRepo, config.DB, config.Validate)
	printAttendanceReportUsecase := usecase.NewPrintAttendanceReport(sessionRepo, config.DB, config.Validate)

	// controller
	userController := http.NewUserController(userUsecase)
	sessionController := http.NewSessionController(sessionUsecase)
	attendanceController := http.NewAttendanceLogController(attendanceUsecase)
	dashboardController := http.NewDashboardController(dashboardUsecase)
	printAttendanceReportController := http.NewPrintAttendanceReportController(printAttendanceReportUsecase)

	routeConfig := route.RouteConfig{
		App:                             config.App,
		UserController:                  userController,
		SessionController:               sessionController,
		AttendanceLogController:         attendanceController,
		DasboardController:              dashboardController,
		PrintAttendanceReportController: printAttendanceReportController,
	}

	routeConfig.Setup()
}
