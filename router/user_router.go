package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/handlers/users"
	"github.com/stevenwijaya/finance-tracker/middleware"
)

func InitUserRouter(r *gin.Engine) {
	r.POST("/register", users.Register)
	r.POST("/login", users.Login)

	user := r.Group("/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", users.Profile)
	}
}
