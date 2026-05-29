package models

import (
	"time"

	"gorm.io/gorm"
)

// PurchaseOrder adalah surat pesanan perusahaan ke Supplier
type PurchaseOrder struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string    `gorm:"type:uuid;not null;index"`
	PONumber     string    `gorm:"type:varchar(100);not null"`
	SupplierName string    `gorm:"type:varchar(255);not null"`
	Date         time.Time `gorm:"not null"`
	TotalAmount  float64   `gorm:"type:decimal(15,2);default:0"`
	Status       string    `gorm:"type:varchar(50);default:'draft'"` // draft (dipesan), received (diterima)
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relasi ke detail barang yang dipesan
	Lines []PurchaseOrderLine `gorm:"foreignKey:PurchaseOrderID"`
}

// PurchaseOrderLine adalah detail barang dalam PO
type PurchaseOrderLine struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchaseOrderID string  `gorm:"type:uuid;not null;index"`
	ItemID          string  `gorm:"type:uuid;not null;index"` // Sambung ke Master Barang
	Quantity        int     `gorm:"not null;default:1"`
	UnitPrice       float64 `gorm:"type:decimal(15,2);not null"` // Harga beli dari Supplier
	SubTotal        float64 `gorm:"type:decimal(15,2);not null"`

	Item Item `gorm:"foreignKey:ItemID"`
}
