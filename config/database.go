package config

import (
	"log"
	"os"

	"enterprise-erp-core/models" // Pastikan import path ini sama di semua file!

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=admin123 dbname=erp_yd port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	// PreferSimpleProtocol: true SANGAT BAGUS untuk Supabase (PgBouncer connection pooling)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, 
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal terhubung ke database! \n", err)
	}
	log.Println("Koneksi Database Berhasil!")

	// ==========================================
	// CTO ADVICE: PISAHKAN MIGRASI TABEL INTI
	// Jangan masukkan tabel FICO (Tenant, Account, Journal) ke sini 
	// karena kita sudah membuatnya dengan SQL khusus (RLS & Constraints).
	// ==========================================
	err = db.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Invoice{},           
		&models.InvoiceLine{},       
		&models.PurchaseOrder{},     
		&models.PurchaseOrderLine{}, 
		&models.Payment{},
		&models.Employee{}, 
		&models.Payroll{},  
		&models.Lead{},     
		&models.Customer{}, 
		&models.AuditLog{},
		&models.Attendance{},
		&models.Warehouse{},     
		&models.Inventory{},     
		&models.StockTransfer{}, 
		&models.Approval{},
		&models.BillOfMaterial{},    
		&models.BOMComponent{},      
		&models.ProductionOrder{},   
		&models.BankStatement{},     
		&models.BankStatementLine{}, 
		&models.Transaction{},
		&models.SystemSettings{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database! Error: ", err)
	}
	log.Println("AutoMigrate (Non-Core Tables) berhasil!")

	// ==========================================
	// SEEDER: MENGGUNAKAN UUID DARI STEP SEBELUMNYA
	// UUID: 620c6b50-8e5d-4da8-8c5c-60c7c4052079 (PT ERP Maju Bersama)
	// ==========================================
	
	// Pastikan Tenant PT ERP Maju Bersama benar-benar ada
	db.Exec(`INSERT INTO tenants (id, company_name, currency) VALUES ('620c6b50-8e5d-4da8-8c5c-60c7c4052079', 'PT ERP Maju Bersama', 'IDR') ON CONFLICT (id) DO NOTHING`)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	db.Exec(`
		INSERT INTO users (id, tenant_id, name, email, password, role) 
		VALUES (gen_random_uuid(), '620c6b50-8e5d-4da8-8c5c-60c7c4052079', 'Super Admin', 'admin@enterprise.com', ?, 'admin') 
		ON CONFLICT (email) DO NOTHING
	`, string(hashedPassword))

	log.Println("Akun Super Admin berhasil disiapkan!")

	DB = db
}
