package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// Protected mensimulasikan login dan menyuntikkan konteks keamanan
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// MOCKING FASE 1:
		// Kita paksa sistem mengira bahwa request ini dari PT ERP Maju Bersama
		// Menggunakan UUID yang Anda dapatkan sebelumnya.
		// Di Fiber, kita menggunakan c.Locals() untuk menyimpan data per-request.
		c.Locals("tenant_id", "620c6b50-8e5d-4da8-8c5c-60c7c4052079")
		
		// Lanjut ke request berikutnya (Controller)
		return c.Next()
	}
}
