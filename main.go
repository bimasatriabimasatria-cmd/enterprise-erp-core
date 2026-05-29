// @title Enterprise ERP API
// @version 1.0
// @description Sistem ERP Komprehensif (Multi-Tenant, HR, Finance, Manufacturing).
// @termsOfService http://swagger.io/terms/
// @contact.name Tim Developer ERP
// @contact.email support@enterprise-erp.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host enterpriseerpapi-u7w0kgyd.b4a.run
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"time"

	"enterprise-erp/config"
	"enterprise-erp/middlewares"
	"enterprise-erp/routes" // Import jalur API kita

	_ "enterprise-erp/docs" // WAJIB ADA: Mengimpor hasil generate Swagger

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
)

func main() {
	// 1. Inisialisasi Database
	config.ConnectDB()

	// 2. Inisialisasi Framework
	app := fiber.New()

	// === TAMBAHKAN BARIS INI UNTUK MEMBUKA PINTU CORS ===
	app.Use(cors.New())
	// ====================================================

	// ==========================================
	// [BARU] TAMENG API GATEWAY (RATE LIMITING)
	// ==========================================
	app.Use(limiter.New(limiter.Config{
		Max:        100,             // Maksimal 100 klik/permintaan
		Expiration: 1 * time.Minute, // Dalam waktu 1 Menit
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Sistem mendeteksi lalu lintas tidak wajar. IP Anda dikunci sementara untuk mencegah serangan DDoS.",
			})
		},
	}))

	// ==========================================
	// 2. [PERBAIKAN] PASANG CCTV SECARA GLOBAL
	// ==========================================
	// Pasang CCTV ke seluruh aplikasi. (Login & Register aman karena
	// sudah kita kecualikan di dalam file audit_middleware.go)
	app.Use(middlewares.AuditLogger())

	// 3. Route dasar (Health Check)
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Sistem ERP Core Engine Berjalan Normal",
		})
	})

	// Rute Halaman Dokumentasi Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Pendaftaran Semua Rute (Ganti 'app' menjadi 'apiGroup' kecuali Auth karena dia mendaftarkan API sendiri di dalamnya)
	routes.AuthRoutes(app) // Auth tetap pakai app agar CCTV tidak merekam password salah berulang

	// 4. Inisialisasi Auth Routes
	routes.ItemRoutes(app)
	routes.AccountRoutes(app)
	routes.JournalRoutes(app)
	routes.InvoiceRoutes(app)
	routes.PORoutes(app)
	routes.PaymentRoutes(app)
	routes.HRRoutes(app)
	routes.CRMRoutes(app)
	routes.ReportRoutes(app)
	routes.AuditRoutes(app)
	routes.AttendanceRoutes(app)
	routes.WarehouseRoutes(app)
	routes.WorkflowRoutes(app)
	routes.ManufacturingRoutes(app)
	routes.FinanceRoutes(app)
	routes.PortalRoutes(app)

	// 5. Menjalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server ERP berjalan di http://localhost:%s", port)
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
