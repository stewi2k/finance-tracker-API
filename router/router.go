package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	//tambahkan router modular disini
	InitUserRouter(router)
	InitTransactionRouter(router)

	return router
}
