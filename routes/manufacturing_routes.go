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
	api.Get("/bom", controllers.GetBOMs) // <--- TAMBAHAN MATA BOM
	api.Delete("/bom/:id", controllers.DeleteBOM)

	// Perintah Produksi
	api.Post("/orders", controllers.CreateProductionOrder)
	api.Get("/orders", controllers.GetProductionOrders) // <--- TAMBAHAN MATA PRODUKSI
	api.Post("/orders/:id/complete", controllers.CompleteProduction)
}
