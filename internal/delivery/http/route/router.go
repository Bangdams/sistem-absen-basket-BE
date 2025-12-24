package route

import (
	"absen-qr-backend/internal/delivery/http"
	"absen-qr-backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                     *fiber.App
	UserController          http.UserController
	SessionController       http.SessionController
	AttendanceLogController http.AttendanceLogController
}

func (config *RouteConfig) Setup() {
	// Api for login
	config.App.Post("/login", config.UserController.Login)
	config.App.Post("/logout", config.UserController.Logout)

	// Group Api
	api := config.App.Group("/api")

	// api check status login
	api.Get("/check-login", config.UserController.CheckLogin)

	// Api for show qr
	api.Get("/show-qrcode/:filename", config.SessionController.ShowQrCode)

	// Api For Management User
	api.Get("/users", util.CheckLevel("admin"), config.UserController.FindAll)
	api.Get("/users/:id", util.CheckLevel("admin"), config.UserController.FindByIdForUpdate)
	api.Post("/users", util.CheckLevel("admin"), config.UserController.Create)
	api.Delete("/users/:id", util.CheckLevel("admin"), config.UserController.Delete)
	api.Put("/users", util.CheckLevel("admin"), config.UserController.Update)

	// Api For Management Session
	api.Get("/session", util.CheckLevel("admin", "coach"), config.SessionController.FindAll)
	api.Get("/session/:id", util.CheckLevel("admin", "coach"), config.SessionController.FindById)
	api.Post("/session", util.CheckLevel("admin", "coach"), config.SessionController.Create)
	api.Delete("/session/:id", util.CheckLevel("admin", "coach"), config.SessionController.Delete)
	api.Put("/session", util.CheckLevel("admin", "coach"), config.SessionController.Update)

	// Api For AttendanceLog
	api.Get("/attendance/student", util.CheckLevel("student"), config.AttendanceLogController.FindAllByStudent)
	api.Patch("/attendance/student", util.CheckLevel("student"), config.AttendanceLogController.UpdateByStudent)
	api.Get("/attendance/:sessionId", util.CheckLevel("coach", "admin"), config.AttendanceLogController.FindAllBySessionId)
	api.Patch("/attendance", util.CheckLevel("coach", "admin"), config.AttendanceLogController.UpdateByCoach)
}
