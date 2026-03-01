package handlers

import (
	"net/http"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

// DevLogin provides a quick login for development/local access without Telegram or Firebase.
// It auto-creates a dev user if one doesn't exist, then returns a JWT token.
func DevLogin(c *gin.Context) {
	telegramService := services.NewTelegramService()

	var user models.User
	result := database.DB.Where("username = ?", "dev_user").First(&user)
	if result.Error != nil {
		// Create the dev user
		user = models.User{
			Username:     "dev_user",
			FirstName:    "Dev",
			LastName:     "User",
			LanguageCode: "en",
			IsVerified:   true,
			IsPremium:    true,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dev user"})
			return
		}
	}

	token, err := telegramService.GenerateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
