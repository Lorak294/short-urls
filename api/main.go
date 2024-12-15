package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	// load env variables
	err := godotenv.Load()
	if(err!= nil) {
		fmt.Println(err)
	}

	// create fiber app instance
	app := fiber.New()

	// registe middlewares
	app.Use(logger.New())

	// setup routing
	setupRoutes(app)

	// run the application
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}