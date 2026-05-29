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

// --- 3. SELESAIKAN PRODUKSI (POTONG BAHAN BAKU, JADIKAN BARANG BARU) ---
func CompleteProduction(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	orderID := c.Params("id")

	tx := config.DB.Begin()

	// 1. Cari Perintah Produksinya
	var order models.ProductionOrder
	if err := tx.Preload("BOM").Preload("BOM.Components").Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Perintah Produksi tidak ditemukan"})
	}

	if order.Status == "completed" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Produksi ini sudah diselesaikan sebelumnya!"})
	}

	// 2. Loop setiap bahan mentah, kalikan dengan target produksi, dan potong stoknya
	for _, comp := range order.BOM.Components {
		totalNeeded := comp.Quantity * order.TargetQuantity

		var inv models.Inventory
		if err := tx.Where("tenant_id = ? AND warehouse_id = ? AND item_id = ?", tenantID, order.WarehouseID, comp.MaterialID).First(&inv).Error; err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Bahan baku tidak ditemukan di gudang pabrik!"})
		}

		if inv.Quantity < totalNeeded {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Stok bahan baku tidak mencukupi untuk produksi!"})
		}

		// Potong Bahan Mentah
		inv.Quantity -= totalNeeded
		tx.Save(&inv)
	}

	// 3. Tambahkan Barang Jadi (Finished Good) ke dalam Gudang
	var fgInv models.Inventory
	err := tx.Where("tenant_id = ? AND warehouse_id = ? AND item_id = ?", tenantID, order.WarehouseID, order.BOM.ItemID).First(&fgInv).Error
	if err != nil {
		// Jika belum pernah ada barang jadi ini di gudang, buat data baru
		fgInv = models.Inventory{
			TenantID:    tenantID,
			WarehouseID: order.WarehouseID,
			ItemID:      order.BOM.ItemID,
			Quantity:    order.TargetQuantity,
		}
		tx.Create(&fgInv)
	} else {
		// Jika sudah ada, tambahkan stoknya
		fgInv.Quantity += order.TargetQuantity
		tx.Save(&fgInv)
	}

	// 4. Ubah status produksi selesai
	order.Status = "completed"
	tx.Save(&order)

	tx.Commit()

	return c.JSON(fiber.Map{"message": "PRODUKSI SELESAI! Bahan baku dipotong, Barang Jadi telah masuk ke Gudang."})
}
