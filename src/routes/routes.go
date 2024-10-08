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
	// Admin ROutes
	admin := api.Group("admin")
	admin.Post("register", controller.Register)
	admin.Post("login", controller.Login)

	adminAuth := admin.Use(middleware.IsAuth)

	adminAuth.Get("user", controller.User)
	adminAuth.Post("logout", controller.Logout)
	adminAuth.Put("users/info", controller.UpdateInfo)
	adminAuth.Put("users/password", controller.UpdatePassword)
	adminAuth.Get("ambassador", controller.Ambassador)
	adminAuth.Get("admin", controller.Admin)
	// Product API
	adminAuth.Get("products", controller.AllProducts)
	adminAuth.Get("product/:id", controller.GetProduct)
	adminAuth.Post("products", controller.CreateProduct)
	adminAuth.Put("product/:id", controller.UpdateProduct)
	adminAuth.Delete("product/:id", controller.DeleteProduct)
	// Link API
	adminAuth.Get("users/:id/Links", controller.Links)
	// Order API
	adminAuth.Get("orders", controller.Orders)

	//ambassador Routes (seller)
	ambass := api.Group("ambassador")
	ambass.Post("register", controller.Register)
	ambass.Post("login", controller.Login)
	ambass.Get("products/frontend", controller.ProductFrontend)
	ambass.Get("products/backend", controller.ProductBackend)
	//
	ambassAuth := ambass.Use(middleware.IsAuth)
	ambassAuth.Get("user", controller.User)
	ambassAuth.Post("logout", controller.Logout)
	ambassAuth.Put("users/info", controller.UpdateInfo)
	ambassAuth.Put("users/password", controller.UpdatePassword)
	// Products API
	ambassAuth.Post("links", controller.CreateLink)
	ambassAuth.Get("stats", controller.Stats)
	ambassAuth.Get("rankings", controller.Rankings)

	// Checkout Routes
	check := api.Group("checkout")
	check.Get("links/:code", controller.GetLink)
	check.Post("orders", controller.CreateOrder)
	check.Post("orders/confirm", controller.CompleteOrder)

}
