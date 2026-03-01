package database

import (
	"fmt"
	"log"
	"os"

	"github.com/bif12/kyrotask/internal/config"
	"github.com/bif12/kyrotask/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	var dialector gorm.Dialector

	cfg := config.AppConfig

	switch cfg.DatabaseDriver {
	case "sqlite":
		// Create data directory if it doesn't exist
		os.MkdirAll("./data", 0755)
		dialector = sqlite.Open(cfg.DatabaseDSN)
	case "postgres":
		dialector = postgres.Open(cfg.DatabaseDSN)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.DatabaseDriver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Retrieve underlying sql.DB to apply connection pooling settings
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve sql.DB for pooling: %w", err)
	}

	// Apply database connection pooling to optimize resource footprint
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	log.Printf("Connected to %s database with pooling limits activated", cfg.DatabaseDriver)
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.TelegramAccount{},
		&models.Project{},
		&models.Task{},
		&models.Comment{},
		&models.Tag{},
		&models.Goal{},
		&models.Milestone{},
		&models.Habit{},
		&models.HabitLog{},
		&models.UserPomodoroProgress{},
		&models.PomodoroSession{},
		&models.Activity{},
		&models.Notification{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
