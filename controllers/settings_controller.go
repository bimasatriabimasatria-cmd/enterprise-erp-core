package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2" // Gunakan Fiber!
)

func GetSettings(c *fiber.Ctx) error {
	var settings models.SystemSettings
	config.DB.FirstOrCreate(&settings, models.SystemSettings{ID: 1})
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
