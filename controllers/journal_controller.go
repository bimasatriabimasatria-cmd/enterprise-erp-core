package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Struktur input baris dari pengguna
type JournalLineInput struct {
	AccountID string  `json:"account_id"`
	Debit     float64 `json:"debit"`
	Credit    float64 `json:"credit"`
}

// Struktur input header dari pengguna
type JournalInput struct {
	Reference   string             `json:"reference"`
	Date        string             `json:"date"` // Format: YYYY-MM-DD
	Description string             `json:"description"`
	Lines       []JournalLineInput `json:"lines"`
}

// Buat Jurnal Baru
func CreateJournal(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input JournalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// 1. Validasi Double-Entry (Debit = Kredit)
	var totalDebit, totalCredit float64
	for _, line := range input.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if totalDebit != totalCredit {
		return c.Status(400).JSON(fiber.Map{
			"error":        "Jurnal Tidak Seimbang (Unbalanced)!",
			"total_debit":  totalDebit,
			"total_credit": totalCredit,
		})
	}

	// Parsing Tanggal
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah. Gunakan YYYY-MM-DD"})
	}

	// 2. Mulai Database Transaction (ACID)
	tx := config.DB.Begin()

	// Simpan Header Jurnal
	journal := models.JournalEntry{
		TenantID:    tenantID,
		Reference:   input.Reference,
		Date:        parsedDate,
		Description: input.Description,
	}

	if err := tx.Create(&journal).Error; err != nil {
		tx.Rollback() // Batalkan semua jika gagal
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Header Jurnal"})
	}

	// Simpan Baris Jurnal (Lines)
	for _, lineInput := range input.Lines {
		line := models.JournalLine{
			JournalEntryID: journal.ID,
			AccountID:      lineInput.AccountID,
			Debit:          lineInput.Debit,
			Credit:         lineInput.Credit,
		}

		if err := tx.Create(&line).Error; err != nil {
			tx.Rollback() // Batalkan semua termasuk Header jika ada 1 baris yang gagal
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Baris Jurnal. Pastikan Account ID valid."})
		}
	}

	// 3. Jika semua aman, Permanenkan data (Commit)
	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message": "Jurnal berhasil dicatat",
		"id":      journal.ID,
	})
}

// Ambil Daftar Jurnal beserta rincian Debit/Kreditnya
func GetJournals(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var journals []models.JournalEntry
	// Preload("Lines") bertugas mengambil baris detail jurnal secara otomatis
	// Preload("Lines.Account") akan menarik nama akun sekalian (contoh: Kas, Beban, dll)
	if err := config.DB.Where("tenant_id = ?", tenantID).
		Preload("Lines").
		Preload("Lines.Account").
		Order("date desc").
		Find(&journals).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data jurnal"})
	}

	return c.JSON(fiber.Map{"data": journals})
}
