package models

import (
	"time"

	"gorm.io/gorm"
)

// Invoice adalah Header Faktur Penjualan
type Invoice struct {
	ID            string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID      string         `gorm:"type:uuid;not null;index"`
	InvoiceNumber string         `gorm:"type:varchar(100);not null"` // Contoh: INV-001
	CustomerName  string         `gorm:"type:varchar(255);not null"`
	Date          time.Time      `gorm:"not null"`
	TotalAmount   float64        `gorm:"type:decimal(15,2);default:0"` // Total harga semua barang
	Status        string         `gorm:"type:varchar(50);default:'unpaid'"` // unpaid, paid, canceled
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Relasi ke baris detail faktur
	Lines []InvoiceLine `gorm:"foreignKey:InvoiceID"`
}

// InvoiceLine adalah detail barang yang dibeli di dalam faktur tersebut
type InvoiceLine struct {
	ID         string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InvoiceID  string  `gorm:"type:uuid;not null;index"`
	ItemID     string  `gorm:"type:uuid;not null;index"` // Menyambung ke Master Barang
	Quantity   int     `gorm:"not null;default:1"`
	UnitPrice  float64 `gorm:"type:decimal(15,2);not null"`
	SubTotal   float64 `gorm:"type:decimal(15,2);not null"`

	// Relasi untuk menarik data barang
	Item Item `gorm:"foreignKey:ItemID"`
}