package routes

import (
	"enterprise-erp/controllers"
	"enterprise-erp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AttendanceRoutes(app *fiber.App) {
	// PERBAIKAN: Gunakan awalan /api/attendance dan lindungi dengan satpam Token
	api := app.Group("/api/attendance", middlewares.Protected())

	api.Post("/clock-in", controllers.ClockIn)
	api.Post("/clock-out", controllers.ClockOut)

	// Data absensi hanya boleh dilihat oleh admin/manajer HRD
	api.Get("/", middlewares.RequireRole("admin", "manager"), controllers.GetAttendances)
}
