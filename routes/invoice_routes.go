package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InvoiceRoutes(app *fiber.App) {
	// Grup /api/invoices diamankan dengan Token JWT
	api := app.Group("/api/invoices", middlewares.Protected())

	api.Post("/", controllers.CreateInvoice)
	api.Get("/", controllers.GetInvoices)
}
