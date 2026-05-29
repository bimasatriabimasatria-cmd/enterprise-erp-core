package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type PaymentInput struct {
	ReceiptNumber string `json:"receipt_number"`
	InvoiceID     string `json:"invoice_id"`
	PaymentDate   string `json:"payment_date"`
}

func ReceivePayment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input PaymentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.PaymentDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah"})
	}

	tx := config.DB.Begin()

	// 1. Cari Faktur yang akan dibayar
	var invoice models.Invoice
	if err := tx.Where("id = ? AND tenant_id = ?", input.InvoiceID, tenantID).First(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Faktur tidak ditemukan"})
	}

	if invoice.Status == "paid" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Faktur ini sudah lunas sebelumnya!"})
	}

	// 2. Buat Catatan Pembayaran
	payment := models.Payment{
		TenantID:      tenantID,
		ReceiptNumber: input.ReceiptNumber,
		InvoiceID:     invoice.ID,
		Amount:        invoice.TotalAmount, // Kita asumsikan pelanggan membayar lunas seketika
		PaymentDate:   parsedDate,
		Method:        "Bank Transfer",
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data pembayaran"})
	}

	// 3. Ubah status Faktur menjadi Lunas
	invoice.Status = "paid"
	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengubah status faktur"})
	}

	// ========================================================
	// [AUTO-JOURNAL] PENERIMAAN UANG (KAS BERTAMBAH, PIUTANG LUNAS)
	// ========================================================

	// Cari Akun Kas Bank (1110) dan Piutang Usaha (1120)
	var bankAccount, arAccount models.Account
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "1110").First(&bankAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Kas Bank (1110) tidak ditemukan."})
	}
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "1120").First(&arAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Piutang Usaha (1120) tidak ditemukan."})
	}

	// Buat Kepala Jurnal
	journal := models.JournalEntry{
		TenantID:    tenantID,
		Reference:   payment.ReceiptNumber,
		Date:        parsedDate,
		Description: "Auto-Jurnal: Pelunasan Faktur " + invoice.InvoiceNumber,
	}
	if err := tx.Create(&journal).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat jurnal penerimaan"})
	}

	// Kas Bank Bertambah (Debit)
	bankLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      bankAccount.ID,
		Debit:          payment.Amount,
		Credit:         0,
	}
	if err := tx.Create(&bankLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Debit Bank"})
	}

	// Piutang Usaha Lunas/Berkurang (Kredit)
	arLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      arAccount.ID,
		Debit:          0,
		Credit:         payment.Amount,
	}
	if err := tx.Create(&arLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Kredit Piutang"})
	}

	// PERMANENKAN TRANSAKSI
	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message": "Pembayaran berhasil diterima. Faktur LUNAS, dan Uang telah masuk ke Jurnal Kas!",
		"payment": payment,
	})
}
