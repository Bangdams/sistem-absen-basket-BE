package http

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type DasboardController interface {
	ShowDashboardData(ctx *fiber.Ctx) error
}

type DasboardControllerImpl struct {
	DashboardUsecase usecase.DashboardUsecase
}

func NewDashboardController(dashboardUsecase usecase.DashboardUsecase) DasboardController {
	return &DasboardControllerImpl{
		DashboardUsecase: dashboardUsecase,
	}
}

// ShowDashboardData implements DasboardController.
func (controller *DasboardControllerImpl) ShowDashboardData(ctx *fiber.Ctx) error {

	response := controller.DashboardUsecase.Dashboard(ctx.UserContext())

	return ctx.JSON(model.WebResponse[model.DashboardResponse]{Data: *response})
}
