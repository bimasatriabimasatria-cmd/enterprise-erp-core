package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Mengambil semua data karyawan
func GetEmployees(c *fiber.Ctx) error {
	var employees []models.Employee
	config.DB.Find(&employees)
	return c.JSON(fiber.Map{"data": employees})
}

// Menambah karyawan baru
func CreateEmployee(c *fiber.Ctx) error {
	var employee models.Employee
	if err := c.BodyParser(&employee); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Simpan ke Supabase
	if err := config.DB.Create(&employee).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data karyawan"})
	}

	return c.JSON(fiber.Map{
		"message": "Karyawan berhasil ditambahkan!",
		"data":    employee,
	})
}
