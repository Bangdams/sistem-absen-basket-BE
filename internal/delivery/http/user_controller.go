package http

import (
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/usecase"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController interface {
	FindByIdForUpdate(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	CheckLogin(ctx *fiber.Ctx) error
}

type UserControllerImpl struct {
	UserUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
	return &UserControllerImpl{
		UserUsecase: userUsecase,
	}
}

// FindByIdForUpdate implements UserController.
func (controller *UserControllerImpl) FindByIdForUpdate(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	role := ctx.Query("role")

	response, err := controller.UserUsecase.FindByIdForUpdate(ctx.UserContext(), uint(id), role)
	if err != nil {
		log.Println("failed to delete session")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Create implements UserController.
func (controller *UserControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.UserRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Update implements UserController.
func (controller *UserControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("error badrequest:", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Delete implements UserController.
func (controller UserControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.UserUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete user")
		return err
	}

	return nil
}

// FindAll implements UserController.
func (controller *UserControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.UserResponse
	var err error

	// get data from jwt token
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	role := ctx.Query("role")
	sortBy := ctx.Query("sort-by")

	order := ctx.Query("order")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")

	responses, currentPage, totalRecords, totalPages, err := controller.UserUsecase.FindAll(ctx.UserContext(), uint(userId), role, order, page, limit, sortBy)
	if err != nil {
		log.Println("failed to FindAll user")
		return err
	}

	return ctx.JSON(model.WebResponsesPagination[model.UserResponse]{
		Data:         responses,
		CurrentPage:  *currentPage,
		TotalRecords: *totalRecords,
		TotalPages:   *totalPages,
	})
}

// Login implements UserController.
func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, token, err := controller.UserUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to login")
		return err
	}

	// durasi token
	duration := os.Getenv("DURATION_JWT_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    *token,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * lifeTime,
	})

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

// Logout implements UserController.
func (controller *UserControllerImpl) Logout(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("token")
	if cookie == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return ctx.JSON(model.WebResponse[string]{Data: "Logout successful"})
}

// CheckLogin implements UserController.
func (controller *UserControllerImpl) CheckLogin(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userId := claims["user_id"].(float64)
	username := claims["username"].(string)
	role := claims["role"].(string)

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: &model.LoginResponse{
		UserID:   uint(userId),
		Username: username,
		Role:     role,
	}})
}
