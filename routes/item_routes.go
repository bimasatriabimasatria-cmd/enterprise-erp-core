package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ItemRoutes(app *fiber.App) {
	// Buat grup /api/items dan pasang gembok pengaman
	api := app.Group("/api/items", middlewares.Protected())

	api.Post("/", controllers.CreateItem) // Tambah barang
	api.Get("/", controllers.GetItems)    // Lihat barang
}