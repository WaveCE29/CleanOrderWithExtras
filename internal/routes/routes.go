package routes

import (
	"github.com/WaveCE29/custom-film-cleaner/internal/models"
	"github.com/WaveCE29/custom-film-cleaner/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/clean", func(c *fiber.Ctx) error {
		var input []models.InputOrder // Expect an array of orders
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		var cleanedOrders []models.CleanedOrder
		processedOrders, err := services.CleanOrder(input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		cleanedOrders = append(cleanedOrders, processedOrders...)

		return c.JSON(cleanedOrders)
	})
}
