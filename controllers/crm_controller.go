package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type LeadInput struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

// 1. Tambah Calon Pelanggan Baru (Prospek)
func CreateLead(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input LeadInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	lead := models.Lead{
		TenantID: tenantID,
		Name:     input.Name,
		Company:  input.Company,
		Email:    input.Email,
		Phone:    input.Phone,
		Status:   "new",
	}

	if err := config.DB.Create(&lead).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Lead"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Calon pelanggan berhasil ditambahkan", "data": lead})
}

// 2. Lihat Daftar Calon Pelanggan
func GetLeads(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var leads []models.Lead
	config.DB.Where("tenant_id = ?", tenantID).Find(&leads)
	return c.JSON(fiber.Map{"data": leads})
}

// 3. Konversi Lead menjadi Customer Resmi (Tanda Tangan Kontrak)
func ConvertLead(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	leadID := c.Params("id") // Ambil ID dari URL

	tx := config.DB.Begin()

	// Cari Lead
	var lead models.Lead
	if err := tx.Where("id = ? AND tenant_id = ?", leadID, tenantID).First(&lead).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Lead tidak ditemukan"})
	}

	// Cek apakah sudah pernah di-convert
	if lead.Status == "converted" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Lead ini sudah menjadi Customer Resmi sebelumnya!"})
	}

	// Pindahkan data ke tabel Customer
	customer := models.Customer{
		TenantID: tenantID,
		Name:     lead.Company, // Nama perusahaan jadi nama pelanggan utama
		Contact:  lead.Name,    // Nama orangnya jadi PIC
		Email:    lead.Email,
		Phone:    lead.Phone,
	}

	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat Customer baru"})
	}

	// Ubah status Lead menjadi "converted" agar tidak dobel
	lead.Status = "converted"
	if err := tx.Save(&lead).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengubah status Lead"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message":  "Prospek berhasil di-Goal-kan! Otomatis terdaftar sebagai Pelanggan Resmi.",
		"customer": customer,
	})
}

// 4. Lihat Daftar Pelanggan Resmi
func GetCustomers(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var customers []models.Customer
	config.DB.Where("tenant_id = ?", tenantID).Find(&customers)
	return c.JSON(fiber.Map{"data": customers})
}
