package users

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	models "github.com/stevenwijaya/finance-tracker/models/users"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"github.com/stevenwijaya/finance-tracker/pkg/response"
	services "github.com/stevenwijaya/finance-tracker/services/users"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("Invalid input data:", err)
		response.Error(c, http.StatusBadRequest, "invalid input format :"+err.Error())
		return
	}

	if err := services.RegisterUser(&input); err != nil {
		log.Error("Failed to register user:", err)
		response.Error(c, http.StatusInternalServerError, "Failed to register user : "+err.Error())
		return
	}

	log.Info("User registered successfully:", input.Username)
	response.Success(c, http.StatusOK, "User registered successfully", input)
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("Invalid input data:", err)
		response.Error(c, http.StatusBadGateway, "invalid input data : "+err.Error())
		return
	}

	user, err := services.LoginUser(input.Username, input.Password)
	if err != nil {
		log.Error("Login failed: ", err)
		response.Error(c, http.StatusUnauthorized, "Login Failed")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	log.Info("User logged in successfully:", input.Username)
	response.Success(c, http.StatusOK, "User logged in successfully", gin.H{
		"token": tokenString,
	})
}

// Profile handler untuk test JWT
func Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	response.Success(c, http.StatusOK, "Authenticated user profile", gin.H{
		"user_id": userID,
	})
}
