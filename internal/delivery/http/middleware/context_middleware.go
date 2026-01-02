package middelware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Untuk kirim context dengan timeout ke request
func TimeoutContextMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c.SetUserContext(ctx)

		return c.Next()
	}
}
