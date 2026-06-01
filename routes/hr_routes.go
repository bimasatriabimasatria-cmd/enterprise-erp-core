package routes

import (
	"enterprise-erp/controllers"

	"github.com/gofiber/fiber/v2"
)

func HRRoutes(router fiber.Router) {
	// Pintu masuknya kita ubah menjadi /hr/employees
	hrGroup := router.Group("/hr/employees")

	hrGroup.Get("/", controllers.GetEmployees)
	hrGroup.Post("/", controllers.CreateEmployee)
}
