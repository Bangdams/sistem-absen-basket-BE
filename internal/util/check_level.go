package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckLevel(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userToken := ctx.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userLevel := claims["role"].(string)

		for _, level := range roles {
			if userLevel == level {
				return ctx.Next()
			}
		}
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
}
