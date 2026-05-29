package models

import (
	"time"

	"gorm.io/gorm"
)

// BillOfMaterial adalah "Resep" untuk membuat sebuah Barang Jadi (Finished Good)
type BillOfMaterial struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"`
	Code      string `gorm:"type:varchar(50);not null"` // Contoh: BOM-MEJA-01
	ItemID    string `gorm:"type:uuid;not null"`        // ID Barang Jadi yang akan dibuat
	Name      string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relasi ke bahan-bahan baku
	Components []BOMComponent `gorm:"foreignKey:BOMID"`
}

// BOMComponent adalah rincian "Bahan Baku" yang dibutuhkan
type BOMComponent struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BOMID      string `gorm:"type:uuid;not null;index"`
	MaterialID string `gorm:"type:uuid;not null"` // ID Barang mentah (Paku, Kayu, dll)
	Quantity   int    `gorm:"not null"`           // Butuh berapa buah untuk 1 resep?
}

// ProductionOrder adalah Surat Perintah Produksi / Perakitan
type ProductionOrder struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID       string    `gorm:"type:uuid;not null;index"`
	OrderNumber    string    `gorm:"type:varchar(100);not null"`         // Contoh: PRD-2026-001
	BOMID          string    `gorm:"type:uuid;not null"`                 // Pakai resep yang mana?
	WarehouseID    string    `gorm:"type:uuid;not null"`                 // Produksi di gudang/pabrik mana?
	TargetQuantity int       `gorm:"not null"`                           // Mau buat berapa unit?
	Status         string    `gorm:"type:varchar(50);default:'planned'"` // planned, in-progress, completed
	StartDate      time.Time `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relasi untuk memudahkan penarikan data
	BOM BillOfMaterial `gorm:"foreignKey:BOMID"`
}
