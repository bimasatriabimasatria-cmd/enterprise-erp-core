package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WarehouseRoutes(app *fiber.App) {
	api := app.Group("/api/warehouse", middlewares.Protected())

	// Ini adalah rute yang dipanggil oleh tombol biru di React Anda!
	api.Post("/", controllers.StockMovement)

	// (Jika nanti Anda butuh rute lain, tambahkan di bawah sini)
}
