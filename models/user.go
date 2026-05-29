package models

import (
	"time"

	"gorm.io/gorm"
)

// User adalah model untuk pengguna aplikasi (Karyawan, Admin, Direktur)
type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"` // Mengikat user ke perusahaannya
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);unique;not null;index"`
	Password  string `gorm:"type:varchar(255);not null"`       // Akan dienkripsi
	Role      string `gorm:"type:varchar(50);default:'staff'"` // admin, manager, staff
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relasi ke tabel Tenant (Opsional tapi direkomendasikan untuk ORM)
	Tenant Tenant `gorm:"foreignKey:TenantID"`
}
