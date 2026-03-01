package handlers

import (
	"net/http"

	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: services.NewAnalyticsService(),
	}
}

// Get godoc
// @Summary Get User Analytics Data
// @Description Fetches aggregated statistics for tasks, habits, pomodoros, and goals for the authenticated user.
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} models.AnalyticsDashboardData
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/analytics [get]
func (h *AnalyticsHandler) Get(c *gin.Context) {
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := userIdStr.(uint)

	data, err := h.analyticsService.GetUserAnalytics(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, data)
}
