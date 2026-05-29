package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WarehouseRoutes(app *fiber.App) {
	api := app.Group("/api/warehouses", middlewares.Protected())

	// Master Gudang
	api.Post("/", controllers.CreateWarehouse)
	api.Get("/", controllers.GetWarehouses)

	// Inventaris dan Mutasi
	api.Post("/inventory", controllers.SetInventory)
	api.Get("/inventory", controllers.GetInventory)
	api.Post("/transfer", controllers.TransferStock)
}
