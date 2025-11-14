package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/middleware"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	//tambahkan router modular disini
	InitUserRouter(router)
	InitTransactionRouter(router)

	return router
}
