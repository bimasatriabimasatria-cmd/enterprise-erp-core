package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type ApprovalRequestInput struct {
	DocumentType string `json:"document_type"`
	DocumentID   string `json:"document_id"`
	Notes        string `json:"notes"`
}

// 1. Staf Mengajukan Persetujuan
func RequestApproval(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string) // Ambil ID staf yang sedang login

	var input ApprovalRequestInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	approval := models.Approval{
		TenantID:     tenantID,
		DocumentType: input.DocumentType,
		DocumentID:   input.DocumentID,
		RequestedBy:  userID,
		Notes:        input.Notes,
		Status:       "pending",
	}

	if err := config.DB.Create(&approval).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat pengajuan persetujuan"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Dokumen berhasil diajukan dan sedang menunggu persetujuan Manajer",
		"data":    approval,
	})
}

// 2. Melihat Daftar Antrean Persetujuan
func GetPendingApprovals(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var approvals []models.Approval
	// Hanya ambil yang statusnya masih pending
	if err := config.DB.Where("tenant_id = ? AND status = ?", tenantID, "pending").Find(&approvals).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data antrean persetujuan"})
	}

	return c.JSON(fiber.Map{"data": approvals})
}

type ProcessApprovalInput struct {
	Status string `json:"status"` // 'approved' atau 'rejected'
	Notes  string `json:"notes"`  // Alasan dari manajer
}

// 3. Manajer Memproses (Setuju / Tolak)
func ProcessApproval(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	managerID := c.Locals("user_id").(string) // ID Manajer/Admin yang memproses
	approvalID := c.Params("id")              // ID tiket persetujuan dari URL

	var input ProcessApprovalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	if input.Status != "approved" && input.Status != "rejected" {
		return c.Status(400).JSON(fiber.Map{"error": "Status harus 'approved' atau 'rejected'"})
	}

	tx := config.DB.Begin()

	var approval models.Approval
	if err := tx.Where("id = ? AND tenant_id = ?", approvalID, tenantID).First(&approval).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Dokumen persetujuan tidak ditemukan"})
	}

	if approval.Status != "pending" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Dokumen ini sudah diproses sebelumnya!"})
	}

	// Update status dan siapa manajer yang menyetujuinya
	approval.Status = input.Status
	approval.ApproverID = &managerID
	approval.Notes = input.Notes

	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses persetujuan"})
	}

	// CATATAN: Di sistem ERP sungguhan, di titik ini kita juga akan mengubah
	// status dokumen aslinya (misal: PO.Status = 'approved').
	// Tapi untuk fondasi, kita mencatat riwayat di tabel Workflow terlebih dahulu.

	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Dokumen berhasil di-" + input.Status,
		"data":    approval,
	})
}
