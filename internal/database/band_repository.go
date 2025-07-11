package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Band represents a band in the database
type Band struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	UserID      int       `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// BandMember represents a band member in the database
type BandMember struct {
	ID        int       `db:"id" json:"id"`
	BandID    int       `db:"band_id" json:"band_id"`
	Name      string    `db:"name" json:"name"`
	Role      string    `db:"role" json:"role"`
	Email     string    `db:"email" json:"email,omitempty"`
	Phone     string    `db:"phone" json:"phone,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// BandWithMembers represents a band with its members
type BandWithMembers struct {
	Band
	Members     []BandMember `json:"members,omitempty"`
	MemberCount int          `json:"member_count"`
}

// CreateBandRequest represents the request to create a new band
type CreateBandRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Members     []BandMember `json:"members,omitempty"`
}

// UpdateBandRequest represents the request to update a band
type UpdateBandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddMemberRequest represents the request to add a new member
type AddMemberRequest struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// UpdateMemberRequest represents the request to update a member
type UpdateMemberRequest struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// BandRepository handles database operations for bands
type BandRepository struct {
	db *sqlx.DB
}

// NewBandRepository creates a new band repository
func NewBandRepository(db *sqlx.DB) *BandRepository {
	return &BandRepository{db: db}
}

// GetBandsByUserID returns all bands for a specific user
func (r *BandRepository) GetBandsByUserID(userID int) ([]BandWithMembers, error) {
	query := `
		SELECT b.id, b.name, b.description, b.user_id, b.created_at, b.updated_at
		FROM bands b
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	var bands []Band
	err := r.db.Select(&bands, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bands: %w", err)
	}

	// Get members for each band
	var bandsWithMembers []BandWithMembers
	for _, band := range bands {
		members, err := r.GetBandMembers(band.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get members for band %d: %w", band.ID, err)
		}

		bandWithMembers := BandWithMembers{
			Band:        band,
			Members:     members,
			MemberCount: len(members),
		}
		bandsWithMembers = append(bandsWithMembers, bandWithMembers)
	}

	return bandsWithMembers, nil
}

// GetBandByID returns a specific band by ID (only if owned by the user)
func (r *BandRepository) GetBandByID(bandID, userID int) (*BandWithMembers, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at
		FROM bands
		WHERE id = $1 AND user_id = $2
	`

	var band Band
	err := r.db.Get(&band, query, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to get band: %w", err)
	}

	// Get members for this band
	members, err := r.GetBandMembers(band.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members for band %d: %w", band.ID, err)
	}

	bandWithMembers := &BandWithMembers{
		Band:        band,
		Members:     members,
		MemberCount: len(members),
	}

	return bandWithMembers, nil
}

// CreateBand creates a new band with optional members
func (r *BandRepository) CreateBand(userID int, req CreateBandRequest) (*BandWithMembers, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create the band
	bandQuery := `
		INSERT INTO bands (name, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, user_id, created_at, updated_at
	`

	var band Band
	err = tx.Get(&band, bandQuery, req.Name, req.Description, userID)
	if err != nil {
		return nil, err
	}

	// Add members if provided
	var members []BandMember
	if len(req.Members) > 0 {
		memberQuery := `
			INSERT INTO band_members (band_id, name, role, email, phone)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, band_id, name, role, email, phone, created_at, updated_at
		`

		for _, memberReq := range req.Members {
			var member BandMember
			err = tx.Get(&member, memberQuery, band.ID, memberReq.Name, memberReq.Role, memberReq.Email, memberReq.Phone)
			if err != nil {
				return nil, err
			}
			members = append(members, member)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	bandWithMembers := &BandWithMembers{
		Band:        band,
		Members:     members,
		MemberCount: len(members),
	}

	return bandWithMembers, nil
}

// UpdateBand updates a specific band
func (r *BandRepository) UpdateBand(bandID, userID int, req UpdateBandRequest) (*Band, error) {
	query := `
		UPDATE bands
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, name, description, user_id, created_at, updated_at
	`

	var band Band
	err := r.db.Get(&band, query, req.Name, req.Description, bandID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Band not found
		}
		return nil, fmt.Errorf("failed to update band: %w", err)
	}

	return &band, nil
}

// DeleteBand deletes a specific band and all its members
func (r *BandRepository) DeleteBand(bandID, userID int) error {
	query := `
		DELETE FROM bands
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(query, bandID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete band: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Band not found
	}

	return nil
}

// GetBandMembers returns all members of a specific band
func (r *BandRepository) GetBandMembers(bandID int) ([]BandMember, error) {
	query := `
		SELECT id, band_id, name, role, email, phone, created_at, updated_at
		FROM band_members
		WHERE band_id = $1
		ORDER BY created_at ASC
	`

	var members []BandMember
	err := r.db.Select(&members, query, bandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get band members: %w", err)
	}

	return members, nil
}

// AddBandMember adds a new member to a band
func (r *BandRepository) AddBandMember(bandID, userID int, req AddMemberRequest) (*BandMember, error) {
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

	// Add the member
	memberQuery := `
		INSERT INTO band_members (band_id, name, role, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, band_id, name, role, email, phone, created_at, updated_at
	`

	var member BandMember
	err = r.db.Get(&member, memberQuery, bandID, req.Name, req.Role, req.Email, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to add band member: %w", err)
	}

	return &member, nil
}

// UpdateBandMember updates a specific band member
func (r *BandRepository) UpdateBandMember(memberID, bandID, userID int, req UpdateMemberRequest) (*BandMember, error) {
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

	// Update the member
	memberQuery := `
		UPDATE band_members
		SET name = $1, role = $2, email = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND band_id = $6
		RETURNING id, band_id, name, role, email, phone, created_at, updated_at
	`

	var member BandMember
	err = r.db.Get(&member, memberQuery, req.Name, req.Role, req.Email, req.Phone, memberID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Member not found
		}
		return nil, fmt.Errorf("failed to update band member: %w", err)
	}

	return &member, nil
}

// DeleteBandMember deletes a specific band member
func (r *BandRepository) DeleteBandMember(memberID, bandID, userID int) error {
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

	// Delete the member
	memberQuery := `DELETE FROM band_members WHERE id = $1 AND band_id = $2`
	result, err := r.db.Exec(memberQuery, memberID, bandID)
	if err != nil {
		return fmt.Errorf("failed to delete band member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Member not found
	}

	return nil
}

// GetBandMemberByID returns a specific band member by ID
func (r *BandRepository) GetBandMemberByID(memberID, bandID, userID int) (*BandMember, error) {
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

	// Get the member
	memberQuery := `
		SELECT id, band_id, name, role, email, phone, created_at, updated_at
		FROM band_members
		WHERE id = $1 AND band_id = $2
	`

	var member BandMember
	err = r.db.Get(&member, memberQuery, memberID, bandID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Member not found
		}
		return nil, fmt.Errorf("failed to get band member: %w", err)
	}

	return &member, nil
}
