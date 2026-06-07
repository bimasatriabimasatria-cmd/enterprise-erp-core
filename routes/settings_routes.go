package routes

import (
	"enterprise-erp/controllers"

	"github.com/gofiber/fiber/v2"
)

func SettingsRoutes(app *fiber.App) {
	// Debug sederhana: Apakah rute ini benar-benar terdaftar?
	app.Get("/api/settings", controllers.GetSettings)
	app.Post("/api/settings", controllers.UpdateSettings)
}
