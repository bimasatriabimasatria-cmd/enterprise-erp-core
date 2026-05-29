package models

import (
	"time"

	"gorm.io/gorm"
)

// BankStatement adalah Kepala File Mutasi yang diunggah
type BankStatement struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string `gorm:"type:uuid;not null;index"`
	BankName  string `gorm:"type:varchar(100);not null"`
	Period    string `gorm:"type:varchar(50);not null"`
	Status    string `gorm:"type:varchar(50);default:'unreconciled'"` // unreconciled, reconciled
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Lines []BankStatementLine `gorm:"foreignKey:StatementID"`
}

// BankStatementLine adalah Baris Transaksi dari mutasi bank (Excel/CSV)
type BankStatementLine struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	StatementID  string    `gorm:"type:uuid;not null;index"`
	Date         time.Time `gorm:"not null"`
	Description  string    `gorm:"type:varchar(255);not null"`
	Amount       float64   `gorm:"type:decimal(15,2);not null"`
	IsReconciled bool      `gorm:"default:false"`
	PaymentID    *string   `gorm:"type:uuid"` // Kosong jika belum ketemu pasangannya di ERP
	CreatedAt    time.Time
}
