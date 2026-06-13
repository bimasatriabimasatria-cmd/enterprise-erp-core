package controllers

import (
	"fmt"
	"net/http"
	"enterprise-erp-core/config"
	"enterprise-erp-core/models"
	"github.com/gin-gonic/gin"
)

// CreateJournal menangani standar Double-Entry Bookkeeping
func CreateJournal(c *gin.Context) {
	var input models.JournalEntry

	// 1. Validasi Input JSON dari Frontend/API
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	// 2. Ambil TenantID dari Middleware Keamanan
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Tenant tidak ditemukan"})
		return
	}
	input.TenantID = tenantID.(string)

	// 3. Validasi Akuntansi: Total Debit HARUS = Total Kredit
	var totalDebit, totalCredit float64
	for _, line := range input.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	// Toleransi kecil untuk floating point
	if totalDebit != totalCredit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jurnal tidak seimbang (Unbalanced): Total Debit tidak sama dengan Total Kredit"})
		return
	}

	// ==========================================
	// 4. ENTERPRISE PATTERN: DATABASE TRANSACTION
	// ==========================================
	tx := config.DB.Begin() // Memulai Transaksi Database

	// [PENTING] Menyuntikkan identitas Tenant ke Postgres untuk melewati RLS
	// Ini akan dibaca oleh fungsi `public.get_jwt_tenant_id()` yang kita buat di SQL
	rlsQuery := fmt.Sprintf(`SET LOCAL "request.jwt.claims" = '{"tenant_id": "%s"}'`, input.TenantID)
	if err := tx.Exec(rlsQuery).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menginisiasi keamanan RLS"})
		return
	}

	// 5. Simpan Data (GORM akan otomatis menyimpan Header dan Lines karena kita sudah atur relasinya di struct)
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback() // Jika ada 1 saja yang gagal (misal: AccountID tidak ada), BATALKAN SEMUA!
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jurnal: " + err.Error()})
		return
	}

	// 6. Jika semua sukses, baru kita "Sahkan" datanya ke database permanen
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jurnal Double-Entry berhasil disimpan!",
		"data":    input,
	})
}
