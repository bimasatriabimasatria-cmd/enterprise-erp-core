package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SettingsRoutes(app *fiber.App) {
	// Tanpa app.Group, langsung definisikan rute lengkap
	app.Get("/api/settings", controllers.GetSettings)
	app.Post("/api/settings", middlewares.AuthMiddleware(), controllers.UpdateSettings)
}
