package http

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/usecase"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AttendanceLogController interface {
	FindAllByStudent(ctx *fiber.Ctx) error
	FindAllBySessionId(ctx *fiber.Ctx) error
	UpdateByStudent(ctx *fiber.Ctx) error
	UpdateByCoach(ctx *fiber.Ctx) error
}

type AttendanceLogControllerImpl struct {
	AttendanceLogUsecase usecase.AttendanceLogUsecase
}

func NewAttendanceLogController(sessionUsecase usecase.AttendanceLogUsecase) AttendanceLogController {
	return &AttendanceLogControllerImpl{
		AttendanceLogUsecase: sessionUsecase,
	}
}

// FindAllByStudent implements AttendanceLogController.
func (controller *AttendanceLogControllerImpl) FindAllByStudent(ctx *fiber.Ctx) error {
	// get data from jwt
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	responses, err := controller.AttendanceLogUsecase.FindAllByStudent(ctx.UserContext(), uint(userId))
	if err != nil {
		log.Println("failed to FindAllByStudent in attendance")
		return err
	}

	return ctx.JSON(model.WebResponses[model.StudentAttendanceLogResponse]{Data: responses})
}

// FindAllBySessionId implements AttendanceLogController.
func (controller *AttendanceLogControllerImpl) FindAllBySessionId(ctx *fiber.Ctx) error {
	sessionId, err := ctx.ParamsInt("sessionId")
	if err != nil {
		return fiber.ErrBadRequest
	}

	responses, err := controller.AttendanceLogUsecase.FindAllBySessionId(ctx.UserContext(), uint(sessionId))
	if err != nil {
		log.Println("failed to FindAllBySessionId in attendance")
		return err
	}

	return ctx.JSON(model.WebResponses[model.AttendanceLogResponse]{Data: responses})
}

// UpdateByStudent implements AttendanceLogController.
func (controller *AttendanceLogControllerImpl) UpdateByStudent(ctx *fiber.Ctx) error {
	request := new(model.AttendanceLogRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	// get data from jwt token
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	idClaim, ok := claims["user_id"].(float64)
	if !ok {
		return fiber.ErrUnauthorized
	}
	studentId := uint(idClaim)

	_, err := controller.AttendanceLogUsecase.Update(ctx.UserContext(), request, uint(studentId))
	if err != nil {
		log.Println("failed to update attendance")
		return err
	}

	return nil
}

// UpdateByStudent implements AttendanceLogController.
func (controller *AttendanceLogControllerImpl) UpdateByCoach(ctx *fiber.Ctx) error {
	request := new(model.UpdateAttendanceLogRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	_, err := controller.AttendanceLogUsecase.ManualUpdate(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to manualUpdate attendance")
		return err
	}

	return nil
}
