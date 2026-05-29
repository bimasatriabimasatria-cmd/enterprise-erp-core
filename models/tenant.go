package models

import (
	"time"
	"gorm.io/gorm"
)

// Tenant adalah model untuk Multi-Tenancy. 
// Setiap klien/perusahaan yang memakai ERP ini akan terdaftar di sini.
type Tenant struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Domain    string         `gorm:"type:varchar(255);unique;not null"` // cth: pt-abadi.erp.com
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}