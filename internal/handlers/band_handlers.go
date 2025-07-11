package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BandMember struct {
	ID     int    `json:"id"`
	BandID int    `json:"band_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
}

type Band struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	UserID      int          `json:"user_id"`
	MemberCount int          `json:"member_count"`
	Members     []BandMember `json:"members,omitempty"`
}

type CreateBandRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Members     []BandMember `json:"members,omitempty"`
}

type UpdateBandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddMemberRequest struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type UpdateMemberRequest struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// In-memory storage (replace with database in production)
var bands []Band = make([]Band, 0)
var bandMembers []BandMember = make([]BandMember, 0)
var nextBandID = 1
var nextMemberID = 1

// GetBands returns all bands for the authenticated user
func GetBands(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var userBands []Band
	for _, band := range bands {
		if band.UserID == userID {
			// Get members for this band
			var bandMembersList []BandMember
			for _, member := range bandMembers {
				if member.BandID == band.ID {
					bandMembersList = append(bandMembersList, member)
				}
			}

			bandWithMembers := band
			bandWithMembers.Members = bandMembersList
			bandWithMembers.MemberCount = len(bandMembersList)
			userBands = append(userBands, bandWithMembers)
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
		MemberCount: 0, // Will be updated after adding members
	}

	nextBandID++
	bands = append(bands, newBand)

	// Add members if provided
	if len(req.Members) > 0 {
		for _, member := range req.Members {
			newMember := BandMember{
				ID:     nextMemberID,
				BandID: newBand.ID,
				Name:   member.Name,
				Role:   member.Role,
				Email:  member.Email,
				Phone:  member.Phone,
			}
			nextMemberID++
			bandMembers = append(bandMembers, newMember)
		}
		// Update member count
		newBand.MemberCount = len(req.Members)
	}

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
			// Get members for this band
			var bandMembersList []BandMember
			for _, member := range bandMembers {
				if member.BandID == band.ID {
					bandMembersList = append(bandMembersList, member)
				}
			}

			bandWithMembers := band
			bandWithMembers.Members = bandMembersList
			bandWithMembers.MemberCount = len(bandMembersList)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bandWithMembers)
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

// GetBandMembers returns all members of a specific band
func GetBandMembers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	// Check if user owns the band
	bandExists := false
	for _, band := range bands {
		if band.ID == bandID && band.UserID == userID {
			bandExists = true
			break
		}
	}

	if !bandExists {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	// Get members for this band
	var bandMembersList []BandMember
	for _, member := range bandMembers {
		if member.BandID == bandID {
			bandMembersList = append(bandMembersList, member)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bandMembersList)
}

// AddBandMember adds a new member to a band
func AddBandMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	bandIDStr := chi.URLParam(r, "bandId")
	bandID, err := strconv.Atoi(bandIDStr)
	if err != nil {
		http.Error(w, "Invalid band ID format", http.StatusBadRequest)
		return
	}

	// Check if user owns the band
	bandExists := false
	for _, band := range bands {
		if band.ID == bandID && band.UserID == userID {
			bandExists = true
			break
		}
	}

	if !bandExists {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	var req AddMemberRequest
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

	newMember := BandMember{
		ID:     nextMemberID,
		BandID: bandID,
		Name:   req.Name,
		Role:   req.Role,
		Email:  req.Email,
		Phone:  req.Phone,
	}

	nextMemberID++
	bandMembers = append(bandMembers, newMember)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMember)
}

// UpdateBandMember updates an existing band member
func UpdateBandMember(w http.ResponseWriter, r *http.Request) {
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

	// Check if user owns the band
	bandExists := false
	for _, band := range bands {
		if band.ID == bandID && band.UserID == userID {
			bandExists = true
			break
		}
	}

	if !bandExists {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	var req UpdateMemberRequest
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

	for i := range bandMembers {
		if bandMembers[i].ID == memberID && bandMembers[i].BandID == bandID {
			bandMembers[i].Name = req.Name
			bandMembers[i].Role = req.Role
			bandMembers[i].Email = req.Email
			bandMembers[i].Phone = req.Phone

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bandMembers[i])
			return
		}
	}

	http.Error(w, "Member not found", http.StatusNotFound)
}

// DeleteBandMember removes a member from a band
func DeleteBandMember(w http.ResponseWriter, r *http.Request) {
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

	// Check if user owns the band
	bandExists := false
	for _, band := range bands {
		if band.ID == bandID && band.UserID == userID {
			bandExists = true
			break
		}
	}

	if !bandExists {
		http.Error(w, "Band not found", http.StatusNotFound)
		return
	}

	for i := range bandMembers {
		if bandMembers[i].ID == memberID && bandMembers[i].BandID == bandID {
			bandMembers = append(bandMembers[:i], bandMembers[i+1:]...) // Remove the element
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Member not found", http.StatusNotFound)
}
