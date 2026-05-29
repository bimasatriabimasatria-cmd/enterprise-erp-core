package models

import (
	"time"

	"gorm.io/gorm"
)

// Warehouse adalah lokasi fisik gudang
type Warehouse struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"`
	Code      string `gorm:"type:varchar(50);not null"`  // Contoh: GDG-JKT
	Name      string `gorm:"type:varchar(255);not null"` // Contoh: Gudang Pusat Jakarta
	Location  string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Inventory adalah detail stok barang per lokasi gudang
type Inventory struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string `gorm:"type:uuid;not null;index"`
	WarehouseID string `gorm:"type:uuid;not null;index"`
	ItemID      string `gorm:"type:uuid;not null;index"`
	Quantity    int    `gorm:"not null;default:0"` // Berapa jumlah barang ini di gudang ini?
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Warehouse Warehouse `gorm:"foreignKey:WarehouseID"`
	Item      Item      `gorm:"foreignKey:ItemID"`
}

// StockTransfer adalah riwayat mutasi/perpindahan barang
type StockTransfer struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID      string    `gorm:"type:uuid;not null;index"`
	Reference     string    `gorm:"type:varchar(100);not null"` // Contoh: TRF-001
	SourceID      string    `gorm:"type:uuid;not null"`         // Gudang Asal
	DestinationID string    `gorm:"type:uuid;not null"`         // Gudang Tujuan
	ItemID        string    `gorm:"type:uuid;not null"`
	Quantity      int       `gorm:"not null"`
	TransferDate  time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
