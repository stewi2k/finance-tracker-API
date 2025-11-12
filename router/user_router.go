package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/handlers"
	"github.com/stevenwijaya/finance-tracker/middleware"
)

func InitUserRouter(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	user := r.Group("/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", handlers.Profile)
	}
}
