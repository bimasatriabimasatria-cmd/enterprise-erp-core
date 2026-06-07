package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetSettings(c *fiber.Ctx) error {
	// 📢 TAMBAHKAN LOG INI
	log.Println("PINTU API /api/settings DITOK SEORANG PENGUNJUNG!")

	var settings models.SystemSettings
	if err := config.DB.First(&settings, 1).Error; err != nil {
		return c.JSON(fiber.Map{"company_name": "ENTERPRISE", "logo": ""})
	}
	return c.JSON(settings)
}

func GetSettings(c *fiber.Ctx) error {
	var settings models.SystemSettings
	// Ambil data pertama. Jika tabel kosong, kembalikan default.
	if err := config.DB.First(&settings, 1).Error; err != nil {
		// Jika belum ada data, beri response default tanpa error 404
		return c.JSON(fiber.Map{"company_name": "ENTERPRISE", "logo": ""})
	}
	return c.JSON(settings)
}

func UpdateSettings(c *fiber.Ctx) error {
	var input models.SystemSettings
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	var settings models.SystemSettings
	config.DB.FirstOrCreate(&settings, models.SystemSettings{ID: 1})
	config.DB.Model(&settings).Updates(input)

	return c.JSON(fiber.Map{"message": "Berhasil!"})
}
