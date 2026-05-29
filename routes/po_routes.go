package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PORoutes(app *fiber.App) {
	api := app.Group("/api/pos", middlewares.Protected())

	api.Post("/", controllers.CreatePO)             // Buat Pesanan
	api.Get("/", controllers.GetPOs)                // Lihat Daftar Pesanan
	api.Post("/:id/receive", controllers.ReceivePO) // Tombol "Terima Barang"
}
