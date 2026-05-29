package models

import (
	"time"

	"gorm.io/gorm"
)

// Approval adalah mesin pencatat alur persetujuan (Maker-Checker)
type Approval struct {
	ID           string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string  `gorm:"type:uuid;not null;index"`
	DocumentType string  `gorm:"type:varchar(100);not null"`         // Contoh: "PurchaseOrder", "Payroll", "Leave"
	DocumentID   string  `gorm:"type:varchar(100);not null"`         // ID atau Nomor Dokumen yang butuh persetujuan
	RequestedBy  string  `gorm:"type:uuid;not null"`                 // ID Staf yang mengajukan
	ApproverID   *string `gorm:"type:uuid"`                          // ID Manajer yang menyetujui (kosong saat diajukan)
	Status       string  `gorm:"type:varchar(50);default:'pending'"` // pending, approved, rejected
	Notes        string  `gorm:"type:text"`                          // Alasan jika ditolak atau catatan manajer
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
