package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Band struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	MemberCount int    `json:"member_count"`
}

type CreateBandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateBandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage (replace with database in production)
var bands []Band = make([]Band, 0)
var nextBandID = 1

// GetBands returns all bands for the authenticated user
func GetBands(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var userBands []Band
	for _, band := range bands {
		if band.UserID == userID {
			userBands = append(userBands, band)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userBands)
}

// CreateBand creates a new band for the authenticated user
func CreateBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var req CreateBandRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Band name is required", http.StatusBadRequest)
		return
	}

	newBand := Band{
		ID:          nextBandID,
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
		MemberCount: 0, // Start with 0 members
	}

	nextBandID++
	bands = append(bands, newBand)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBand)
}

// GetBand returns a specific band by ID (only if owned by the authenticated user)
func GetBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for _, band := range bands {
		if band.ID == id && band.UserID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(band)
			return
		}
	}

	http.Error(w, "Band not found", http.StatusNotFound)
}

// UpdateBand updates a specific band (only if owned by the authenticated user)
func UpdateBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req UpdateBandRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Band name is required", http.StatusBadRequest)
		return
	}

	for i := range bands {
		if bands[i].ID == id && bands[i].UserID == userID {
			bands[i].Name = req.Name
			bands[i].Description = req.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bands[i])
			return
		}
	}

	http.Error(w, "Band not found", http.StatusNotFound)
}

// DeleteBand deletes a specific band (only if owned by the authenticated user)
func DeleteBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for i := range bands {
		if bands[i].ID == id && bands[i].UserID == userID {
			bands = append(bands[:i], bands[i+1:]...) // Remove the element
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Band not found", http.StatusNotFound)
}
