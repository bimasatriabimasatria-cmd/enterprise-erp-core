package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AccountRoutes(app *fiber.App) {
	// Buat grup /api/accounts dan pasang pelindung token
	api := app.Group("/api/accounts", middlewares.Protected())

	api.Post("/", controllers.CreateAccount) // Tambah Akun
	api.Get("/", controllers.GetAccounts)    // Lihat Daftar Akun
}
