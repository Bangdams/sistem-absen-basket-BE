package middelware

import (
	"log"
	"os"
	"strconv"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Middleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("SERVER_Allow_Origins"),
		AllowCredentials: true,
		AllowHeaders:     "Content-Type",
	}))

	durationString := os.Getenv("TIMEOUT_CONTEXT")
	durationInt, err := strconv.Atoi(durationString)
	if err != nil {
		log.Println("Warning: TIMEOUT_CONTEXT not found or invalid, using default 10s")
		durationInt = 10
	}
	app.Use(TimeoutContextMiddleware(time.Duration(durationInt) * time.Second))

	app.Use("/api", jwtware.New(jwtware.Config{
		TokenLookup: "cookie:token",
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("SECRET_KEY")),
		},
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		},
	}))
}
