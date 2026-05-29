package models

import (
	"time"
)

// AuditLog adalah rekaman "CCTV" dari semua aktivitas penting di sistem
type AuditLog struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID  string    `gorm:"type:uuid;not null;index"`
	UserID    string    `gorm:"type:uuid;not null;index"`   // Siapa pelakunya?
	Action    string    `gorm:"type:varchar(50);not null"`  // Contoh: POST (Tambah), PUT (Edit), DELETE (Hapus)
	Resource  string    `gorm:"type:varchar(255);not null"` // Di menu mana? (cth: /api/invoices)
	Payload   string    `gorm:"type:text"`                  // Data apa yang dia ketik/kirim?
	IPAddress string    `gorm:"type:varchar(50)"`           // Dari IP/Komputer mana?
	CreatedAt time.Time `gorm:"index"`                      // Kapan terjadinya?
}
