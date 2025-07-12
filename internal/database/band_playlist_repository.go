package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// BandPlaylist represents a playlist for a specific band
type BandPlaylist struct {
	ID          int       `db:"id" json:"id"`
	BandID      int       `db:"band_id" json:"band_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// BandPlaylistSong represents a song in a band playlist
type BandPlaylistSong struct {
	ID         int       `db:"id" json:"id"`
	PlaylistID int       `db:"playlist_id" json:"playlist_id"`
	Artist     string    `db:"artist" json:"artist"`
	Song       string    `db:"song" json:"song"`
	Notes      string    `db:"notes" json:"notes"`
	Position   int       `db:"position" json:"position"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// BandPlaylistWithSongs represents a playlist with its songs
type BandPlaylistWithSongs struct {
	BandPlaylist
	Songs     []BandPlaylistSong `json:"songs"`
	SongCount int                `json:"song_count"`
}

// CreatePlaylistRequest represents the request to create a new playlist
type CreatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdatePlaylistRequest represents the request to update a playlist
type UpdatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddSongRequest represents the request to add a new song
type AddSongRequest struct {
	Artist   string `json:"artist"`
	Song     string `json:"song"`
	Notes    string `json:"notes"`
	Position int    `json:"position"`
}

// UpdateSongRequest represents the request to update a song
type UpdateSongRequest struct {
	Artist   string `json:"artist"`
	Song     string `json:"song"`
	Notes    string `json:"notes"`
	Position int    `json:"position"`
}

// BandPlaylistRepository handles database operations for band playlists
type BandPlaylistRepository struct {
	db *sqlx.DB
}

// NewBandPlaylistRepository creates a new band playlist repository
func NewBandPlaylistRepository(db *sqlx.DB) *BandPlaylistRepository {
	return &BandPlaylistRepository{db: db}
}

// GetPlaylistsByBandID returns all playlists for a specific band
func (r *BandPlaylistRepository) GetPlaylistsByBandID(bandID, userID int) ([]BandPlaylistWithSongs, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		SELECT id, band_id, name, description, created_at, updated_at
		FROM band_playlists
		WHERE band_id = $1
		ORDER BY created_at DESC
	`

	var playlists []BandPlaylist
	err = r.db.Select(&playlists, query, bandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	// Get songs for each playlist
	var playlistsWithSongs []BandPlaylistWithSongs
	for _, playlist := range playlists {
		songs, err := r.GetPlaylistSongs(playlist.ID, bandID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get songs for playlist %d: %w", playlist.ID, err)
		}

		playlistWithSongs := BandPlaylistWithSongs{
			BandPlaylist: playlist,
			Songs:        songs,
			SongCount:    len(songs),
		}
		playlistsWithSongs = append(playlistsWithSongs, playlistWithSongs)
	}

	return playlistsWithSongs, nil
}

// GetPlaylistByID returns a specific playlist by ID (only if owned by the user)
func (r *BandPlaylistRepository) GetPlaylistByID(playlistID, bandID, userID int) (*BandPlaylistWithSongs, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		SELECT id, band_id, name, description, created_at, updated_at
		FROM band_playlists
		WHERE id = $1 AND band_id = $2
	`

	var playlist BandPlaylist
	err = r.db.Get(&playlist, query, playlistID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Playlist not found
		}
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	// Get songs for this playlist
	songs, err := r.GetPlaylistSongs(playlist.ID, bandID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get songs for playlist %d: %w", playlist.ID, err)
	}

	playlistWithSongs := &BandPlaylistWithSongs{
		BandPlaylist: playlist,
		Songs:        songs,
		SongCount:    len(songs),
	}

	return playlistWithSongs, nil
}

// CreatePlaylist creates a new playlist for a band
func (r *BandPlaylistRepository) CreatePlaylist(bandID, userID int, req CreatePlaylistRequest) (*BandPlaylistWithSongs, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		INSERT INTO band_playlists (band_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, band_id, name, description, created_at, updated_at
	`

	var playlist BandPlaylist
	err = r.db.Get(&playlist, query, bandID, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	playlistWithSongs := &BandPlaylistWithSongs{
		BandPlaylist: playlist,
		Songs:        []BandPlaylistSong{},
		SongCount:    0,
	}

	return playlistWithSongs, nil
}

// UpdatePlaylist updates a specific playlist
func (r *BandPlaylistRepository) UpdatePlaylist(playlistID, bandID, userID int, req UpdatePlaylistRequest) (*BandPlaylist, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		UPDATE band_playlists
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND band_id = $4
		RETURNING id, band_id, name, description, created_at, updated_at
	`

	var playlist BandPlaylist
	err = r.db.Get(&playlist, query, req.Name, req.Description, playlistID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Playlist not found
		}
		return nil, fmt.Errorf("failed to update playlist: %w", err)
	}

	return &playlist, nil
}

// DeletePlaylist deletes a specific playlist and all its songs
func (r *BandPlaylistRepository) DeletePlaylist(playlistID, bandID, userID int) error {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Band not found
		}
		return fmt.Errorf("failed to verify band ownership: %w", err)
	}

	// Delete the playlist (songs will be deleted automatically due to CASCADE)
	query := `DELETE FROM band_playlists WHERE id = $1 AND band_id = $2`
	result, err := r.db.Exec(query, playlistID, bandID)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Playlist not found
	}

	return nil
}

