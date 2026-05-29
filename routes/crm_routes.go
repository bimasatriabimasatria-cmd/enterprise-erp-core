package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CRMRoutes(app *fiber.App) {
	api := app.Group("/api/crm", middlewares.Protected())

	// Jalur Lead
	api.Post("/leads", controllers.CreateLead)
	api.Get("/leads", controllers.GetLeads)
	api.Post("/leads/:id/convert", controllers.ConvertLead) // Tombol Konversi

	// Jalur Customer
	api.Get("/customers", controllers.GetCustomers)
}
