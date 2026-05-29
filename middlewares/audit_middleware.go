package middlewares

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// AuditLogger mencatat setiap kali pengguna menambah/mengubah/menghapus data
func AuditLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Biarkan sistem melayani permintaan pengguna terlebih dahulu
		err := c.Next()

		// 2. Baca jenis aksinya
		method := c.Method()

		// Kita hanya peduli jika ada perubahan data (Tambah, Ubah, Hapus)
		// Kita tidak perlu mencatat aktivitas cuma "Melihat-lihat" (GET) agar database tidak cepat penuh
		if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {

			// Jangan catat password saat registrasi dan login
			if c.Path() == "/api/auth/login" || c.Path() == "/api/auth/register" {
				return err
			}

			// 3. Ambil data identitas dari Token JWT (Siapa pelakunya?)
			tenantID, _ := c.Locals("tenant_id").(string)
			userID, _ := c.Locals("user_id").(string)

			// Pastikan pelakunya sudah login
			if tenantID != "" && userID != "" {
				log := models.AuditLog{
					TenantID:  tenantID,
					UserID:    userID,
					Action:    method,
					Resource:  c.Path(),
					Payload:   string(c.Body()), // Rekam semua tombol/ketikan yang dia kirim
					IPAddress: c.IP(),
				}

				// 4. GOROUTINE (Keajaiban Golang)
				// Perintah "go" membuat proses simpan CCTV dikerjakan asisten di latar belakang.
				// Karyawan tidak perlu menunggu proses simpan ini selesai! (Mempertahankan <200ms API)
				go config.DB.Create(&log)
			}
		}

		return err
	}
}
