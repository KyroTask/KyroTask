package handlers

import (
	"net/http"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

type FirebaseAuthRequest struct {
	IDToken     string `json:"id_token" binding:"required"`
	LinkAccount bool   `json:"link_account"`
}

// VerifyFirebaseAuth verifies a Firebase token and returns a JWT
func (h *AuthHandler) VerifyFirebaseAuth(c *gin.Context) {
	var req FirebaseAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	fbToken, err := services.VerifyIDToken(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify Firebase token"})
		return
	}

	firebaseUID := fbToken.UID
	email := ""
	if emailClaim, ok := fbToken.Claims["email"].(string); ok {
		email = emailClaim
	}

	db := database.DB
	var user models.User

	if req.LinkAccount {
		userID, exists := middleware.GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Must be logged in to link account"})
			return
		}

		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Link account
		user.FirebaseUID = &firebaseUID
		if email != "" {
			user.Email = &email
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Firebase account"})
			return
		}

	} else {
		// Normal Login/Signup flow
		result := db.Where("firebase_uid = ?", firebaseUID).First(&user)
		if result.Error != nil {
			// Find by email if exist for some reason? Let's just create new if UID not found
			user = models.User{
				FirebaseUID:  &firebaseUID,
				Username:     "User_" + firebaseUID[:6],
				FirstName:    "New",
				LastName:     "User",
				LanguageCode: "en",
			}
			if email != "" {
				user.Email = &email
			}

			// extract name from token if available
			if nameClaim, ok := fbToken.Claims["name"].(string); ok {
				user.FirstName = nameClaim
				user.LastName = ""
			}
			if pictureClaim, ok := fbToken.Claims["picture"].(string); ok {
				user.PhotoURL = pictureClaim
			}

			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		}
	}

	// Generate standard JWT token for our app
	token, err := h.telegramService.GenerateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
