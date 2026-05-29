package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Struktur data dari input form
type ItemInput struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// CreateItem membuat barang baru
// @Summary Tambah barang ke gudang
// @Description Memasukkan barang baru ke dalam sistem, otomatis terikat dengan Tenant ID user
// @Tags Inventory & Warehouse
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer TOKEN_JWT_KAMU"
// @Param request body map[string]interface{} true "Format JSON Barang"
// @Success 201 {object} map[string]interface{}
// @Router /api/items [post]
func CreateItem(c *fiber.Ctx) error {
	// Ambil identitas perusahaan dari Satpam (Middleware)
	tenantID := c.Locals("tenant_id").(string)

	var input ItemInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	item := models.Item{
		TenantID:    tenantID,
		SKU:         input.SKU,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan barang. Pastikan SKU unik dalam perusahaan Anda."})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Barang berhasil ditambahkan",
		"data":    item,
	})
}

// Lihat Semua Barang (Hanya milik perusahaan sendiri)
func GetItems(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var items []models.Item
	// FILTERING OTOMATIS: Kunci utama Multi-Tenant
	if err := config.DB.Where("tenant_id = ?", tenantID).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data barang"})
	}

	return c.JSON(fiber.Map{"data": items})
}
