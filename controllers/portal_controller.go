package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// GetMyInvoices adalah halaman utama Customer Portal
func GetMyInvoices(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string) // ID pengguna yang sedang login (Pelanggan)

	// 1. Cari profil User untuk mendapatkan Email-nya
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	// 2. Validasi Keamanan Tingkat Tinggi:
	// Cek apakah Email pengguna ini benar-benar terdaftar di Master Data Pelanggan Resmi (CRM)
	var customer models.Customer
	if err := config.DB.Where("tenant_id = ? AND email = ?", tenantID, user.Email).First(&customer).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses Ditolak (403)! Email Anda tidak terdaftar sebagai Pelanggan Resmi perusahaan kami.",
		})
	}

	// 3. Tarik HANYA tagihan (Invoice) milik perusahaan pelanggan tersebut
	var invoices []models.Invoice
	if err := config.DB.Preload("Lines").
		Preload("Lines.Item").
		Where("tenant_id = ? AND customer_name = ?", tenantID, customer.Name).
		Find(&invoices).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data tagihan pelanggan"})
	}

	// 4. Kirim data yang sudah difilter
	return c.JSON(fiber.Map{
		"message": "Selamat datang di Customer Portal",
		"customer_info": fiber.Map{
			"company_name": customer.Name,
			"pic_name":     customer.Contact,
			"email":        customer.Email,
		},
		"my_invoices": invoices,
	})
}
