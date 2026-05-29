package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func JournalRoutes(app *fiber.App) {
	api := app.Group("/api/journals", middlewares.Protected())

	api.Post("/", controllers.CreateJournal)
	api.Get("/", controllers.GetJournals)
}
