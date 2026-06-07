package models

import "time"

type SystemSettings struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyName string    `json:"company_name"`
	Logo        string    `json:"logo"` // Menyimpan Base64 atau URL
	UpdatedAt   time.Time `json:"updated_at"`
}
