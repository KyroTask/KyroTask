package main

import (
	"log"

	"github.com/bif12/kyrotask/internal/config"
	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/handlers"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	config.Load()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Firebase Auth
	services.InitFirebaseAuth()

	// Start Scheduler
	scheduler := services.NewScheduler()
	scheduler.Start()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// GZIP compression with custom size exclusion
	// gin-contrib/gzip compressest by default, but we can bypass it for tiny payloads using an Ignore logic wrapper if we wanted,
	// or we can use the library's built-in exclusion functions.
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4"}), gzip.WithExcludedPaths([]string{"/api/v1/health"})))

	// Global Security Headers
	r.Use(middleware.SecurityHeaders())

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// API v1 routes
	// Limit API to 10 requests per second with a burst of 50
	limiter := middleware.NewRateLimiter(rate.Limit(10), 50)

	v1 := r.Group("/api/v1")
	v1.Use(limiter.Limit())
	{
		// Public routes
		v1.GET("/health", handlers.HealthCheck)

		// Handlers
		authHandler := handlers.NewAuthHandler()
		telegramHandler := handlers.NewTelegramHandler()
		projectHandler := handlers.NewProjectHandler()
		taskHandler := handlers.NewTaskHandler()
		goalHandler := handlers.NewGoalHandler()
		habitHandler := handlers.NewHabitHandler()
		activityHandler := handlers.NewActivityHandler()
		dashboardHandler := handlers.NewDashboardHandler()
		calendarHandler := handlers.NewCalendarHandler()
		milestoneHandler := handlers.NewMilestoneHandler()
		commentHandler := handlers.NewCommentHandler()
		tagHandler := handlers.NewTagHandler()
		userHandler := handlers.NewUserHandler()
		pomodoroHandler := handlers.NewPomodoroHandler()
		analyticsHandler := handlers.NewAnalyticsHandler()

		// Start Telegram Polling in background
		go telegramHandler.StartPolling()

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/telegram/verify", telegramHandler.VerifyTelegramWebApp)
			auth.POST("/telegram/widget", telegramHandler.VerifyTelegramWidget)
			auth.POST("/firebase", authHandler.VerifyFirebaseAuth)
			auth.POST("/dev-login", handlers.DevLogin)
		}

		// Telegram bot webhook
		v1.POST("/telegram/webhook", telegramHandler.HandleWebhook)

		// WebSocket — auth handled inside handler via token query param
		v1.GET("/ws/pomodoro", handlers.PomodoroWS)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// User
			protected.GET("/user/me", authHandler.GetMe)
			protected.PUT("/user/me", userHandler.UpdateProfile)

			// Dashboard & Calendar
			protected.GET("/dashboard", dashboardHandler.Get)
			protected.GET("/calendar", calendarHandler.Get)

			// Pomodoro
			pomodoro := protected.Group("/pomodoro")
			{
				pomodoro.GET("/status", pomodoroHandler.GetStatus)
				pomodoro.POST("/session/start", pomodoroHandler.StartSession)
				pomodoro.POST("/session/cycle", pomodoroHandler.LogCycle)
				pomodoro.POST("/session/complete", pomodoroHandler.CompleteLevel)
				pomodoro.POST("/session/abandon", pomodoroHandler.AbandonSession)
				pomodoro.POST("/session/pause", pomodoroHandler.PauseSession)
				pomodoro.POST("/session/resume", pomodoroHandler.ResumeSession)
				pomodoro.POST("/session/resume-work", pomodoroHandler.ResumeWork)
			}

			// Projects
			protected.GET("/projects", projectHandler.List)
			protected.POST("/projects", projectHandler.Create)
			protected.GET("/projects/:id", projectHandler.Get)
			protected.PUT("/projects/:id", projectHandler.Update)
			protected.DELETE("/projects/:id", projectHandler.Delete)

			// Tasks
			protected.GET("/tasks", taskHandler.List)
			protected.POST("/tasks", taskHandler.Create)
			protected.GET("/tasks/:id", taskHandler.Get)
			protected.PUT("/tasks/:id", taskHandler.Update)
			protected.DELETE("/tasks/:id", taskHandler.Delete)

			// Task Comments
			protected.GET("/tasks/:id/comments", commentHandler.ListByTask)
			protected.POST("/tasks/:id/comments", commentHandler.Create)
			protected.DELETE("/comments/:comment_id", commentHandler.Delete)

			// Goals
			protected.GET("/goals", goalHandler.List)
			protected.GET("/goals/:id", goalHandler.Get)
			protected.POST("/goals", goalHandler.Create)
			protected.PUT("/goals/:id", goalHandler.Update)
			protected.DELETE("/goals/:id", goalHandler.Delete)

			// Milestones
			protected.GET("/milestones", milestoneHandler.List)
			protected.POST("/milestones", milestoneHandler.Create)
			protected.PUT("/milestones/:id", milestoneHandler.Update)
			protected.DELETE("/milestones/:id", milestoneHandler.Delete)

			// Habits
			protected.GET("/habits", habitHandler.List)
			protected.POST("/habits", habitHandler.Create)
			protected.GET("/habits/:id", habitHandler.Get)
			protected.PUT("/habits/:id", habitHandler.Update)
			protected.DELETE("/habits/:id", habitHandler.Delete)
			protected.POST("/habits/:id/log", habitHandler.Log)

			// Tags
			protected.GET("/tags", tagHandler.List)
			protected.POST("/tags", tagHandler.Create)
			protected.DELETE("/tags/:id", tagHandler.Delete)

			// Activity
			protected.GET("/activities", activityHandler.List)

			// Analytics
			protected.GET("/analytics", analyticsHandler.Get)
		}
	}

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
