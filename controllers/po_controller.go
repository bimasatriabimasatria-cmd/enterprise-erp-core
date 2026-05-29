package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type POLineInput struct {
	ItemID    string  `json:"item_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type POInput struct {
	PONumber     string        `json:"po_number"`
	SupplierName string        `json:"supplier_name"`
	Date         string        `json:"date"`
	Lines        []POLineInput `json:"lines"`
}

// 1. Membuat Purchase Order (Pesanan Pembelian)
func CreatePO(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input POInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah"})
	}

	tx := config.DB.Begin()

	po := models.PurchaseOrder{
		TenantID:     tenantID,
		PONumber:     input.PONumber,
		SupplierName: input.SupplierName,
		Date:         parsedDate,
		Status:       "draft", // Masih dipesan, belum sampai
	}

	if err := tx.Create(&po).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat PO"})
	}

	var grandTotal float64

	for _, lineInput := range input.Lines {
		subTotal := lineInput.UnitPrice * float64(lineInput.Quantity)
		grandTotal += subTotal

		poLine := models.PurchaseOrderLine{
			PurchaseOrderID: po.ID,
			ItemID:          lineInput.ItemID,
			Quantity:        lineInput.Quantity,
			UnitPrice:       lineInput.UnitPrice,
			SubTotal:        subTotal,
		}

		if err := tx.Create(&poLine).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan baris PO"})
		}
	}

	po.TotalAmount = grandTotal
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan total PO"})
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message": "Purchase Order berhasil dibuat (Barang belum masuk gudang)",
		"po_id":   po.ID,
	})
}

// 2. Menerima Barang di Gudang (Stok Bertambah)
func ReceivePO(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	poID := c.Params("id") // Mengambil ID PO dari URL

	tx := config.DB.Begin()

	var po models.PurchaseOrder
	// Cari PO dan tarik detail barisnya
	if err := tx.Preload("Lines").Where("id = ? AND tenant_id = ?", poID, tenantID).First(&po).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Purchase Order tidak ditemukan"})
	}

	// Cegah penerimaan ganda
	if po.Status == "received" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "PO ini sudah diterima sebelumnya!"})
	}

	// Tambahkan stok untuk setiap barang di dalam PO
	for _, line := range po.Lines {
		var item models.Item
		if err := tx.First(&item, "id = ?", line.ItemID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Data barang tidak valid"})
		}

		// LOGIKA UTAMA: Tambah Stok
		item.Stock += line.Quantity
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menambah stok gudang"})
		}
	}

	// Ubah status PO menjadi diterima
	po.Status = "received"
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengupdate status PO"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{"message": "Barang telah diterima di gudang, Stok berhasil ditambah!"})
}

// 3. Lihat Semua PO
func GetPOs(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var pos []models.PurchaseOrder
	if err := config.DB.Where("tenant_id = ?", tenantID).Preload("Lines").Preload("Lines.Item").Order("date desc").Find(&pos).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data PO"})
	}

	return c.JSON(fiber.Map{"data": pos})
}
