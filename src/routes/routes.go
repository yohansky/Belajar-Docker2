package routes

import (
	"ambassador/src/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("register", controller.Register)
	admin.Post("login", controller.Login)
}
