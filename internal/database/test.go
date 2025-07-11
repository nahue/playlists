package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// TestConnection tests the database connection and prints basic info
func TestConnection(db *sqlx.DB) error {
	// if db == nil {
	// 	return fmt.Errorf("database connection not established")
	// }

	// Test basic query
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to query database version: %w", err)
	}

	log.Printf("Database connection successful. PostgreSQL version: %s", version)

	// Test if migrations table exists
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'goose_db_version'
		)
	`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check migrations table: %w", err)
	}

	if tableExists {
		log.Println("Migrations table exists - migrations have been run")
	} else {
		log.Println("Migrations table not found - run migrations first")
	}

	return nil
}
