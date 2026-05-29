package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PortalRoutes(app *fiber.App) {
	// Pintu masuk /api/portal khusus untuk pihak eksternal
	api := app.Group("/api/portal", middlewares.Protected())

	// Rute Pelanggan: WAJIB memiliki peran (role) 'customer'
	customerPortal := api.Group("/customer", middlewares.RequireRole("customer"))
	customerPortal.Get("/invoices", controllers.GetMyInvoices)
}
