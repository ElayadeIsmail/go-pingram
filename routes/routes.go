package routes

import (
	"github.com/ElayadeIsmail/go-pingram/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Group API
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server Runing on Port 3000")
	})

	// Auth group
	auth := app.Group("/auth")

	auth.Post("/register", controllers.Register)
}
