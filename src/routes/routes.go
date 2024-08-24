package routes

import (
	"ambassador/src/controller"
	"ambassador/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("products", controller.AllProducts)
	app.Get("product/:id", controller.GetProduct)

	api := app.Group("api")
	admin := api.Group("admin")
	admin.Post("register", controller.Register)
	admin.Post("login", controller.Login)

	adminAuth := admin.Use(middleware.IsAuth)

	adminAuth.Get("user", controller.User)
	adminAuth.Post("logout", controller.Logout)
	adminAuth.Put("users/info", controller.UpdateInfo)
	adminAuth.Put("users/password", controller.UpdatePassword)
	adminAuth.Get("ambassador", controller.Ambassador)
	// Product API
	adminAuth.Get("products", controller.AllProducts)
	adminAuth.Post("products", controller.CreateProduct)
	adminAuth.Put("product/:id", controller.UpdateProduct)
	adminAuth.Delete("product/:id", controller.DeleteProduct)
	// Link API
	adminAuth.Get("users/:id/Links", controller.Links)
	// Order API
	adminAuth.Get("orders", controller.Orders)

}
