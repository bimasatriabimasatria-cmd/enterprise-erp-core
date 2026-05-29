package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ManufacturingRoutes(app *fiber.App) {
	api := app.Group("/api/manufacturing", middlewares.Protected())

	// Resep BOM
	api.Post("/bom", controllers.CreateBOM)

	// Perintah Produksi
	api.Post("/orders", controllers.CreateProductionOrder)
	api.Post("/orders/:id/complete", controllers.CompleteProduction) // Tombol Eksekusi
}
