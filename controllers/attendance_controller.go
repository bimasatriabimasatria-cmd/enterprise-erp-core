package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

type ClockInInput struct {
	EmployeeID string `json:"employee_id"`
	Time       string `json:"time"` // Format: "2026-05-29T08:30:00Z" (Format ISO 8601)
}

// 1. Absen Masuk (Clock In)
func ClockIn(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input ClockInInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Parsing Waktu Absen
	clockInTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format waktu salah. Gunakan standar RFC3339 (Contoh: 2026-05-29T08:30:00Z)"})
	}

	// Ambil tanggalnya saja (tanpa jam) untuk direkam ke kolom Date
	dateOnly := time.Date(clockInTime.Year(), clockInTime.Month(), clockInTime.Day(), 0, 0, 0, 0, clockInTime.Location())

	tx := config.DB.Begin()

	// Pastikan karyawan valid
	var employee models.Employee
	if err := tx.Where("id = ? AND tenant_id = ?", input.EmployeeID, tenantID).First(&employee).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Karyawan tidak ditemukan"})
	}

	// Cek apakah hari ini sudah absen masuk?
	var existing models.Attendance
	if err := tx.Where("employee_id = ? AND date = ?", employee.ID, dateOnly).First(&existing).Error; err == nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Karyawan ini sudah melakukan Clock In hari ini!"})
	}

	// LOGIKA ERP: Jam Masuk Standar adalah 09:00:00
	// Jika lebih dari jam 9, maka statusnya "late" (Terlambat)
	status := "on-time"
	if clockInTime.Hour() >= 9 && clockInTime.Minute() > 0 {
		status = "late"
	}

	attendance := models.Attendance{
		TenantID:   tenantID,
		EmployeeID: employee.ID,
		Date:       dateOnly,
		CheckIn:    &clockInTime,
		Status:     status,
	}

	if err := tx.Create(&attendance).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data absensi"})
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message": "Absen Masuk Berhasil",
		"status":  status,
		"time":    clockInTime,
	})
}

type ClockOutInput struct {
	EmployeeID string `json:"employee_id"`
	Time       string `json:"time"`
}

// 2. Absen Pulang (Clock Out)
func ClockOut(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input ClockOutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	clockOutTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format waktu salah"})
	}

	dateOnly := time.Date(clockOutTime.Year(), clockOutTime.Month(), clockOutTime.Day(), 0, 0, 0, 0, clockOutTime.Location())

	// Cari data absensi masuk karyawan di hari tersebut
	var attendance models.Attendance
	if err := config.DB.Where("tenant_id = ? AND employee_id = ? AND date = ?", tenantID, input.EmployeeID, dateOnly).First(&attendance).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Anda belum melakukan Clock In hari ini!"})
	}

	// Update data dengan jam pulang
	attendance.CheckOut = &clockOutTime
	if err := config.DB.Save(&attendance).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data absensi pulang"})
	}

	return c.JSON(fiber.Map{
		"message": "Absen Pulang Berhasil",
		"time":    clockOutTime,
	})
}

// 3. Lihat Data Absensi (Hanya bisa diakses admin/HRD)
func GetAttendances(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var attendances []models.Attendance
	if err := config.DB.Where("tenant_id = ?", tenantID).Preload("Employee").Order("date desc").Find(&attendances).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data absensi"})
	}

	return c.JSON(fiber.Map{"data": attendances})
}
