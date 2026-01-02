package config

import (
	middelware "absen-qr-backend/internal/delivery/http/middleware"
	"absen-qr-backend/internal/model"
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      "AbsenQR API",
		ErrorHandler: NewErrorHandler(),
	})

	middelware.Middleware(app)

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if errors.Is(err, context.DeadlineExceeded) {
			code = fiber.StatusRequestTimeout
			message = "Process took too long"
		} else if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		} else {
			message = err.Error()
		}

		var errorResponse model.ErrorResponse

		if jsonErr := json.Unmarshal([]byte(message), &errorResponse); jsonErr != nil {
			errorResponse.Message = message
		}

		return ctx.Status(code).JSON(model.WebResponse[any]{
			Errors: &errorResponse,
		})
	}
}
