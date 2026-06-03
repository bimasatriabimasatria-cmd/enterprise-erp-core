package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WarehouseRoutes(app *fiber.App) {
	api := app.Group("/api/warehouse", middlewares.Protected())

	// TAMBAHKAN KATA "/movement" DI SINI
	api.Post("/movement", controllers.StockMovement)
	api.Get("/", controllers.GetWarehouses)
}
