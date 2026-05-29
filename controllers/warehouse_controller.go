package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// 1. Buat Gudang Baru
type WarehouseInput struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func CreateWarehouse(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var input WarehouseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	wh := models.Warehouse{
		TenantID: tenantID,
		Code:     input.Code,
		Name:     input.Name,
		Location: input.Location,
	}

	if err := config.DB.Create(&wh).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan Gudang"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Gudang berhasil ditambahkan", "data": wh})
}

// 2. Set Stok Awal di Gudang Tertentu (Untuk Inisialisasi)
type InventoryInput struct {
	WarehouseID string `json:"warehouse_id"`
	ItemID      string `json:"item_id"`
	Quantity    int    `json:"quantity"`
}

func SetInventory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var input InventoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Cari apakah barang ini sudah ada di gudang tersebut
	var inv models.Inventory
	err := config.DB.Where("tenant_id = ? AND warehouse_id = ? AND item_id = ?", tenantID, input.WarehouseID, input.ItemID).First(&inv).Error

	if err != nil {
		// Jika belum ada, buat baru
		inv = models.Inventory{
			TenantID:    tenantID,
			WarehouseID: input.WarehouseID,
			ItemID:      input.ItemID,
			Quantity:    input.Quantity,
		}
		config.DB.Create(&inv)
	} else {
		// Jika sudah ada, update jumlahnya
		inv.Quantity = input.Quantity
		config.DB.Save(&inv)
	}

	return c.JSON(fiber.Map{"message": "Stok gudang berhasil diatur", "data": inv})
}

// 3. Mutasi Barang (Transfer Stok Antar Gudang)
type TransferInput struct {
	Reference     string `json:"reference"`
	SourceID      string `json:"source_id"`
	DestinationID string `json:"destination_id"`
	ItemID        string `json:"item_id"`
	Quantity      int    `json:"quantity"`
}

func TransferStock(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var input TransferInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	if input.SourceID == input.DestinationID {
		return c.Status(400).JSON(fiber.Map{"error": "Gudang asal dan tujuan tidak boleh sama!"})
	}

	tx := config.DB.Begin()

	// 1. Cek stok di Gudang Asal
	var sourceInv models.Inventory
	if err := tx.Where("tenant_id = ? AND warehouse_id = ? AND item_id = ?", tenantID, input.SourceID, input.ItemID).First(&sourceInv).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan di Gudang Asal"})
	}

	if sourceInv.Quantity < input.Quantity {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Stok di Gudang Asal tidak mencukupi untuk mutasi!"})
	}

	// 2. Kurangi stok di Gudang Asal
	sourceInv.Quantity -= input.Quantity
	if err := tx.Save(&sourceInv).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengurangi stok gudang asal"})
	}

	// 3. Tambahkan stok di Gudang Tujuan
	var destInv models.Inventory
	if err := tx.Where("tenant_id = ? AND warehouse_id = ? AND item_id = ?", tenantID, input.DestinationID, input.ItemID).First(&destInv).Error; err != nil {
		// Jika belum pernah ada barang ini di gudang tujuan, buat catatannya
		destInv = models.Inventory{
			TenantID:    tenantID,
			WarehouseID: input.DestinationID,
			ItemID:      input.ItemID,
			Quantity:    input.Quantity,
		}
		if err := tx.Create(&destInv).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat stok di gudang tujuan"})
		}
	} else {
		// Jika sudah ada, tinggal tambahkan
		destInv.Quantity += input.Quantity
		if err := tx.Save(&destInv).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menambah stok gudang tujuan"})
		}
	}

	// 4. Catat Riwayat Mutasi
	transfer := models.StockTransfer{
		TenantID:      tenantID,
		Reference:     input.Reference,
		SourceID:      input.SourceID,
		DestinationID: input.DestinationID,
		ItemID:        input.ItemID,
		Quantity:      input.Quantity,
		TransferDate:  time.Now(),
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat riwayat mutasi"})
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message":          "Mutasi Barang Berhasil! Stok telah berpindah lokasi.",
		"source_remaining": sourceInv.Quantity,
		"destination_new":  destInv.Quantity,
	})
}

// 4. Lihat Daftar Gudang
func GetWarehouses(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var warehouses []models.Warehouse
	config.DB.Where("tenant_id = ?", tenantID).Find(&warehouses)
	return c.JSON(fiber.Map{"data": warehouses})
}

// 5. Lihat Inventaris per Gudang
func GetInventory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var inventories []models.Inventory
	config.DB.Where("tenant_id = ?", tenantID).Preload("Warehouse").Preload("Item").Find(&inventories)
	return c.JSON(fiber.Map{"data": inventories})
}
