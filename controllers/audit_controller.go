package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Lihat Rekaman CCTV (Audit Logs) - Diurutkan dari yang paling baru
func GetAuditLogs(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var logs []models.AuditLog

	// Kita batasi hanya mengambil 100 aktivitas terakhir agar tidak terlalu berat
	if err := config.DB.Where("tenant_id = ?", tenantID).
		Order("created_at desc").
		Limit(100).
		Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil rekaman Audit"})
	}

	return c.JSON(fiber.Map{
		"message": "Menampilkan 100 aktivitas terakhir",
		"data":    logs,
	})
}
