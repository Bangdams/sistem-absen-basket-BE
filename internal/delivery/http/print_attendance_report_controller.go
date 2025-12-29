package http

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PrintAttendanceReportController interface {
	GetAttendanceReport(ctx *fiber.Ctx) error
}

type PrintAttendanceReportControllerImpl struct {
	PrintAttendanceReportUsecase usecase.PrintAttendanceReportUsecase
}

func NewPrintAttendanceReportController(printAttendanceReportUsecase usecase.PrintAttendanceReportUsecase) PrintAttendanceReportController {
	return &PrintAttendanceReportControllerImpl{
		PrintAttendanceReportUsecase: printAttendanceReportUsecase,
	}
}

// GetAttendanceReport implements [PrintAttendanceReportController].
func (controller *PrintAttendanceReportControllerImpl) GetAttendanceReport(ctx *fiber.Ctx) error {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	response, err := controller.PrintAttendanceReportUsecase.GetAttendanceReport(ctx.UserContext(), startDate, endDate)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponses[model.PrintAttendanceReportResponse]{
		Data: response,
	})
}
