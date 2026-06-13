package routes

import (
	"enterprise-erp-core/controllers"
	"enterprise-erp-core/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupJournalRoutes mengelola semua endpoint untuk modul FICO (Jurnal)
func SetupJournalRoutes(r *gin.Engine) {
	// Kita buat group "/api/v1/finance" untuk standar API modern
	financeGroup := r.Group("/api/v1/finance")
	
	// Terapkan Middleware Keamanan (RBAC & RLS Mocking) untuk seluruh rute di grup ini
	financeGroup.Use(middlewares.AuthMiddleware())
	{
		// Endpoint: POST /api/v1/finance/journals
		financeGroup.POST("/journals", controllers.CreateJournal)
	}
}
