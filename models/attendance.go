package models

import (
	"time"

	"gorm.io/gorm"
)

// Attendance mencatat jam masuk dan keluar karyawan setiap harinya
type Attendance struct {
	ID         string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID   string     `gorm:"type:uuid;not null;index"`
	EmployeeID string     `gorm:"type:uuid;not null;index"`
	Date       time.Time  `gorm:"type:date;not null"` // Tanggal Absen (tanpa jam)
	CheckIn    *time.Time // Pointer (*) agar bisa bernilai kosong (null) jika belum absen
	CheckOut   *time.Time // Pointer (*) agar bisa bernilai kosong (null) jika belum pulang
	Status     string     `gorm:"type:varchar(50);default:'on-time'"` // on-time, late
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relasi ke tabel Master Karyawan
	Employee Employee `gorm:"foreignKey:EmployeeID"`
}
