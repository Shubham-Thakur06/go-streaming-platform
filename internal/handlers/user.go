package handlers

import (
	"net/http"
	"time"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:  db,
		cfg: cfg,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !service.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := service.GenerateToken(&user, h.cfg.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
