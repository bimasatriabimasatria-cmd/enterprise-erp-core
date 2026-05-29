package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FinanceRoutes(app *fiber.App) {
	// Buat grup khusus keuangan tingkat lanjut
	api := app.Group("/api/finance", middlewares.Protected(), middlewares.RequireRole("admin", "manager"))

	api.Post("/reconciliation/upload", controllers.UploadBankStatement)
	api.Get("/reconciliation", controllers.GetBankStatements)
	api.Post("/reconciliation/:id/auto-match", controllers.AutoReconcile)
}
