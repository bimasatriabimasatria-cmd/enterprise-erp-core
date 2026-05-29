package models

import (
	"time"

	"gorm.io/gorm"
)

// JournalEntry adalah Header Jurnal (Kapan transaksi terjadi dan apa buktinya)
type JournalEntry struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	Reference   string    `gorm:"type:varchar(100);not null"` // Contoh: INV-2026-001
	Date        time.Time `gorm:"not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relasi ke baris-baris jurnal (Satu Header punya banyak Baris)
	Lines []JournalLine `gorm:"foreignKey:JournalEntryID"`
}

// JournalLine adalah Baris Detail Jurnal (Debit/Kredit per akun)
type JournalLine struct {
	ID             string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JournalEntryID string  `gorm:"type:uuid;not null;index"`
	AccountID      string  `gorm:"type:uuid;not null;index"`
	Debit          float64 `gorm:"type:decimal(15,2);default:0"`
	Credit         float64 `gorm:"type:decimal(15,2);default:0"`

	// Relasi untuk menarik data nama akun nantinya
	Account Account `gorm:"foreignKey:AccountID"`
}
