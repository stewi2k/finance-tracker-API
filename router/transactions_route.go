package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/stevenwijaya/finance-tracker/handlers/transactions"
	"github.com/stevenwijaya/finance-tracker/middleware"
)

func InitTransactionRouter(r *gin.Engine) {
	transaction := r.Group("/transaction")
	transaction.Use(middleware.JWTAuth())
	{
		transaction.POST("/", handlers.CreateTransaction)
		transaction.GET("/", handlers.GetAllTransaction)
		transaction.GET("/:id", handlers.GetTransactionById)
		transaction.GET("/summary", handlers.GetTransactionSummary)
		transaction.GET("/summary/category", handlers.GetTransactionSummaryByCategory)
		transaction.PUT("/:id", handlers.UpdateTransaction)
		transaction.DELETE("/:id", handlers.DeleteTransaction)
	}
}
