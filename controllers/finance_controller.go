package controllers

import (
	"encoding/csv"
	"strconv"
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// 1. Mengunggah dan Membaca File CSV dari Bank
func UploadBankStatement(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Tangkap file bernama "file" dari request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File tidak ditemukan atau format salah"})
	}

	// Buka isi filenya
	f, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file"})
	}
	defer f.Close()

	// Baca sebagai CSV
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca isi file CSV"})
	}

	tx := config.DB.Begin()

	// Buat Kepala Rekaman Mutasi
	statement := models.BankStatement{
		TenantID: tenantID,
		BankName: "Bank BCA",
		Period:   "Juni 2026",
	}
	tx.Create(&statement)

	// Masukkan setiap baris di Excel/CSV ke database kita
	for _, row := range records {
		// Asumsi format CSV: Kolom 1 = Tanggal, Kolom 2 = Deskripsi, Kolom 3 = Nominal
		if len(row) < 3 {
			continue
		}

		parsedDate, _ := time.Parse("2006-01-02", row[0])
		amount, _ := strconv.ParseFloat(row[2], 64)

		line := models.BankStatementLine{
			StatementID: statement.ID,
			Date:        parsedDate,
			Description: row[1],
			Amount:      amount,
		}
		tx.Create(&line)
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message":      "File Mutasi Bank berhasil diimpor!",
		"statement_id": statement.ID,
	})
}

// 2. Mesin Pencocok Otomatis (Auto-Reconcile)
func AutoReconcile(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	statementID := c.Params("id")

	tx := config.DB.Begin()

	// Ambil semua baris mutasi dari Bank yang belum dicocokkan
	var unMatchedLines []models.BankStatementLine
	tx.Where("statement_id = ? AND is_reconciled = ?", statementID, false).Find(&unMatchedLines)

	matchCount := 0

	// Loop setiap transaksi Bank
	for _, line := range unMatchedLines {
		var payment models.Payment

		// KEAJAIBAN: Cari data Pembayaran di ERP yang:
		// 1. Nominalnya sama persis dengan uang yang masuk di Bank
		// 2. Belum pernah dicocokkan (ID-nya belum ada di tabel mutasi bank)
		err := tx.Raw(`
			SELECT * FROM payments 
			WHERE tenant_id = ? 
			AND amount = ? 
			AND id NOT IN (SELECT payment_id FROM bank_statement_lines WHERE payment_id IS NOT NULL) 
			LIMIT 1
		`, tenantID, line.Amount).Scan(&payment).Error

		if err == nil && payment.ID != "" {
			// JIKA KETEMU! Kita kawinkan datanya
			line.IsReconciled = true
			line.PaymentID = &payment.ID
			tx.Save(&line)
			matchCount++
		}
	}

	// Update status kepala mutasi jika semua sudah cocok
	var countRemaining int64
	tx.Model(&models.BankStatementLine{}).Where("statement_id = ? AND is_reconciled = ?", statementID, false).Count(&countRemaining)

	if countRemaining == 0 {
		tx.Model(&models.BankStatement{}).Where("id = ?", statementID).Update("status", "reconciled")
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message":              "Proses Rekonsiliasi Selesai!",
		"matched_transactions": matchCount,
		"unmatched_remaining":  countRemaining,
	})
}

// 3. Lihat Hasil Mutasi Bank
func GetBankStatements(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var statements []models.BankStatement
	config.DB.Where("tenant_id = ?", tenantID).Preload("Lines").Find(&statements)
	return c.JSON(fiber.Map{"data": statements})
}
