package controllers

import (
	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
)

// Dapatkan Laporan Laba Rugi (Income Statement)
func GetIncomeStatement(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// 1. Tarik seluruh jurnal keuangan milik perusahaan ini
	var entries []models.JournalEntry
	if err := config.DB.Preload("Lines").Preload("Lines.Account").
		Where("tenant_id = ?", tenantID).Find(&entries).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menarik data buku besar"})
	}

	// 2. Siapkan wadah untuk menghitung total
	var totalRevenue float64 = 0
	var totalExpense float64 = 0

	// Kita gunakan Map (kamus) untuk mengelompokkan total per akun
	accountTotals := make(map[string]float64)
	accountNames := make(map[string]string)
	accountTypes := make(map[string]string)

	// 3. Proses setiap baris jurnal
	for _, entry := range entries {
		for _, line := range entry.Lines {
			accCode := line.Account.AccountCode

			// Simpan nama dan tipe akun untuk laporan
			if _, exists := accountNames[accCode]; !exists {
				accountNames[accCode] = line.Account.AccountName
				accountTypes[accCode] = line.Account.AccountType
			}

			// Rumus Akuntansi:
			// PENDAPATAN (Revenue) bertambah di Kredit, berkurang di Debit
			if line.Account.AccountType == "revenue" {
				accountTotals[accCode] += (line.Credit - line.Debit)
			}

			// BEBAN (Expense) bertambah di Debit, berkurang di Kredit
			if line.Account.AccountType == "expense" {
				accountTotals[accCode] += (line.Debit - line.Credit)
			}
		}
	}

	// 4. Susun bentuk laporannya
	var revenueDetails []fiber.Map
	var expenseDetails []fiber.Map

	for code, total := range accountTotals {
		if accountTypes[code] == "revenue" {
			totalRevenue += total
			revenueDetails = append(revenueDetails, fiber.Map{
				"account_code": code,
				"account_name": accountNames[code],
				"total":        total,
			})
		} else if accountTypes[code] == "expense" {
			totalExpense += total
			expenseDetails = append(expenseDetails, fiber.Map{
				"account_code": code,
				"account_name": accountNames[code],
				"total":        total,
			})
		}
	}

	// 5. Hitung Laba Bersih
	netIncome := totalRevenue - totalExpense

	// 6. Kirim Laporan ke Direktur!
	return c.JSON(fiber.Map{
		"report_name": "Laporan Laba Rugi (Income Statement)",
		"currency":    "IDR",
		"data": fiber.Map{
			"revenues": fiber.Map{
				"details":       revenueDetails,
				"total_revenue": totalRevenue,
			},
			"expenses": fiber.Map{
				"details":       expenseDetails,
				"total_expense": totalExpense,
			},
			"net_income": netIncome,
		},
	})
}
