package models

import (
	"time"

	"gorm.io/gorm"
)

// Employee adalah Master Data Karyawan
type Employee struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	NIK         string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_nik"` // Nomor Induk Karyawan
	Name        string    `gorm:"type:varchar(255);not null"`
	Position    string    `gorm:"type:varchar(100);not null"`
	BasicSalary float64   `gorm:"type:decimal(15,2);not null"` // Gaji Pokok Bulanan
	HireDate    time.Time `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Payroll adalah Catatan Slip Gaji (Penggajian)
type Payroll struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	EmployeeID  string    `gorm:"type:uuid;not null;index"`
	Period      string    `gorm:"type:varchar(20);not null"` // Contoh: "2026-05" (Bulan Mei 2026)
	BasicSalary float64   `gorm:"type:decimal(15,2);not null"`
	Allowances  float64   `gorm:"type:decimal(15,2);default:0"` // Tunjangan (Bonus, Transport)
	Deductions  float64   `gorm:"type:decimal(15,2);default:0"` // Potongan (Absen, Pajak, BPJS)
	NetPay      float64   `gorm:"type:decimal(15,2);not null"`  // Gaji Bersih yang ditransfer
	PaymentDate time.Time `gorm:"not null"`
	Status      string    `gorm:"type:varchar(50);default:'paid'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relasi ke tabel Karyawan
	Employee Employee `gorm:"foreignKey:EmployeeID"`
}
