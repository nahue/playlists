package app

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/nahue/playlists/internal/database"
)

// Application represents the main application instance
type Application struct {
	Logger *log.Logger
	Config *Config
	DB     *sqlx.DB
}

// Config holds application configuration
type Config struct {
	Port string
	Host string
}

// NewConfig creates a new application config from environment variables
func NewConfig() *Config {
	return &Config{
		Port: getEnv("SERVER_PORT", "8080"),
		Host: getEnv("SERVER_HOST", ""),
	}
}

// New creates a new Application instance
func NewApplication() *Application {
	// Load environment variables from config file
	if err := godotenv.Load("config.env"); err != nil {
		log.Printf("Warning: Could not load config.env file: %v", err)
	}

	// Initialize database connection
	dbConfig := database.NewConfig()
	db, err := database.Open(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test database connection
	if err := database.TestConnection(db); err != nil {
		log.Printf("Database connection test failed: %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Application{
		Logger: logger,
		Config: NewConfig(),
		DB:     db,
	}
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	// Close database connection
	if err := database.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
