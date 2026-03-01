package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	DatabaseDriver          string
	DatabaseDSN             string
	JWTSecret               string
	JWTExpiry               time.Duration
	TelegramBotToken        string
	TelegramBotUsername     string
	FirebaseCredentialsPath string
	AllowedOrigins          []string
	GinMode                 string

	// Connection Pooling
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

var AppConfig *Config

func Load() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port:                    getEnv("PORT", "3001"),
		DatabaseDriver:          getEnv("DB_DRIVER", "sqlite"),
		DatabaseDSN:             getEnv("DB_DSN", "./data/dev.db"),
		JWTSecret:               getEnv("JWT_SECRET", "change-this-secret-key"),
		TelegramBotToken:        getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramBotUsername:     getEnv("TELEGRAM_BOT_USERNAME", "TaskBeastBot"),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		GinMode:                 getEnv("GIN_MODE", "debug"),

		// Connection Pooling Defaults
		DBMaxOpenConns:    25,
		DBMaxIdleConns:    10,
		DBConnMaxLifetime: 5 * time.Minute,
	}

	// Parse JWT expiry
	expiryStr := getEnv("JWT_EXPIRY", "168h")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Printf("Invalid JWT_EXPIRY format, using default 168h: %v", err)
		expiry = 168 * time.Hour
	}
	AppConfig.JWTExpiry = expiry

	// Parse allowed origins
	originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3001")
	AppConfig.AllowedOrigins = parseCommaSeparated(originsStr)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCommaSeparated(s string) []string {
	var result []string
	for i, j := 0, 0; j <= len(s); j++ {
		if j == len(s) || s[j] == ',' {
			if j > i {
				result = append(result, s[i:j])
			}
			i = j + 1
		}
	}
	return result
}
