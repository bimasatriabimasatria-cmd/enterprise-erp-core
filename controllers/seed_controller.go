package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(c *fiber.Ctx) error {
	// 1. Ambil Perusahaan Utama (Tenant)
	var tenant models.Tenant
	if err := config.DB.First(&tenant).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Tenant belum ada, jalankan ulang server untuk trigger auto-migrate."})
	}

	// 2. Buat Password Universal: "password123"
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	passStr := string(password)

	// 3. Pasukan Karyawan Baru
	users := []models.User{
		{TenantID: tenant.ID, Name: "Bapak Direktur", Email: "admin@erp.com", Password: passStr, Role: "super_admin"},
		{TenantID: tenant.ID, Name: "Pak Gudang", Email: "gudang@erp.com", Password: passStr, Role: "warehouse_staff"},
		{TenantID: tenant.ID, Name: "Ibu HRD", Email: "hrd@erp.com", Password: passStr, Role: "hr_staff"},
		{TenantID: tenant.ID, Name: "Mbak Kasir", Email: "finance@erp.com", Password: passStr, Role: "finance_staff"},
	}

	// 4. Suntikkan ke Database (Abaikan jika sudah ada)
	for _, u := range users {
		config.DB.Where("email = ?", u.Email).FirstOrCreate(&u)
	}

	return c.JSON(fiber.Map{
		"message": "SELDING RBAC BERHASIL! 4 Akun telah diciptakan. Gunakan password: password123",
		"accounts": []string{
			"admin@erp.com (Super Admin)",
			"gudang@erp.com (Gudang & Pabrik)",
			"hrd@erp.com (HR & Payroll)",
			"finance@erp.com (Finance & CRM)",
		},
	})
}
