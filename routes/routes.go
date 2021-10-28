package routes

import (
	"github.com/ElayadeIsmail/go-pingram/controllers"
	"github.com/ElayadeIsmail/go-pingram/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// add currentuserMiddleware
	app.Use(middlewares.CurrentUser)

	// Group API
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World")
	})

	// Auth group
	auth := api.Group("/auth")
	// Auth Routes
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Get("/currentuser", controllers.CurrentUser)
}
