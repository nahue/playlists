package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type PlaylistEntry struct {
	ID       int    `json:"id"`
	Artist   string `json:"artist"`
	Song     string `json:"song"`
	UserName string `json:"user_name"`
}

var playlist []PlaylistEntry = make([]PlaylistEntry, 0) // Initialize as an empty slice
var nextPlaylistID = 1

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

func AddToPlaylist(w http.ResponseWriter, r *http.Request) {
	var newEntry PlaylistEntry
	err := json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newEntry.Artist == "" || newEntry.Song == "" || newEntry.UserName == "" {
		http.Error(w, "Artist, song, and user name are required", http.StatusBadRequest)
		return
	}

	newEntry.ID = nextPlaylistID
	nextPlaylistID++
	playlist = append(playlist, newEntry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEntry)
}

func UpdatePlaylistEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedEntry PlaylistEntry
	err = json.NewDecoder(r.Body).Decode(&updatedEntry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if updatedEntry.Artist == "" || updatedEntry.Song == "" || updatedEntry.UserName == "" {
		http.Error(w, "Artist, song, and user name are required", http.StatusBadRequest)
		return
	}

	for i := range playlist {
		if playlist[i].ID == id {
			updatedEntry.ID = id // Ensure the ID matches
			playlist[i] = updatedEntry
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedEntry)
			return
		}
	}

	http.Error(w, "Playlist entry not found", http.StatusNotFound)
}

func DeletePlaylistEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for i := range playlist {
		if playlist[i].ID == id {
			playlist = append(playlist[:i], playlist[i+1:]...) // Remove the element
			w.WriteHeader(http.StatusNoContent)                // Indicate successful deletion
			return
		}
	}

	http.Error(w, "Playlist entry not found", http.StatusNotFound)
}

func GetPlaylistEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for _, entry := range playlist {
		if entry.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	http.Error(w, "Playlist entry not found", http.StatusNotFound)
}

// GetArtists handles the artist autocomplete endpoint
func GetArtists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		// Return empty array if no query provided
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
		return
	}

	query = strings.ToLower(strings.TrimSpace(query))

	// Create a map to store unique artists
	artistMap := make(map[string]bool)
	var artists []string

	// Search through playlist entries
	for _, entry := range playlist {
		artist := strings.ToLower(entry.Artist)
		if strings.Contains(artist, query) {
			// Use original case from the entry
			if !artistMap[entry.Artist] {
				artistMap[entry.Artist] = true
				artists = append(artists, entry.Artist)
			}
		}
	}

	// Limit results to 10 suggestions
	if len(artists) > 10 {
		artists = artists[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// GetUserNames handles the user name autocomplete endpoint
func GetUserNames(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		// Return empty array if no query provided
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
		return
	}

	query = strings.ToLower(strings.TrimSpace(query))

	// Create a map to store unique user names
	userMap := make(map[string]bool)
	var userNames []string

	// Search through playlist entries
	for _, entry := range playlist {
		userName := strings.ToLower(entry.UserName)
		if strings.Contains(userName, query) {
			// Use original case from the entry
			if !userMap[entry.UserName] {
				userMap[entry.UserName] = true
				userNames = append(userNames, entry.UserName)
			}
		}
	}

	// Limit results to 10 suggestions
	if len(userNames) > 10 {
		userNames = userNames[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userNames)
}
