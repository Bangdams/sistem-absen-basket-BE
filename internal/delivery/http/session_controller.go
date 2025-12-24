package http

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/usecase"
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SessionController interface {
	ShowQrCode(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
}

type SessionControllerImpl struct {
	SessionUsecase usecase.SessionUsecase
}

func NewSessionController(sessionUsecase usecase.SessionUsecase) SessionController {
	return &SessionControllerImpl{
		SessionUsecase: sessionUsecase,
	}
}

// ShowQrCode implements SessionController.
func (controller *SessionControllerImpl) ShowQrCode(ctx *fiber.Ctx) error {
	filename := ctx.Params("filename")
	filepath := filepath.Join("./upload", filename)

	err := ctx.SendFile(filepath)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return nil
}

// FindById implements SessionController.
func (controller *SessionControllerImpl) FindById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	response, err := controller.SessionUsecase.FindById(ctx.UserContext(), uint(id))
	if err != nil {
		log.Println("failed to findbyid session")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SessionResponse]{Data: response})
}

// FindAll implements SessionController.
func (controller *SessionControllerImpl) FindAll(ctx *fiber.Ctx) error {
	order := ctx.Query("order")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")

	responses, currentPage, totalRecords, totalPages, err := controller.SessionUsecase.FindAll(ctx.UserContext(), order, page, limit)
	if err != nil {
		log.Println("failed to FindAll session")
		return err
	}

	return ctx.JSON(model.WebResponsesPagination[model.SessionResponse]{
		Data:         responses,
		CurrentPage:  *currentPage,
		TotalRecords: *totalRecords,
		TotalPages:   *totalPages,
	})
}

// Delete implements SessionController.
func (controller *SessionControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.SessionUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete session")
		return err
	}

	return nil
}

// Update implements SessionController.
func (controller *SessionControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateSessionRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.SessionUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update Session")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SessionResponse]{Data: response})
}

// Create implements SessionController.
func (controller *SessionControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.SessionRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	// get data from jwt
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	request.CoachId = uint(userId)

	response, err := controller.SessionUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create session")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SessionResponse]{Data: response})
}
