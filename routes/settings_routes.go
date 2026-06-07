package routes

import (
	"enterprise-erp/controllers" // Sesuaikan dengan path project Anda
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SettingsRoutes(app *fiber.App) {
	// Gunakan app.Get langsung jika ingin menghindari masalah trailing slash
	app.Get("/api/settings", controllers.GetSettings)
	app.Post("/api/settings", middlewares.AuthMiddleware(), controllers.UpdateSettings)
}
