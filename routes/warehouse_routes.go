package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WarehouseRoutes(app *fiber.App) {
	api := app.Group("/api/warehouse", middlewares.Protected())

	// 1. Master Lokasi Gudang
	api.Post("/", controllers.CreateWarehouse) // <--- INI PINTU YANG KEMARIN TERHAPUS (Penyebab 405)
	api.Get("/", controllers.GetWarehouses)

	// 2. Transaksi & Mutasi Stok
	api.Post("/movement", controllers.StockMovement)

	// 3. Fitur Lanjutan (Sesuai dengan controller asli Anda)
	api.Post("/transfer", controllers.TransferStock)
	api.Post("/inventory", controllers.SetInventory)
	api.Get("/inventory", controllers.GetInventory)
}
