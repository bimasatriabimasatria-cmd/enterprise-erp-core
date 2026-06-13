package controllers

import (
	"fmt"
	"enterprise-erp/config" // Pastikan import path ini sesuai dengan modul Anda
	"enterprise-erp/models" // Pastikan import path ini sesuai dengan modul Anda

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

	// 2. Ambil TenantID dari Middleware Keamanan
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses ditolak: Tenant tidak ditemukan",
		})
	}
	input.TenantID = tenantID.(string)

	// 3. Validasi Akuntansi: Total Debit HARUS = Total Kredit
	var totalDebit, totalCredit float64
	for _, line := range input.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	
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
			"error": "Gagal menginisiasi keamanan RLS",
		})
	}

	// 5. Simpan Data (GORM akan otomatis menyimpan Header dan Lines)
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback() // BATALKAN SEMUA jika ada 1 baris yang gagal
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan jurnal: " + err.Error(),
		})
	}

	// 6. Sahkan data ke database permanen
	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jurnal Double-Entry berhasil disimpan!",
		"data":    input,
	})
}

// GetJournals (Dummy response agar routes Anda tidak error)
func GetJournals(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "List jurnal akan diimplementasikan nanti"})
}
