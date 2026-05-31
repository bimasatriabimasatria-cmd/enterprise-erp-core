package routes

import (
	"enterprise-erp/controllers"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(router fiber.Router) {
	router.Post("/transactions", controllers.CreateTransaction)
	// Anda bisa menambah rute transaksi lain di sini nanti (cth: GET /transactions)
}
