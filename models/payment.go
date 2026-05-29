package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment adalah bukti pembayaran atau penerimaan uang
type Payment struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID      string    `gorm:"type:uuid;not null;index"`
	ReceiptNumber string    `gorm:"type:varchar(100);not null"` // Contoh: PAY-2026-001
	InvoiceID     string    `gorm:"type:uuid;not null;index"`   // Membayar tagihan yang mana?
	Amount        float64   `gorm:"type:decimal(15,2);not null"`
	PaymentDate   time.Time `gorm:"not null"`
	Method        string    `gorm:"type:varchar(50);default:'Bank Transfer'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Relasi ke Faktur
	Invoice Invoice `gorm:"foreignKey:InvoiceID"`
}
