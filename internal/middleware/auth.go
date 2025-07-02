package middleware

import (
	"net/http"
	"strings"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(cfg config.JWTConfig, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		claims, err := service.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		var user models.User
		if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserNotFound.Error()})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*models.User)
}
