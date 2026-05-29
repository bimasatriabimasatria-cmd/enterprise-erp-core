package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole memeriksa apakah jabatan (role) pengguna diizinkan mengakses rute ini
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil jabatan dari token yang sudah dibongkar oleh Satpam Pertama (auth_middleware)
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal membaca hak akses pengguna",
			})
		}

		// Cocokkan jabatan pengguna dengan daftar jabatan yang diizinkan
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next() // Jabatan cocok, silakan masuk!
			}
		}

		// Jika perulangan selesai dan tidak ada yang cocok, tolak aksesnya
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses Ditolak (403 Forbidden)! Anda tidak memiliki jabatan yang diizinkan untuk melihat menu ini.",
		})
	}
}
