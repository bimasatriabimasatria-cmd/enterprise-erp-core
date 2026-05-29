package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Struktur input dari klien
type AccountInput struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
}

// Tambah Akun Baru
func CreateAccount(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input AccountInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Validasi Tipe Akun Standar Akuntansi
	validTypes := map[string]bool{"asset": true, "liability": true, "equity": true, "revenue": true, "expense": true}
	if !validTypes[input.AccountType] {
		return c.Status(400).JSON(fiber.Map{"error": "Tipe akun tidak valid. Gunakan: asset, liability, equity, revenue, atau expense"})
	}

	account := models.Account{
		TenantID:    tenantID,
		AccountCode: input.AccountCode,
		AccountName: input.AccountName,
		AccountType: input.AccountType,
	}

	if err := config.DB.Create(&account).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan akun. Pastikan Kode Akun belum digunakan."})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Akun berhasil ditambahkan ke Buku Besar",
		"data":    account,
	})
}

// Lihat Daftar Akun (Hanya milik perusahaan sendiri, diurutkan berdasarkan kode)
func GetAccounts(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var accounts []models.Account
	// Perhatikan .Order("account_code asc") agar laporan keuangan urut dari Harta sampai Beban
	if err := config.DB.Where("tenant_id = ?", tenantID).Order("account_code asc").Find(&accounts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data akun"})
	}

	return c.JSON(fiber.Map{"data": accounts})
}
