package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WorkflowRoutes(app *fiber.App) {
	// Semua rute di bawah ini wajib login
	api := app.Group("/api/workflow", middlewares.Protected())

	// Semua karyawan (staff/manager/admin) bisa meminta persetujuan
	api.Post("/request", controllers.RequestApproval)

	// HANYA Admin dan Manager yang bisa melihat antrean dan menyetujui
	api.Get("/pending", middlewares.RequireRole("admin", "manager"), controllers.GetPendingApprovals)
	api.Post("/:id/process", middlewares.RequireRole("admin", "manager"), controllers.ProcessApproval)
}
