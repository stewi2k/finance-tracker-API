package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Missing Authorization header
		if authHeader == "" {
			log.Warn("Unauthorized access attempt - missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Ambil token dari header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			log.Warn("Unauthorized access attempt - empty token after Bearer prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Empty token"})
			c.Abort()
			return
		}

		// Parse dan verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Error("JWTAuth - invalid signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Errorf("JWTAuth - failed to parse token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Warn("JWTAuth - token invalid or expired")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Validasi claim token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			// Cek apakah token expired
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					log.Warn("JWTAuth - token expired")
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
					c.Abort()
					return
				}
			}

			// Ambil user_id dari token dan konversi ke uint
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				log.Error("JWTAuth - invalid user_id type in token claims")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			userID := uint(userIDFloat)
			c.Set("user_id", userID)

			// Log sukses autentikasi user
			log.Infof("JWTAuth - user authenticated successfully (user_id=%d, path=%s, method=%s)", userID, c.Request.URL.Path, c.Request.Method)

			c.Next()
		} else {
			log.Error("JWTAuth - invalid token claims structure")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}
