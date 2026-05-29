package models

import (
	"time"
	"gorm.io/gorm"
)

// Item adalah cetak biru untuk Master Barang/Produk
type Item struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string         `gorm:"type:uuid;not null;uniqueIndex:idx_tenant_sku"` // Kunci Multi-Tenant
	SKU         string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_sku"` // Stock Keeping Unit
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"type:decimal(15,2);default:0"` // Decimal untuk uang (akuntansi)
	Stock       int            `gorm:"default:0"`
	IsActive    bool           `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relasi Opsional
	Tenant      Tenant         `gorm:"foreignKey:TenantID"`
}