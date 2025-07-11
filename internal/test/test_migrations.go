package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationsApplied(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Check that all expected tables exist
	tables := []string{"users", "bands", "band_members", "playlist_entries"}

	for _, table := range tables {
		var exists bool
		err := db.Get(&exists, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`, table)
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", table)
	}

	// Check that users table has the expected columns
	var columnCount int
	err := db.Get(&columnCount, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'users' 
		AND table_schema = 'public'
	`)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, columnCount, 6, "Users table should have at least 6 columns (id, first_name, last_name, email, password_hash, created_at, updated_at)")

	// Check that bands table has the expected columns
	err = db.Get(&columnCount, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'bands' 
		AND table_schema = 'public'
	`)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, columnCount, 5, "Bands table should have at least 5 columns (id, name, description, user_id, created_at, updated_at)")

	// Check that band_members table has the expected columns
	err = db.Get(&columnCount, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'band_members' 
		AND table_schema = 'public'
	`)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, columnCount, 6, "Band_members table should have at least 6 columns (id, band_id, name, role, email, phone, created_at, updated_at)")

	// Check that playlist_entries table has the expected columns
	err = db.Get(&columnCount, `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'playlist_entries' 
		AND table_schema = 'public'
	`)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, columnCount, 6, "Playlist_entries table should have at least 6 columns (id, artist, song, user_name, user_id, created_at, updated_at)")
}
