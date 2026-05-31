package models

import "time"

type Employee struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	EmployeeID string    `json:"employee_id" gorm:"unique"` // Contoh: EMP-001
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	Status     string    `json:"status"` // "Aktif" atau "Cuti" atau "Resign"
	CreatedAt  time.Time `json:"created_at"`
}
