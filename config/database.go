package config

import (
	"log"
	"os"

	"enterprise-erp/models"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Println("Peringatan: File .env tidak ditemukan")
	//}

	//dsn := fmt.Sprintf(
	//	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_NAME"),
	//	os.Getenv("DB_PORT"),
	//)

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatal("Gagal terhubung ke Database! Error: ", err)
	//}

	//log.Println("Koneksi Database Berhasil!")

	// 1. Sistem akan mencoba membaca URL Database rahasia dari Server Cloud
	dsn := os.Getenv("DATABASE_URL")

	// 2. Jika kosong (berarti sedang dijalankan di komputer laptop Anda), gunakan localhost
	if dsn == "" {
		// PERHATIAN: Pastikan password dan dbname di bawah ini sesuai dengan milik lokal Anda!
		dsn = "host=localhost user=postgres password=admin123 dbname=erp_yd port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database! \n", err)
	}

	log.Println("Koneksi Database Berhasil!")
	DB = db

	// === TAMBAHKAN BARIS INI UNTUK MERESET TOTAL DATABASE ===
	//db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	// =========================================================

	// UPDATE: Menambahkan models.User{} untuk migrasi
	err = db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Item{},
		&models.Account{},
		&models.JournalEntry{},
		&models.JournalLine{},
		&models.Invoice{},           // BARU
		&models.InvoiceLine{},       // BARU
		&models.PurchaseOrder{},     // BARU
		&models.PurchaseOrderLine{}, // BARU
		&models.Payment{},
		&models.Employee{}, // BARU
		&models.Payroll{},  // BARU
		&models.Lead{},     // BARU
		&models.Customer{}, // BARU
		&models.AuditLog{},
		&models.Attendance{},
		&models.Warehouse{},     // BARU
		&models.Inventory{},     // BARU
		&models.StockTransfer{}, // BARU
		&models.Approval{},
		&models.BillOfMaterial{},    // BARU
		&models.BOMComponent{},      // BARU
		&models.ProductionOrder{},   // BARU
		&models.BankStatement{},     // BARU: Modul Bank
		&models.BankStatementLine{}, // BARU: Modul Bank
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database! Error: ", err)
	}

	log.Println("Migrasi tabel Tenant dan User berhasil!")

	db.Exec(`INSERT INTO tenants (id, name, domain) VALUES ('550e8400-e29b-41d4-a716-446655440000', 'PT Enterprise Sejahtera', 'enterprise.com') ON CONFLICT (domain) DO NOTHING`)

	// === TAMBAHKAN PEMBUATAN ADMIN OTOMATIS DI SINI ===
	// 1. Kita buat password "admin123" yang diacak (hash) demi keamanan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)

	// 2. Masukkan akun Super Admin ke database dan kaitkan dengan Tenant di atas
	db.Exec(`
		INSERT INTO users (id, tenant_id, name, email, password, role) 
		VALUES (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440000', 'Super Admin', 'admin@enterprise.com', ?, 'admin') 
		ON CONFLICT (email) DO NOTHING
	`, string(hashedPassword))

	log.Println("Akun Super Admin berhasil disiapkan!")

	DB = db
}
