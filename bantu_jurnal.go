package main

import (
	"fmt"
	"log"
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"
)

func main() {
	// 1. Konek ke Database
	config.ConnectDB()

	// 2. Ambil Perusahaan Pertama (Tenant)
	var tenant models.Tenant
	if err := config.DB.First(&tenant).Error; err != nil {
		log.Fatal("Belum ada perusahaan. Silakan jalankan API Register dulu!")
	}

	// 3. Cari atau Buat Akun Kas Bank (1110)
	var kas models.Account
	err := config.DB.Where("tenant_id = ? AND account_code = ?", tenant.ID, "1110").First(&kas).Error
	if err != nil {
		// Jika tidak ketemu/lupa, buatkan otomatis
		kas = models.Account{
			TenantID:    tenant.ID,
			AccountCode: "1110",
			AccountName: "Kas Bank BCA",
			AccountType: "asset",
		}
		config.DB.Create(&kas)
		fmt.Println("[+] Akun Kas (1110) otomatis dibuat karena belum ada.")
	}

	// 4. Cari atau Buat Akun Modal (3100)
	var modal models.Account
	err = config.DB.Where("tenant_id = ? AND account_code = ?", tenant.ID, "3100").First(&modal).Error
	if err != nil {
		// Jika tidak ketemu/lupa, buatkan otomatis
		modal = models.Account{
			TenantID:    tenant.ID,
			AccountCode: "3100",
			AccountName: "Modal Disetor",
			AccountType: "equity",
		}
		config.DB.Create(&modal)
		fmt.Println("[+] Akun Modal (3100) otomatis dibuat karena belum ada.")
	}

	// 5. Cetak JSON Siap Copy-Paste dengan tanggal hari ini
	tanggalHariIni := time.Now().Format("2006-01-02")
	fmt.Println("\n============= COPY JSON DI BAWAH INI =============")
	fmt.Printf(`{
    "reference": "TRX-MODAL-001",
    "date": "%s",
    "description": "Setoran Modal Awal dari Pemilik",
    "lines": [
        {
            "account_id": "%s",
            "debit": 50000000,
            "credit": 0
        },
        {
            "account_id": "%s",
            "debit": 0,
            "credit": 50000000
        }
    ]
}`, tanggalHariIni, kas.ID, modal.ID)
	fmt.Println("\n==================================================")
}
