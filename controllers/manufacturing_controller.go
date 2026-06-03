package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// --- 1. BILL OF MATERIALS (RESEP) ---
type BOMComponentInput struct {
	MaterialID string `json:"material_id"`
	Quantity   int    `json:"quantity"`
}

type BOMInput struct {
	Code       string              `json:"code"`
	Name       string              `json:"name"`
	ItemID     string              `json:"item_id"` // Barang jadinya
	Components []BOMComponentInput `json:"components"`
}

func CreateBOM(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var input BOMInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	tx := config.DB.Begin()

	bom := models.BillOfMaterial{
		TenantID: tenantID,
		Code:     input.Code,
		Name:     input.Name,
		ItemID:   input.ItemID,
	}

	if err := tx.Create(&bom).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Kepala Resep"})
	}

	for _, comp := range input.Components {
		bomComp := models.BOMComponent{
			BOMID:      bom.ID,
			MaterialID: comp.MaterialID,
			Quantity:   comp.Quantity,
		}
		if err := tx.Create(&bomComp).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Bahan Baku"})
		}
	}

	tx.Commit()
	return c.Status(201).JSON(fiber.Map{"message": "Resep BOM berhasil dibuat!", "bom_id": bom.ID})
}

// --- 2. PERINTAH PRODUKSI (PRODUCTION ORDER) ---
type ProductionInput struct {
	OrderNumber    string `json:"order_number"`
	BOMID          string `json:"bom_id"`
	WarehouseID    string `json:"warehouse_id"`
	TargetQuantity int    `json:"target_quantity"`
	StartDate      string `json:"start_date"`
}

func CreateProductionOrder(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var input ProductionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, _ := time.Parse("2006-01-02", input.StartDate)

	po := models.ProductionOrder{
		TenantID:       tenantID,
		OrderNumber:    input.OrderNumber,
		BOMID:          input.BOMID,
		WarehouseID:    input.WarehouseID,
		TargetQuantity: input.TargetQuantity,
		Status:         "planned",
		StartDate:      parsedDate,
	}

	// [PERBAIKAN] Kita tangkap pesan penolakan dari database
	if err := config.DB.Create(&po).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan pesanan produksi: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Perintah Produksi diterbitkan", "data": po})
}

// --- 3. SELESAIKAN PRODUKSI (MENGGUNAKAN STOK GLOBAL) ---
func CompleteProduction(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	orderID := c.Params("id")

	tx := config.DB.Begin()

	var order models.ProductionOrder
	if err := tx.Preload("BOM").Preload("BOM.Components").Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Perintah Produksi tidak ditemukan"})
	}

	if order.Status == "completed" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Produksi ini sudah diselesaikan sebelumnya!"})
	}

	// 1. Potong Stok Bahan Baku Global (Tabel Items)
	for _, comp := range order.BOM.Components {
		totalNeeded := comp.Quantity * order.TargetQuantity

		var item models.Item
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, comp.MaterialID).First(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Bahan baku tidak ditemukan di Master Barang!"})
		}

		if item.Stock < totalNeeded {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Stok bahan baku tidak mencukupi untuk produksi!"})
		}

		item.Stock -= totalNeeded
		tx.Save(&item) // Simpan pemotongan stok
	}

	// 2. Tambah Stok Barang Jadi Global (Tabel Items)
	var fgItem models.Item
	if err := tx.Where("tenant_id = ? AND id = ?", tenantID, order.BOM.ItemID).First(&fgItem).Error; err == nil {
		fgItem.Stock += order.TargetQuantity
		tx.Save(&fgItem)
	}

	// 3. Ubah status produksi selesai
	order.Status = "completed"
	tx.Save(&order)

	tx.Commit()

	return c.JSON(fiber.Map{"message": "PRODUKSI SELESAI! Bahan baku dipotong, Barang Jadi telah masuk ke Gudang."})
}

// --- 4. LIHAT DAFTAR RESEP (BOM) ---
func GetBOMs(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var boms []models.BillOfMaterial
	config.DB.Where("tenant_id = ?", tenantID).Find(&boms)
	return c.JSON(fiber.Map{"data": boms})
}

// --- 5. LIHAT DAFTAR PERINTAH PRODUKSI ---
func GetProductionOrders(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var orders []models.ProductionOrder
	config.DB.Where("tenant_id = ?", tenantID).Order("start_date desc").Find(&orders)
	return c.JSON(fiber.Map{"data": orders})
}

// --- 6. HAPUS RESEP (BOM) ---
func DeleteBOM(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	bomID := c.Params("id")

	tx := config.DB.Begin()

	// Hapus komponen anak-anaknya dulu (Relasi)
	if err := tx.Where("bom_id = ?", bomID).Delete(&models.BOMComponent{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus komponen resep"})
	}

	// Hapus resep induknya
	if err := tx.Where("id = ? AND tenant_id = ?", bomID, tenantID).Delete(&models.BillOfMaterial{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus resep BOM"})
	}

	tx.Commit()
	return c.JSON(fiber.Map{"message": "Resep berhasil dihapus!"})
}
