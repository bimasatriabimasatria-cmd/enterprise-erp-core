package models

import "time"

// JournalEntry adalah Header dari transaksi (Map/Amplop)
type JournalEntry struct {
	ID              string        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID        string        `gorm:"type:uuid;not null" json:"tenant_id"` // Kunci Multi-Tenant
	EntryDate       time.Time     `gorm:"type:date;not null" json:"entry_date"`
	ReferenceNumber string        `gorm:"type:varchar(100)" json:"reference_number"`
	Description     string        `gorm:"type:text" json:"description"`
	Status          string        `gorm:"type:varchar(20);default:'DRAFT'" json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	
	// Relasi ke Journal Lines (1 Jurnal punya banyak baris Debit/Kredit)
	Lines           []JournalLine `gorm:"foreignKey:JournalEntryID;constraint:OnDelete:CASCADE;" json:"lines"`
}

// JournalLine adalah Detail Transaksi (Isi dari Amplop)
type JournalLine struct {
	ID             string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	JournalEntryID string  `gorm:"type:uuid;not null" json:"journal_entry_id"`
	AccountID      string  `gorm:"type:uuid;not null" json:"account_id"`
	Description    string  `gorm:"type:text" json:"description"`
	Debit          float64 `gorm:"type:decimal(19,4);default:0.0000" json:"debit"`
	Credit         float64 `gorm:"type:decimal(19,4);default:0.0000" json:"credit"`
}
