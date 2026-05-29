package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah satpam penjaga jalur API.
// Jika tidak punya token atau tokennya salah, akses ditolak.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization dari permintaan klien
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak. Token tidak ditemukan"})
		}

		// 2. Format standar harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Format token tidak valid. Gunakan format: Bearer <token>"})
		}

		tokenString := parts[1]

		// 3. Bongkar dan validasi token menggunakan rahasia dari .env
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma enkripsi sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau sudah kadaluarsa"})
		}

		// 4. Ekstrak isi Kartu ID (Karyawan ini dari perusahaan mana? Jabatannya apa?)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Gagal membaca isi token"})
		}

		// 5. Simpan data ini di memori Fiber (Locals) agar bisa dipakai oleh Controller nanti
		c.Locals("tenant_id", claims["tenant_id"])
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		// 6. Silakan lewat! (Lanjut ke controller/proses berikutnya)
		return c.Next()
	}
}