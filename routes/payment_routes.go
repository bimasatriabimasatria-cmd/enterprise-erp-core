package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PaymentRoutes(app *fiber.App) {
	api := app.Group("/api/payments", middlewares.Protected())

	api.Post("/receive", controllers.ReceivePayment)
}
