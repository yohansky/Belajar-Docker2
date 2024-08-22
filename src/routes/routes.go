package routes

import (
	"ambassador/src/controller"
	"ambassador/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("register", controller.Register)
	admin.Post("login", controller.Login)

	adminAuth := admin.Use(middleware.IsAuth)

	adminAuth.Get("user", controller.User)
	adminAuth.Post("logout", controller.Logout)
	adminAuth.Put("users/info", controller.UpdateInfo)
	adminAuth.Put("users/password", controller.UpdatePassword)
	adminAuth.Put("ambassador", controller.Ambassador)
}
