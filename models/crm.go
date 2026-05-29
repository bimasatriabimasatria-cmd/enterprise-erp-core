package models

import (
	"time"

	"gorm.io/gorm"
)

// Lead adalah Calon Pelanggan (Prospek)
type Lead struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"`
	Name      string `gorm:"type:varchar(255);not null"` // Nama PIC
	Company   string `gorm:"type:varchar(255);not null"` // Nama Perusahaan
	Email     string `gorm:"type:varchar(255)"`
	Phone     string `gorm:"type:varchar(50)"`
	Status    string `gorm:"type:varchar(50);default:'new'"` // new, contacted, qualified, converted, lost
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Customer adalah Pelanggan Resmi (Master Data Pelanggan)
type Customer struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"`
	Name      string `gorm:"type:varchar(255);not null"` // Nama Perusahaan / Individu
	Contact   string `gorm:"type:varchar(255)"`          // Nama PIC
	Email     string `gorm:"type:varchar(255)"`
	Phone     string `gorm:"type:varchar(50)"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
