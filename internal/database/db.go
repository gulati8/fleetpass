package database

import (
	"fleetpass/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "fleetpass_user"),
		Password: getEnv("DB_PASSWORD", "fleetpass_password"),
		DBName:   getEnv("DB_NAME", "fleetpass"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect establishes a connection to the database using GORM
func Connect(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting underlying db: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * 60) // 5 minutes in seconds
	sqlDB.SetConnMaxIdleTime(5 * 60)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("error enabling UUID extension: %w", err)
	}

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.Organization{},
		&models.Location{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Vehicle{},
	)
	if err != nil {
		return fmt.Errorf("error running auto-migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Init initializes the global database connection
func Init() error {
	config := LoadConfigFromEnv()
	db, err := Connect(config)
	if err != nil {
		return err
	}

	DB = db

	// Run migrations
	if err := AutoMigrate(db); err != nil {
		return err
	}

	// Seed database with initial data
	if err := SeedDatabase(db); err != nil {
		log.Printf("Warning: Error seeding database: %v", err)
		// Don't fail on seed errors (might already be seeded)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
