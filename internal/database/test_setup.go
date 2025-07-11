package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// setupTestDB ensures the test database is migrated and returns a connection
func setupTestDB(t *testing.T) *sqlx.DB {
	// Test database configuration
	config := &Config{
		Host:     "localhost",
		Port:     "5455",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	// Connect to database
	db, err := Open(config)
	require.NoError(t, err)

	// Clean up tables before each test
	cleanupTables(t, db)

	return db
}

// findProjectRoot finds the project root directory by looking for go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

// cleanupTables removes all data from tables in the correct order
func cleanupTables(t *testing.T, db *sqlx.DB) {
	// Delete in reverse order due to foreign key constraints
	db.MustExec("DELETE FROM band_members")
	db.MustExec("DELETE FROM bands")
	db.MustExec("DELETE FROM playlist_entries")
	db.MustExec("DELETE FROM users")

	// Verify cleanup worked
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users")
	require.NoError(t, err)
	require.Equal(t, 0, count, "Users table should be empty after cleanup")
}

// createTestUser creates a test user and returns the user ID
func createTestUser(t *testing.T, db *sqlx.DB, email string) int {
	var userID int
	err := db.Get(&userID, `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, "Test", "User", email, "hashed_password")
	require.NoError(t, err)
	return userID
}

// TestMain runs once before all tests to set up the test environment
func TestMain(m *testing.M) {
	// Run migrations before any tests
	config := &Config{
		Host:     "localhost",
		Port:     "5455",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	// Get the migrations directory path
	projectRoot, err := findProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Failed to find project root: %v", err))
	}
	migrationsDir := filepath.Join(projectRoot, "migrations")

	// Build the goose command
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	cmd := exec.Command("goose", "-dir", migrationsDir, "postgres", dsn, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run migrations
	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("Failed to run database migrations: %v", err))
	}

	// Run the tests
	os.Exit(m.Run())
}
