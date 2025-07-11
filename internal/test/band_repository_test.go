package test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/nahue/playlists/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB, cleanupTables, and createTestUser are now defined in test_setup.go

func TestBandRepository_CreateBand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID := createTestUser(t, db, "test@example.com")

	req := database.CreateBandRequest{
		Name:        "Test Band",
		Description: "A test band",
		Members: []database.BandMember{
			{Name: "John Doe", Role: "Guitarist", Email: "john@example.com"},
			{Name: "Jane Smith", Role: "Singer", Email: "jane@example.com"},
		},
	}

	band, err := repo.CreateBand(userID, req)
	require.NoError(t, err)
	assert.NotNil(t, band)
	assert.Equal(t, "Test Band", band.Name)
	assert.Equal(t, "A test band", band.Description)
	assert.Equal(t, userID, band.UserID)
	assert.Equal(t, 2, band.MemberCount)
	assert.Len(t, band.Members, 2)
	assert.Equal(t, "John Doe", band.Members[0].Name)
	assert.Equal(t, "Guitarist", band.Members[0].Role)
	assert.Equal(t, "Jane Smith", band.Members[1].Name)
	assert.Equal(t, "Singer", band.Members[1].Role)
}

func TestBandRepository_GetBandsByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID1 := createTestUser(t, db, "user1@example.com")
	userID2 := createTestUser(t, db, "user2@example.com")

	// Create bands for user1
	req1 := database.CreateBandRequest{Name: "Band 1", Description: "First band"}
	req2 := database.CreateBandRequest{Name: "Band 2", Description: "Second band"}

	_, err := repo.CreateBand(userID1, req1)
	require.NoError(t, err)
	_, err = repo.CreateBand(userID1, req2)
	require.NoError(t, err)

	// Create band for user2
	req3 := database.CreateBandRequest{Name: "Band 3", Description: "Third band"}
	_, err = repo.CreateBand(userID2, req3)
	require.NoError(t, err)

	// Get bands for user1
	bands, err := repo.GetBandsByUserID(userID1)
	require.NoError(t, err)
	assert.Len(t, bands, 2)

	// Check that both bands exist (order doesn't matter)
	bandNames := make([]string, len(bands))
	for i, band := range bands {
		bandNames[i] = band.Name
	}
	assert.Contains(t, bandNames, "Band 1")
	assert.Contains(t, bandNames, "Band 2")

	// Get bands for user2
	bands2, err := repo.GetBandsByUserID(userID2)
	require.NoError(t, err)
	assert.Len(t, bands2, 1)
	assert.Equal(t, "Band 3", bands2[0].Name)
}

func TestBandRepository_GetBandByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID1 := createTestUser(t, db, "user1@example.com")
	userID2 := createTestUser(t, db, "user2@example.com")

	req := database.CreateBandRequest{
		Name:        "Test Band",
		Description: "A test band",
		Members: []database.BandMember{
			{Name: "John Doe", Role: "Guitarist"},
		},
	}

	createdBand, err := repo.CreateBand(userID1, req)
	require.NoError(t, err)

	// Get band by ID (correct user)
	band, err := repo.GetBandByID(createdBand.ID, userID1)
	require.NoError(t, err)
	assert.NotNil(t, band)
	assert.Equal(t, "Test Band", band.Name)
	assert.Equal(t, 1, band.MemberCount)

	// Try to get band with wrong user
	band, err = repo.GetBandByID(createdBand.ID, userID2)
	require.NoError(t, err)
	assert.Nil(t, band) // Should return nil for unauthorized access

	// Try to get non-existent band
	band, err = repo.GetBandByID(999, userID1)
	require.NoError(t, err)
	assert.Nil(t, band)
}

func TestBandRepository_UpdateBand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID1 := createTestUser(t, db, "user1@example.com")
	userID2 := createTestUser(t, db, "user2@example.com")

	req := database.CreateBandRequest{Name: "Original Name", Description: "Original description"}
	createdBand, err := repo.CreateBand(userID1, req)
	require.NoError(t, err)

	// Update band
	updateReq := database.UpdateBandRequest{
		Name:        "Updated Name",
		Description: "Updated description",
	}

	updatedBand, err := repo.UpdateBand(createdBand.ID, userID1, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, updatedBand)
	assert.Equal(t, "Updated Name", updatedBand.Name)
	assert.Equal(t, "Updated description", updatedBand.Description)

	// Try to update with wrong user
	updatedBand, err = repo.UpdateBand(createdBand.ID, userID2, updateReq)
	require.NoError(t, err)
	assert.Nil(t, updatedBand) // Should return nil for unauthorized access
}

