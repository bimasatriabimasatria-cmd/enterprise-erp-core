package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"
	"github.com/gofiber/fiber/v2"
)

func JournalRoutes(app *fiber.App) {
	// Buat grup /api/v1/journals
	journalGroup := app.Group("/api/v1/journals")

	// Pasang CCTV (Auth Middleware) DI SINI. 
	// Setiap request ke /api/v1/journals/... akan dicegat oleh middleware ini dulu.
	journalGroup.Use(middlewares.AuthMiddleware())

	// Endpoint untuk membuat jurnal baru
	journalGroup.Post("/", controllers.CreateJournal)
	journalGroup.Post("", controllers.CreateJournal)
}
