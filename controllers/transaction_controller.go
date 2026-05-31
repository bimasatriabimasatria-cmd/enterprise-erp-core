package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTransaction(c *fiber.Ctx) error {
	// Format data yang akan diterima dari React
	type Request struct {
		SKU  string `json:"sku"`
		Type string `json:"type"` // "IN" atau "OUT"
		Qty  int    `json:"qty"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// 1. Cari barang di gudang berdasarkan SKU
	var item models.Item
	if err := config.DB.Where("sku = ?", req.SKU).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Barang dengan SKU tersebut tidak ditemukan"})
	}

	// 2. Jika barang KELUAR, pastikan stoknya cukup!
	if req.Type == "OUT" && item.Stock < req.Qty {
		return c.Status(400).JSON(fiber.Map{"error": "Stok di gudang tidak mencukupi!"})
	}

	// 3. Update jumlah stok master secara otomatis
	if req.Type == "IN" {
		item.Stock += req.Qty
	} else if req.Type == "OUT" {
		item.Stock -= req.Qty
	}
	config.DB.Save(&item) // Simpan stok terbaru ke database

	// 4. Catat riwayat pergerakan ke Buku Besar
	trx := models.Transaction{
		SKU:  req.SKU,
		Type: req.Type,
		Qty:  req.Qty,
	}
	config.DB.Create(&trx)

	return c.JSON(fiber.Map{
		"message":   "Transaksi berhasil diproses!",
		"new_stock": item.Stock,
	})
}