// GetPlaylistSongs returns all songs for a specific playlist
func (r *BandPlaylistRepository) GetPlaylistSongs(playlistID, bandID, userID int) ([]BandPlaylistSong, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		SELECT s.id, s.playlist_id, s.artist, s.song, s.notes, s.position, s.created_at, s.updated_at
		FROM band_playlist_songs s
		JOIN band_playlists p ON s.playlist_id = p.id
		WHERE s.playlist_id = $1 AND p.band_id = $2
		ORDER BY s.position ASC, s.created_at ASC
	`

	var songs []BandPlaylistSong
	err = r.db.Select(&songs, query, playlistID, bandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist songs: %w", err)
	}

	return songs, nil
}

// AddSong adds a new song to a playlist
func (r *BandPlaylistRepository) AddSong(playlistID, bandID, userID int, req AddSongRequest) (*BandPlaylistSong, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	// Verify that the playlist belongs to the band
	playlistQuery := `SELECT id FROM band_playlists WHERE id = $1 AND band_id = $2`
	var playlistIDCheck int
	err = r.db.Get(&playlistIDCheck, playlistQuery, playlistID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Playlist not found
		}
		return nil, fmt.Errorf("failed to verify playlist ownership: %w", err)
	}

	query := `
		INSERT INTO band_playlist_songs (playlist_id, artist, song, notes, position)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, playlist_id, artist, song, notes, position, created_at, updated_at
	`

	var song BandPlaylistSong
	err = r.db.Get(&song, query, playlistID, req.Artist, req.Song, req.Notes, req.Position)
	if err != nil {
		return nil, fmt.Errorf("failed to add song: %w", err)
	}

	return &song, nil
}

// UpdateSong updates a specific song in a playlist
func (r *BandPlaylistRepository) UpdateSong(songID, playlistID, bandID, userID int, req UpdateSongRequest) (*BandPlaylistSong, error) {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to verify band ownership: %w", err)
	}

	// Verify that the song belongs to the playlist and the playlist belongs to the band
	query := `
		UPDATE band_playlist_songs
		SET artist = $1, song = $2, notes = $3, position = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND playlist_id = $6 AND playlist_id IN (
			SELECT id FROM band_playlists WHERE band_id = $7
		)
		RETURNING id, playlist_id, artist, song, notes, position, created_at, updated_at
	`

	var song BandPlaylistSong
	err = r.db.Get(&song, query, req.Artist, req.Song, req.Notes, req.Position, songID, playlistID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Song not found
		}
		return nil, fmt.Errorf("failed to update song: %w", err)
	}

	return &song, nil
}

// DeleteSong deletes a specific song from a playlist
func (r *BandPlaylistRepository) DeleteSong(songID, playlistID, bandID, userID int) error {
	// First verify that the band belongs to the user
	bandQuery := `SELECT id FROM bands WHERE id = $1 AND user_id = $2`
	var bandIDCheck int
	err := r.db.Get(&bandIDCheck, bandQuery, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Band not found
		}
		return fmt.Errorf("failed to verify band ownership: %w", err)
	}

	query := `
		DELETE FROM band_playlist_songs
		WHERE id = $1 AND playlist_id = $2 AND playlist_id IN (
			SELECT id FROM band_playlists WHERE band_id = $3
		)
	`

	result, err := r.db.Exec(query, songID, playlistID, bandID)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Song not found
	}

	return nil
}
