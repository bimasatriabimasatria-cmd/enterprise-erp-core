package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuditRoutes(app *fiber.App) {
	// Pasang gembok ganda: Harus Login (Protected) DAN Harus Admin (RequireRole)
	api := app.Group("/api/audit", middlewares.Protected(), middlewares.RequireRole("admin"))

	api.Get("/logs", controllers.GetAuditLogs)
}
