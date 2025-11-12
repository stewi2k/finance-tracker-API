package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/config"
	"github.com/stevenwijaya/finance-tracker/database"
	"github.com/stevenwijaya/finance-tracker/router"
)

func initApp() {
	// Initialize configurations, database connections, etc.
	config.LoadConfig()
	database.ConnectDatabase()
}

func main() {
	initApp()
	router := router.InitRouter()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Finance Tracker API",
		})
	})
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "server is up and running",
		})
	})

	router.Run(":8081")
}
