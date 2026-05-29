package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Nama fungsi WAJIB diawali dengan kata "Test"
func TestHealthCheckAPI(t *testing.T) {
	// 1. Siapkan mesin Fiber simulasi (Tanpa perlu database)
	app := fiber.New()

	// 2. Daftarkan rute yang akan diuji
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success", "message": "OK"})
	})

	// 3. Buat tembakan Request HTTP palsu (Mock Request)
	req := httptest.NewRequest("GET", "/api/health", nil)

	// 4. Tembakkan ke dalam sistem (Sistem akan memprosesnya dalam memori)
	resp, err := app.Test(req)

	// 5. ASURANSI KODE: Pastikan tidak ada error dan statusnya WAJIB 200 (OK)
	assert.NoError(t, err, "Tidak boleh ada error saat melakukan request")
	assert.Equal(t, 200, resp.StatusCode, "Sistem harus membalas dengan status 200 OK")
}