func TestBandRepository_DeleteBand(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID1 := createTestUser(t, db, "user1@example.com")
	userID2 := createTestUser(t, db, "user2@example.com")

	req := database.CreateBandRequest{
		Name: "Test Band",
		Members: []database.BandMember{
			{Name: "John Doe", Role: "Guitarist"},
		},
	}

	createdBand, err := repo.CreateBand(userID1, req)
	require.NoError(t, err)

	// Verify band exists
	band, err := repo.GetBandByID(createdBand.ID, userID1)
	require.NoError(t, err)
	assert.NotNil(t, band)

	// Delete band
	err = repo.DeleteBand(createdBand.ID, userID1)
	require.NoError(t, err)

	// Verify band is deleted
	band, err = repo.GetBandByID(createdBand.ID, userID1)
	require.NoError(t, err)
	assert.Nil(t, band)

	// Try to delete with wrong user (should not error but not delete)
	err = repo.DeleteBand(createdBand.ID, userID2)
	require.NoError(t, err)
}

func TestBandRepository_AddBandMember(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID := createTestUser(t, db, "test@example.com")

	req := database.CreateBandRequest{Name: "Test Band"}
	createdBand, err := repo.CreateBand(userID, req)
	require.NoError(t, err)

	// Add member
	memberReq := database.AddMemberRequest{
		Name:  "New Member",
		Role:  "Drummer",
		Email: "drummer@example.com",
		Phone: "123-456-7890",
	}

	member, err := repo.AddBandMember(createdBand.ID, userID, memberReq)
	require.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "New Member", member.Name)
	assert.Equal(t, "Drummer", member.Role)
	assert.Equal(t, "drummer@example.com", member.Email)
	assert.Equal(t, "123-456-7890", member.Phone)

	// Verify member was added
	band, err := repo.GetBandByID(createdBand.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, 1, band.MemberCount)
	assert.Equal(t, "New Member", band.Members[0].Name)
}

func TestBandRepository_UpdateBandMember(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID := createTestUser(t, db, "test@example.com")

	req := database.CreateBandRequest{
		Name: "Test Band",
		Members: []database.BandMember{
			{Name: "Original Name", Role: "Original Role"},
		},
	}

	createdBand, err := repo.CreateBand(userID, req)
	require.NoError(t, err)
	require.Len(t, createdBand.Members, 1)

	memberID := createdBand.Members[0].ID

	// Update member
	updateReq := database.UpdateMemberRequest{
		Name:  "Updated Name",
		Role:  "Updated Role",
		Email: "updated@example.com",
		Phone: "987-654-3210",
	}

	updatedMember, err := repo.UpdateBandMember(memberID, createdBand.ID, userID, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, updatedMember)
	assert.Equal(t, "Updated Name", updatedMember.Name)
	assert.Equal(t, "Updated Role", updatedMember.Role)
	assert.Equal(t, "updated@example.com", updatedMember.Email)
	assert.Equal(t, "987-654-3210", updatedMember.Phone)
}

func TestBandRepository_DeleteBandMember(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID := createTestUser(t, db, "test@example.com")

	req := database.CreateBandRequest{
		Name: "Test Band",
		Members: []database.BandMember{
			{Name: "Member 1", Role: "Role 1"},
			{Name: "Member 2", Role: "Role 2"},
		},
	}

	createdBand, err := repo.CreateBand(userID, req)
	require.NoError(t, err)
	require.Len(t, createdBand.Members, 2)

	memberID := createdBand.Members[0].ID

	// Delete member
	err = repo.DeleteBandMember(memberID, createdBand.ID, userID)
	require.NoError(t, err)

	// Verify member was deleted
	band, err := repo.GetBandByID(createdBand.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, 1, band.MemberCount)
	assert.Equal(t, "Member 2", band.Members[0].Name)
}

func TestBandRepository_GetBandMemberByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewBandRepository(db)
	userID := createTestUser(t, db, "test@example.com")

	req := database.CreateBandRequest{
		Name: "Test Band",
		Members: []database.BandMember{
			{Name: "Test Member", Role: "Test Role"},
		},
	}

	createdBand, err := repo.CreateBand(userID, req)
	require.NoError(t, err)
	require.Len(t, createdBand.Members, 1)

	memberID := createdBand.Members[0].ID

	// Get member by ID
	member, err := repo.GetBandMemberByID(memberID, createdBand.ID, userID)
	require.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "Test Member", member.Name)
	assert.Equal(t, "Test Role", member.Role)

	// Try to get member with wrong user
	member, err = repo.GetBandMemberByID(memberID, createdBand.ID, 999)
	require.NoError(t, err)
	assert.Nil(t, member) // Should return nil for unauthorized access
}
