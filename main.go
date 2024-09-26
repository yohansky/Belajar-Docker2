package main

import (
	"ambassador/src/database"
	"ambassador/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	database.Connect()
	database.AutoMigrate()
	database.SetupRedis()
	database.SetupCacheChannel()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World aku yohanes")
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000,http://localhost:4000",
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	routes.Setup(app)

	app.Listen(":8080")

}
