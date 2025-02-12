package main

import (
	"log"

	"github.com/WaveCE29/custom-film-cleaner/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
