package routes

import (
	"enterprise-erp/controllers" // Sesuaikan dengan path project Anda
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SettingsRoutes(app *fiber.App) {
	api := app.Group("/api/settings")

	// Siapapun bisa baca (untuk tampilan login & sidebar)
	api.Get("/", controllers.GetSettings)

	// HANYA Super Admin yang bisa edit
	api.Post("/", middlewares.AuthMiddleware(), controllers.UpdateSettings)
}
