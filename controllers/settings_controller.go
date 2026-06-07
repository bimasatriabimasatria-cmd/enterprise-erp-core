package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

func GetSettings(c *fiber.Ctx) error {
	var settings models.SystemSettings
	// Ambil data pertama. Jika belum ada, buat default.
	if err := config.DB.FirstOrCreate(&settings, models.SystemSettings{ID: 1}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal akses database"})
	}
	return c.JSON(settings)
}

func UpdateSettings(c *fiber.Ctx) error {
	var input models.SystemSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	var settings models.SystemSettings
	config.DB.FirstOrCreate(&settings, models.SystemSettings{ID: 1})
	config.DB.Model(&settings).Updates(input)

	return c.JSON(fiber.Map{"message": "Identitas berhasil diperbarui di server!"})
}
