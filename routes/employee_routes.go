package routes

import (
	"enterprise-erp/controllers"

	"github.com/gofiber/fiber/v2"
)

func EmployeeRoutes(router fiber.Router) {
	employeeGroup := router.Group("/employees")

	// GET /api/employees
	employeeGroup.Get("/", controllers.GetEmployees)

	// POST /api/employees
	employeeGroup.Post("/", controllers.CreateEmployee)
}
