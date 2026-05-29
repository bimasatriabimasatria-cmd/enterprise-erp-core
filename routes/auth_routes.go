package routes

import (
	"enterprise-erp/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	// Buat grup URL khusus API
	api := app.Group("/api/auth")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
}