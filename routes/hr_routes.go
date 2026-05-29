package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func HRRoutes(app *fiber.App) {
	api := app.Group("/api/hr", middlewares.Protected())

	api.Post("/employees", controllers.CreateEmployee)
	api.Get("/employees", controllers.GetEmployees)

	api.Post("/payroll", controllers.ProcessPayroll)
}
