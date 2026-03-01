package handlers

import (
	"log"
	"net/http"

	"github.com/bif12/kyrotask/internal/config"
	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// PomodoroWS upgrades the HTTP connection to WebSocket for real-time pomodoro sync.
// Authentication is done via ?token=<JWT> query parameter (since WS can't use headers).
func PomodoroWS(c *gin.Context) {
	// Authenticate via JWT token in query param
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token query param required"})
		return
	}

	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := claims.UserID

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	Hub.Register(userID, conn)
	log.Printf("WS connected: user %d", userID)

	// Send initial pomodoro state immediately
	sendPomodoroStatus(userID, conn)

	// Read pump — just drains incoming messages and detects disconnect
	go func() {
		defer func() {
			Hub.Unregister(userID, conn)
			conn.Close()
			log.Printf("WS disconnected: user %d", userID)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// sendPomodoroStatus sends the current pomodoro status to a single connection.
func sendPomodoroStatus(userID uint, conn *websocket.Conn) {
	payload := buildPomodoroPayload(userID)
	conn.WriteJSON(payload)
}

// BroadcastPomodoroStatus pushes updated pomodoro status to ALL of a user's tabs/devices.
func BroadcastPomodoroStatus(userID uint) {
	payload := buildPomodoroPayload(userID)
	Hub.BroadcastToUser(userID, payload)
}

// buildPomodoroPayload creates the status JSON (same shape as GET /status).
func buildPomodoroPayload(userID uint) map[string]interface{} {
	var progress models.UserPomodoroProgress
	database.DB.Where("user_id = ?", userID).Attrs(models.UserPomodoroProgress{
		CurrentLevel: 1,
	}).FirstOrCreate(&progress)

	workDur, breakDur, targets, reqProj := models.GetLevelConfig(progress.CurrentLevel)

	var activeSession *models.PomodoroSession
	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err == nil {
		if session.CycleStartedAt == nil {
			session.CycleStartedAt = &session.CreatedAt
			database.DB.Save(&session)
		}
		activeSession = &session
	}

	payload := map[string]interface{}{
		"type":                    "pomodoro_status",
		"current_level":           progress.CurrentLevel,
		"required_work_duration":  workDur,
		"required_break_duration": breakDur,
		"target_cycles":           targets,
		"requires_project":        reqProj,
		"active_session":          activeSession,
	}

	if activeSession != nil {
		payload["seconds_remaining"] = activeSession.SecondsRemaining()
		payload["is_on_break"] = activeSession.IsOnBreak
		payload["is_paused"] = activeSession.IsPaused
	}

	return payload
}
