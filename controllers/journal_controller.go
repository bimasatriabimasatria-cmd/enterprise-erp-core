package controllers

import (
	"fmt"
	"enterprise-erp/config" // Sesuaikan nama module dengan go.mod Anda (enterprise-erp)
	"enterprise-erp/models" // Sesuaikan nama module dengan go.mod Anda
	"github.com/gofiber/fiber/v2"
)

// CreateJournal menangani standar Double-Entry Bookkeeping
func CreateJournal(c *fiber.Ctx) error {
	var input models.JournalEntry

	// 1. Validasi Input JSON dari Frontend/API
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid: " + err.Error(),
		})
	}

	// 2. Ambil TenantID dari Middleware Keamanan (Fiber menggunakan c.Locals)
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses ditolak: Tenant tidak ditemukan di sesi ini",
		})
	}
	input.TenantID = tenantID.(string)

	// 3. Validasi Akuntansi: Total Debit HARUS = Total Kredit
	var totalDebit, totalCredit float64
	for _, line := range input.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	
	// Toleransi kecil untuk floating point
	if totalDebit != totalCredit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Jurnal tidak seimbang (Unbalanced): Total Debit tidak sama dengan Total Kredit",
		})
	}

	// ==========================================
	// 4. ENTERPRISE PATTERN: DATABASE TRANSACTION
	// ==========================================
	tx := config.DB.Begin() // Memulai Transaksi Database

	// [PENTING] Menyuntikkan identitas Tenant ke Postgres untuk melewati RLS
	rlsQuery := fmt.Sprintf(`SET LOCAL "request.jwt.claims" = '{"tenant_id": "%s"}'`, input.TenantID)
	if err := tx.Exec(rlsQuery).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menginisiasi keamanan RLS: " + err.Error(),
		})
	}

	// 5. Simpan Data (GORM otomatis menyimpan Header dan Lines)
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback() // BATALKAN SEMUA jika ada 1 baris gagal!
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan jurnal: " + err.Error(),
		})
	}

	// 6. Jika semua sukses, Sahkan!
	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jurnal Double-Entry berhasil disimpan!",
		"data":    input,
	})
}
