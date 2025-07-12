package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nahue/playlists/internal/database"
)

// BandPlaylistHandler handles HTTP requests for band playlist operations
type BandPlaylistHandler struct {
	playlistRepo *database.BandPlaylistRepository
	logger       *log.Logger
}

// NewBandPlaylistHandler creates a new BandPlaylistHandler with the given repository
func NewBandPlaylistHandler(playlistRepo *database.BandPlaylistRepository, logger *log.Logger) *BandPlaylistHandler {
	return &BandPlaylistHandler{
		playlistRepo: playlistRepo,
		logger:       logger,
	}
}

// GetPlaylists returns all playlists for a specific band
func (h *BandPlaylistHandler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlists, err := h.playlistRepo.GetPlaylistsByBandID(bandID, userID)
	if err != nil {
		h.logger.Printf("Failed to get playlists: %v", err)
		http.Error(w, "Failed to get playlists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

// GetPlaylist returns a specific playlist by ID
func (h *BandPlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	playlist, err := h.playlistRepo.GetPlaylistByID(playlistID, bandID, userID)
	if err != nil {
		h.logger.Printf("Failed to get playlist: %v", err)
		http.Error(w, "Failed to get playlist", http.StatusInternalServerError)
		return
	}

	if playlist == nil {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

// CreatePlaylist creates a new playlist for a band
func (h *BandPlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	var req database.CreatePlaylistRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Playlist name is required", http.StatusBadRequest)
		return
	}

	playlist, err := h.playlistRepo.CreatePlaylist(bandID, userID, req)
	if err != nil {
		h.logger.Printf("Failed to create playlist: %v", err)
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	if playlist == nil {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// UpdatePlaylist updates a specific playlist
func (h *BandPlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	var req database.UpdatePlaylistRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Playlist name is required", http.StatusBadRequest)
		return
	}

	playlist, err := h.playlistRepo.UpdatePlaylist(playlistID, bandID, userID, req)
	if err != nil {
		h.logger.Printf("Failed to update playlist: %v", err)
		http.Error(w, "Failed to update playlist", http.StatusInternalServerError)
		return
	}

	if playlist == nil {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

// DeletePlaylist deletes a specific playlist
func (h *BandPlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	err = h.playlistRepo.DeletePlaylist(playlistID, bandID, userID)
	if err != nil {
		h.logger.Printf("Failed to delete playlist: %v", err)
		http.Error(w, "Failed to delete playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPlaylistSongs returns all songs for a specific playlist
func (h *BandPlaylistHandler) GetPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	songs, err := h.playlistRepo.GetPlaylistSongs(playlistID, bandID, userID)
	if err != nil {
		h.logger.Printf("Failed to get playlist songs: %v", err)
		http.Error(w, "Failed to get playlist songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// AddSong adds a new song to a playlist
func (h *BandPlaylistHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	var req database.AddSongRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Artist == "" || req.Song == "" {
		http.Error(w, "Artist and song are required", http.StatusBadRequest)
		return
	}

	song, err := h.playlistRepo.AddSong(playlistID, bandID, userID, req)
	if err != nil {
		h.logger.Printf("Failed to add song: %v", err)
		http.Error(w, "Failed to add song", http.StatusInternalServerError)
		return
	}

	if song == nil {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// UpdateSong updates a specific song in a playlist
func (h *BandPlaylistHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	songIDStr := chi.URLParam(r, "songId")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID format", http.StatusBadRequest)
		return
	}

	var req database.UpdateSongRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Artist == "" || req.Song == "" {
		http.Error(w, "Artist and song are required", http.StatusBadRequest)
		return
	}

	song, err := h.playlistRepo.UpdateSong(songID, playlistID, bandID, userID, req)
	if err != nil {
		h.logger.Printf("Failed to update song: %v", err)
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	if song == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// DeleteSong deletes a specific song from a playlist
func (h *BandPlaylistHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	playlistIDStr := chi.URLParam(r, "playlistId")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID format", http.StatusBadRequest)
		return
	}

	songIDStr := chi.URLParam(r, "songId")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID format", http.StatusBadRequest)
		return
	}

	err = h.playlistRepo.DeleteSong(songID, playlistID, bandID, userID)
	if err != nil {
		h.logger.Printf("Failed to delete song: %v", err)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
