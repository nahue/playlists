package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nahue/playlists/internal/database"
)

// BandHandler handles HTTP requests for band operations
type BandHandler struct {
	bandRepo *database.BandRepository
	logger   *log.Logger
}

// NewBandHandler creates a new BandHandler with the given repository
func NewBandHandler(bandRepo *database.BandRepository, logger *log.Logger) *BandHandler {
	return &BandHandler{
		bandRepo: bandRepo,
		logger:   logger,
	}
}

// GetBands returns all bands for the authenticated user
func (h *BandHandler) GetBands(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	bands, err := h.bandRepo.GetBandsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get bands", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bands)
}

// CreateBand creates a new band for the authenticated user
func (h *BandHandler) CreateBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var req database.CreateBandRequest
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

	band, err := h.bandRepo.CreateBand(userID, req)
	if err != nil {
		h.logger.Printf("Failed to create band: %v", err)
		http.Error(w, "Failed to create band", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(band)
}

// GetBand returns a specific band by ID (only if owned by the authenticated user)
func (h *BandHandler) GetBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	band, err := h.bandRepo.GetBandByID(id, userID)
	if err != nil {
		http.Error(w, "Failed to get band", http.StatusInternalServerError)
		return
	}

	if band == nil {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(band)
}

// UpdateBand updates a specific band (only if owned by the authenticated user)
func (h *BandHandler) UpdateBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req database.UpdateBandRequest
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

	band, err := h.bandRepo.UpdateBand(id, userID, req)
	if err != nil {
		http.Error(w, "Failed to update band", http.StatusInternalServerError)
		return
	}

	if band == nil {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(band)
}

// DeleteBand deletes a specific band (only if owned by the authenticated user)
func (h *BandHandler) DeleteBand(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.bandRepo.DeleteBand(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete band", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetBandMembers returns all members of a specific band
func (h *BandHandler) GetBandMembers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	// Check if user owns the band
	band, err := h.bandRepo.GetBandByID(bandID, userID)
	if err != nil {
		http.Error(w, "Failed to verify band ownership", http.StatusInternalServerError)
		return
	}

	if band == nil {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(band.Members)
}

// AddBandMember adds a new member to a band
func (h *BandHandler) AddBandMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	var req database.AddMemberRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Role == "" {
		http.Error(w, "Name and role are required", http.StatusBadRequest)
		return
	}

	member, err := h.bandRepo.AddBandMember(bandID, userID, req)
	if err != nil {
		http.Error(w, "Failed to add band member", http.StatusInternalServerError)
		return
	}

	if member == nil {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// UpdateBandMember updates a specific band member
func (h *BandHandler) UpdateBandMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	memberIDStr := chi.URLParam(r, "memberId")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	var req database.UpdateMemberRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Role == "" {
		http.Error(w, "Name and role are required", http.StatusBadRequest)
		return
	}

	member, err := h.bandRepo.UpdateBandMember(memberID, bandID, userID, req)
	if err != nil {
		http.Error(w, "Failed to update band member", http.StatusInternalServerError)
		return
	}

	if member == nil {
		http.Error(w, "Band member not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// DeleteBandMember deletes a specific band member
func (h *BandHandler) DeleteBandMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	memberIDStr := chi.URLParam(r, "memberId")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		http.Error(w, "Invalid member ID format", http.StatusBadRequest)
		return
	}

	err = h.bandRepo.DeleteBandMember(memberID, bandID, userID)
	if err != nil {
		http.Error(w, "Failed to delete band member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
