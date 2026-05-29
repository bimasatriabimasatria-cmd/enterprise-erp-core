package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type InvoiceLineInput struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type InvoiceInput struct {
	InvoiceNumber string             `json:"invoice_number"`
	CustomerName  string             `json:"customer_name"`
	Date          string             `json:"date"`
	Lines         []InvoiceLineInput `json:"lines"`
}

func CreateInvoice(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input InvoiceInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah"})
	}

	tx := config.DB.Begin()

	invoice := models.Invoice{
		TenantID:      tenantID,
		InvoiceNumber: input.InvoiceNumber,
		CustomerName:  input.CustomerName,
		Date:          parsedDate,
		Status:        "unpaid",
	}

	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Faktur"})
	}

	var grandTotal float64

	// Proses Pemotongan Stok
	for _, lineInput := range input.Lines {
		var item models.Item

		if err := tx.Where("id = ? AND tenant_id = ?", lineInput.ItemID, tenantID).First(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
		}

		if item.Stock < lineInput.Quantity {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Stok tidak cukup untuk: " + item.Name})
		}

		subTotal := item.Price * float64(lineInput.Quantity)
		grandTotal += subTotal

		item.Stock -= lineInput.Quantity
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui stok"})
		}

		invoiceLine := models.InvoiceLine{
			InvoiceID: invoice.ID,
			ItemID:    item.ID,
			Quantity:  lineInput.Quantity,
			UnitPrice: item.Price,
			SubTotal:  subTotal,
		}

		if err := tx.Create(&invoiceLine).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan detail faktur"})
		}
	}

	invoice.TotalAmount = grandTotal
	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan total harga"})
	}

	// ========================================================
	// [BARU] KEAJAIBAN ERP: INTEGRASI AKUNTANSI OTOMATIS
	// ========================================================

	// 1. Cari Akun Piutang Usaha (1120) & Pendapatan (4100)
	var arAccount, revAccount models.Account
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "1120").First(&arAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Piutang Usaha (1120) belum dibuat. Harap buat di menu Chart of Accounts."})
	}
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "4100").First(&revAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Pendapatan Penjualan (4100) belum dibuat. Harap buat di menu Chart of Accounts."})
	}

	// 2. Buat Kepala Jurnal Keuangan
	journal := models.JournalEntry{
		TenantID:    tenantID,
		Reference:   invoice.InvoiceNumber,
		Date:        parsedDate,
		Description: "Auto-Jurnal: Penjualan ke " + invoice.CustomerName,
	}
	if err := tx.Create(&journal).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat jurnal keuangan"})
	}

	// 3. Catat Debit (Piutang bertambah)
	arLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      arAccount.ID,
		Debit:          grandTotal,
		Credit:         0,
	}
	if err := tx.Create(&arLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Debit Piutang"})
	}

	// 4. Catat Kredit (Pendapatan bertambah)
	revLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      revAccount.ID,
		Debit:          0,
		Credit:         grandTotal,
	}
	if err := tx.Create(&revLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Kredit Pendapatan"})
	}

	// ========================================================

	// Permanenkan semua proses (Stok + Faktur + Jurnal Keuangan)
	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message":      "Faktur dibuat, Stok dipotong, dan Jurnal Keuangan otomatis tercatat!",
		"invoice_id":   invoice.ID,
		"total_amount": grandTotal,
	})
}

func GetInvoices(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var invoices []models.Invoice
	if err := config.DB.Where("tenant_id = ?", tenantID).
		Preload("Lines").
		Preload("Lines.Item").
		Order("date desc").
		Find(&invoices).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data faktur"})
	}

	return c.JSON(fiber.Map{"data": invoices})
}
