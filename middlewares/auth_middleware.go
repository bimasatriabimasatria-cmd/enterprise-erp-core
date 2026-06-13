package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mocking Tenant ID (Masukkan UUID perusahaan dummy Anda di sini)
		// Nanti di Fase 2, nilai ini didapat dari dekripsi JWT Token.
		c.Locals("tenant_id", "620c6b50-8e5d-4da8-8c5c-60c7c4052079")

		// Lanjutkan ke Controller (CreateJournal)
		return c.Next()
	}
}
