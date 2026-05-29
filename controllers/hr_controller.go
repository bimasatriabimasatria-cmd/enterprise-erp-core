package controllers

import (
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// --- BAGIAN KARYAWAN ---
type EmployeeInput struct {
	NIK         string  `json:"nik"`
	Name        string  `json:"name"`
	Position    string  `json:"position"`
	BasicSalary float64 `json:"basic_salary"`
	HireDate    string  `json:"hire_date"`
}

func CreateEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input EmployeeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.HireDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah"})
	}

	employee := models.Employee{
		TenantID:    tenantID,
		NIK:         input.NIK,
		Name:        input.Name,
		Position:    input.Position,
		BasicSalary: input.BasicSalary,
		HireDate:    parsedDate,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data Karyawan (NIK mungkin sudah dipakai)"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Karyawan berhasil didaftarkan", "data": employee})
}

func GetEmployees(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var employees []models.Employee
	config.DB.Where("tenant_id = ?", tenantID).Find(&employees)
	return c.JSON(fiber.Map{"data": employees})
}

// --- BAGIAN PAYROLL (PENGGAJIAN) ---
type PayrollInput struct {
	EmployeeID  string  `json:"employee_id"`
	Period      string  `json:"period"` // cth: "2026-05"
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
	PaymentDate string  `json:"payment_date"`
}

func ProcessPayroll(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input PayrollInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.PaymentDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal salah"})
	}

	tx := config.DB.Begin()

	// 1. Cari Karyawan untuk mengambil Gaji Pokoknya
	var employee models.Employee
	if err := tx.Where("id = ? AND tenant_id = ?", input.EmployeeID, tenantID).First(&employee).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Karyawan tidak ditemukan"})
	}

	// 2. Hitung Gaji Bersih (Net Pay) = Pokok + Tunjangan - Potongan
	netPay := employee.BasicSalary + input.Allowances - input.Deductions

	// 3. Simpan Catatan Slip Gaji
	payroll := models.Payroll{
		TenantID:    tenantID,
		EmployeeID:  employee.ID,
		Period:      input.Period,
		BasicSalary: employee.BasicSalary,
		Allowances:  input.Allowances,
		Deductions:  input.Deductions,
		NetPay:      netPay,
		PaymentDate: parsedDate,
		Status:      "paid",
	}

	if err := tx.Create(&payroll).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses slip gaji"})
	}

	// ========================================================
	// [AUTO-JOURNAL] BIAYA GAJI (BEBAN BERTAMBAH, KAS BERKURANG)
	// ========================================================

	// Cari Akun Beban Gaji (5100) dan Kas Bank (1110)
	var expenseAccount, bankAccount models.Account
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "5100").First(&expenseAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Beban Gaji (5100) belum dibuat di Chart of Accounts."})
	}
	if err := tx.Where("tenant_id = ? AND account_code = ?", tenantID, "1110").First(&bankAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Akun Kas Bank (1110) belum dibuat."})
	}

	// Buat Kepala Jurnal Keuangan
	journal := models.JournalEntry{
		TenantID:    tenantID,
		Reference:   "PR-" + input.Period + "-" + employee.NIK,
		Date:        parsedDate,
		Description: "Auto-Jurnal: Pembayaran Gaji " + employee.Name + " Periode " + input.Period,
	}
	if err := tx.Create(&journal).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat jurnal penggajian"})
	}

	// Beban Gaji Bertambah (Debit)
	expenseLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      expenseAccount.ID,
		Debit:          netPay, // Perusahaan kehilangan uang sebesar netPay
		Credit:         0,
	}
	if err := tx.Create(&expenseLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Debit Beban Gaji"})
	}

	// Uang di Bank Berkurang (Kredit)
	bankLine := models.JournalLine{
		JournalEntryID: journal.ID,
		AccountID:      bankAccount.ID,
		Debit:          0,
		Credit:         netPay,
	}
	if err := tx.Create(&bankLine).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mencatat Kredit Kas Bank"})
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"message": "Gaji berhasil dicairkan dan Jurnal Keuangan otomatis terpotong!",
		"net_pay": netPay,
	})
}
