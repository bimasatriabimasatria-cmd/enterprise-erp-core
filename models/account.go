package models

import (
	"time"

	"gorm.io/gorm"
)

// Account adalah cetak biru untuk Chart of Accounts (Bagan Akun)
type Account struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string `gorm:"type:uuid;not null;uniqueIndex:idx_tenant_account"`
	AccountCode string `gorm:"type:varchar(50);not null;uniqueIndex:idx_tenant_account"` // Contoh: 1100, 2100, 5000
	AccountName string `gorm:"type:varchar(255);not null"`                               // Contoh: Kas Kecil, Hutang Usaha
	AccountType string `gorm:"type:varchar(50);not null"`                                // asset, liability, equity, revenue, expense
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relasi
	Tenant Tenant `gorm:"foreignKey:TenantID"`
}
