package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App) {
	api := app.Group("/api/reports", middlewares.Protected())

	// BARIS INI BERUBAH: Kita tambahkan RequireRole("admin", "manager")
	// Artinya staf biasa tidak akan bisa membuka Laba Rugi
	api.Get(
		"/income-statement",
		middlewares.RequireRole("admin", "manager"),
		controllers.GetIncomeStatement,
	)
}
