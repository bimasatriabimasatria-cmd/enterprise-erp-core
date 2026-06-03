package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func HRRoutes(app *fiber.App) {
	// KITA PAKSA AWALANNYA HARUS /api/hr
	api := app.Group("/api/hr", middlewares.Protected())

	api.Get("/employees", controllers.GetEmployees)
	api.Post("/employees", controllers.CreateEmployee)
	api.Post("/payroll", controllers.ProcessPayroll)
}
