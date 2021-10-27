package main

import (
	"log"

	"github.com/ElayadeIsmail/go-pingram/database"
	"github.com/ElayadeIsmail/go-pingram/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// create new APP
	app := fiber.New()

	// Apply Cors Middleware
	app.Use(cors.New())

	// Connect To Database
	database.Connect()

	// Setup Routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

	// Close DB after Exit programme
	defer database.Close()
}
