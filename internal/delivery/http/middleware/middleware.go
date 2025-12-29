package middelware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Middelware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("SERVER_Allow_Origins"),
		AllowCredentials: true,
		AllowHeaders:     "Content-Type",
	}))

	app.Use("/api", jwtware.New(jwtware.Config{
		TokenLookup: "cookie:token",
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("SECRET_KEY")),
		},
		ContextKey: "user",
	}))
}
