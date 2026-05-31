package models

import "time"

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SKU       string    `json:"sku"`
	Type      string    `json:"type"` // Hanya boleh berisi "IN" (Masuk) atau "OUT" (Keluar)
	Qty       int       `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
}
