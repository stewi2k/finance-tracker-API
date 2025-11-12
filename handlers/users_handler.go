package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stevenwijaya/finance-tracker/models"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"github.com/stevenwijaya/finance-tracker/services"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("Invalid input data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RegisterUser(&input); err != nil {
		log.Error("Failed to register user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	log.Info("User registered successfully:", input.Username)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("Invalid input data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.LoginUser(input.Username, input.Password)
	if err != nil {
		log.Error("Login failed: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	log.Info("User logged in successfully:", input.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// Profile handler untuk test JWT
func Profile(c *gin.Context) {
	userID := c.GetInt("user_id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Authenticated user profile",
		"user_id": userID,
	})
}
