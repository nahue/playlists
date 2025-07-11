package database

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	user, err := repo.CreateUser(req)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john@example.com", user.Email)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	// Create first user
	_, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Try to create second user with same email
	_, err = repo.CreateUser(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
		Password:  "password123",
	}

	createdUser, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Get user by ID
	user, err := repo.GetUserByID(createdUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, "jane@example.com", user.Email)

	// Try to get non-existent user
	user, err = repo.GetUserByID(999)
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     "bob@example.com",
		Password:  "password123",
	}

	_, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Get user by email
	user, err := repo.GetUserByEmail("bob@example.com")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Bob", user.FirstName)
	assert.Equal(t, "Johnson", user.LastName)
	assert.Equal(t, "bob@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)

	// Try to get non-existent user
	user, err = repo.GetUserByEmail("nonexistent@example.com")
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_AuthenticateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "Alice",
		LastName:  "Brown",
		Email:     "alice@example.com",
		Password:  "password123",
	}

	_, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Test successful authentication
	loginReq := LoginRequest{
		Email:    "alice@example.com",
		Password: "password123",
	}

	user, err := repo.AuthenticateUser(loginReq)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Alice", user.FirstName)
	assert.Equal(t, "Brown", user.LastName)
	assert.Equal(t, "alice@example.com", user.Email)

	// Test wrong password
	loginReq.Password = "wrongpassword"
	user, err = repo.AuthenticateUser(loginReq)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	assert.Nil(t, user)

	// Test non-existent user
	loginReq.Email = "nonexistent@example.com"
	loginReq.Password = "password123"
	user, err = repo.AuthenticateUser(loginReq)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	assert.Nil(t, user)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Debug: Check what's in the database before creating user
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users")
	require.NoError(t, err)
	t.Logf("Users in database before test: %d", count)

	req := CreateUserRequest{
		FirstName: "Charlie",
		LastName:  "Wilson",
		Email:     "charlie_update@example.com",
		Password:  "password123",
	}

	createdUser, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Update user
	updateReq := UpdateUserRequest{
		FirstName: "Charles",
		LastName:  "Williams",
		Email:     "charles_updated@example.com",
	}

	updatedUser, err := repo.UpdateUser(createdUser.ID, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "Charles", updatedUser.FirstName)
	assert.Equal(t, "Williams", updatedUser.LastName)
	assert.Equal(t, "charles_updated@example.com", updatedUser.Email)

	// Try to update non-existent user
	updatedUser, err = repo.UpdateUser(999, updateReq)
	require.NoError(t, err)
	assert.Nil(t, updatedUser)
}

func TestUserRepository_UpdateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create first user
	req1 := CreateUserRequest{
		FirstName: "User1",
		LastName:  "Last1",
		Email:     "user1@example.com",
		Password:  "password123",
	}
	_, err := repo.CreateUser(req1)
	require.NoError(t, err)

	// Create second user
	req2 := CreateUserRequest{
		FirstName: "User2",
		LastName:  "Last2",
		Email:     "user2@example.com",
		Password:  "password123",
	}
	user2, err := repo.CreateUser(req2)
	require.NoError(t, err)

	// Try to update second user with first user's email
	updateReq := UpdateUserRequest{
		Email: "user1@example.com",
	}

	updatedUser, err := repo.UpdateUser(user2.ID, updateReq)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	assert.Nil(t, updatedUser)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "David",
		LastName:  "Miller",
		Email:     "david@example.com",
		Password:  "oldpassword",
	}

	createdUser, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Update password
	err = repo.UpdatePassword(createdUser.ID, "newpassword")
	require.NoError(t, err)

	// Verify new password works
	loginReq := LoginRequest{
		Email:    "david@example.com",
		Password: "newpassword",
	}

	user, err := repo.AuthenticateUser(loginReq)
	require.NoError(t, err)
	assert.NotNil(t, user)

	// Verify old password doesn't work
	loginReq.Password = "oldpassword"
	user, err = repo.AuthenticateUser(loginReq)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	req := CreateUserRequest{
		FirstName: "Eve",
		LastName:  "Davis",
		Email:     "eve@example.com",
		Password:  "password123",
	}

	createdUser, err := repo.CreateUser(req)
	require.NoError(t, err)

	// Verify user exists
	user, err := repo.GetUserByID(createdUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, user)

	// Delete user
	err = repo.DeleteUser(createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	user, err = repo.GetUserByID(createdUser.ID)
	require.NoError(t, err)
	assert.Nil(t, user)

	// Try to delete non-existent user
	err = repo.DeleteUser(999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create multiple users
	users := []CreateUserRequest{
		{FirstName: "User1", LastName: "Last1", Email: "user1@example.com", Password: "password123"},
		{FirstName: "User2", LastName: "Last2", Email: "user2@example.com", Password: "password123"},
		{FirstName: "User3", LastName: "Last3", Email: "user3@example.com", Password: "password123"},
	}

	for _, req := range users {
		_, err := repo.CreateUser(req)
		require.NoError(t, err)
	}

	// Get all users
	allUsers, err := repo.GetAllUsers()
	require.NoError(t, err)
	assert.Len(t, allUsers, 3)

	// Verify all three users exist (order doesn't matter)
	userNames := make([]string, len(allUsers))
	for i, user := range allUsers {
		userNames[i] = user.FirstName
	}
	assert.Contains(t, userNames, "User1")
	assert.Contains(t, userNames, "User2")
	assert.Contains(t, userNames, "User3")
}

func TestUserRepository_GetUsersCount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Initially should be 0
	count, err := repo.GetUsersCount()
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Create a user
	req := CreateUserRequest{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "password123",
	}

	_, err = repo.CreateUser(req)
	require.NoError(t, err)

	// Count should be 1
	count, err = repo.GetUsersCount()
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
