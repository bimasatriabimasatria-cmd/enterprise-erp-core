package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware mensimulasikan login dan menyuntikkan konteks keamanan
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Di masa depan (Fase 2), di sini kita akan memvalidasi JWT Token (Bearer Auth).
		
		// UNTUK SEKARANG: Kita paksa sistem mengira bahwa user yang me-request 
		// berasal dari PT ERP Maju Bersama (Sesuai UUID yang Anda berikan)
		c.Set("tenant_id", "620c6b50-8e5d-4da8-8c5c-60c7c4052079")
		
		// Lanjut ke request berikutnya
		c.Next()
	}
}
