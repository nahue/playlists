# Database Package

This package provides database connectivity and repository patterns for the playlists application.

## Overview

The database package contains:
- Database connection management
- Repository implementations for data access
- Database models and request/response types

## Database Connection

### Configuration

The database connection is configured using environment variables:

```bash
DB_HOST=localhost
DB_PORT=5454
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
```

### Usage

```go
import "github.com/nahue/playlists/internal/database"

// Create database configuration
config := database.NewConfig()

// Open database connection
db, err := database.Open(config)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

## Repositories

### Band Repository

The `BandRepository` provides a clean interface for managing bands and band members in PostgreSQL.

#### Features

- **CRUD Operations**: Create, read, update, and delete bands
- **Member Management**: Add, update, and delete band members
- **User Isolation**: All operations are scoped to the authenticated user
- **Transaction Support**: Atomic operations for band creation with members
- **Error Handling**: Comprehensive error handling with meaningful messages

#### Usage

```go
import "github.com/nahue/playlists/internal/database"

// Create repository instance
repo := database.NewBandRepository(db)

// Create a new band with members
req := database.CreateBandRequest{
    Name:        "My Band",
    Description: "A great band",
    Members: []database.BandMember{
        {Name: "John Doe", Role: "Guitarist", Email: "john@example.com"},
        {Name: "Jane Smith", Role: "Singer", Email: "jane@example.com"},
    },
}

band, err := repo.CreateBand(userID, req)
if err != nil {
    log.Printf("Failed to create band: %v", err)
    return
}
```

#### Available Methods

##### Band Operations

- `GetBandsByUserID(userID int) ([]BandWithMembers, error)` - Get all bands for a user
- `GetBandByID(bandID, userID int) (*BandWithMembers, error)` - Get a specific band
- `CreateBand(userID int, req CreateBandRequest) (*BandWithMembers, error)` - Create a new band
- `UpdateBand(bandID, userID int, req UpdateBandRequest) (*Band, error)` - Update a band
- `DeleteBand(bandID, userID int) error` - Delete a band

##### Member Operations

- `GetBandMembers(bandID int) ([]BandMember, error)` - Get all members of a band
- `GetBandMemberByID(memberID, bandID, userID int) (*BandMember, error)` - Get a specific member
- `AddBandMember(bandID, userID int, req AddMemberRequest) (*BandMember, error)` - Add a new member
- `UpdateBandMember(memberID, bandID, userID int, req UpdateMemberRequest) (*BandMember, error)` - Update a member
- `DeleteBandMember(memberID, bandID, userID int) error` - Delete a member

### User Repository

The `UserRepository` provides a clean interface for managing users and authentication in PostgreSQL.

#### Features

- **User Management**: Create, read, update, and delete users
- **Authentication**: Secure password hashing and verification
- **Email Uniqueness**: Ensures email addresses are unique
- **Password Security**: Uses bcrypt for password hashing
- **Comprehensive CRUD**: Full user lifecycle management

#### Usage

```go
import "github.com/nahue/playlists/internal/database"

// Create repository instance
repo := database.NewUserRepository(db)

// Create a new user
req := database.CreateUserRequest{
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
    Password:  "securepassword",
}

user, err := repo.CreateUser(req)
if err != nil {
    log.Printf("Failed to create user: %v", err)
    return
}

// Authenticate user
loginReq := database.LoginRequest{
    Email:    "john@example.com",
    Password: "securepassword",
}

authUser, err := repo.AuthenticateUser(loginReq)
if err != nil {
    log.Printf("Authentication failed: %v", err)
    return
}
```

#### Available Methods

##### User Operations

- `CreateUser(req CreateUserRequest) (*UserResponse, error)` - Create a new user
- `GetUserByID(userID int) (*UserResponse, error)` - Get user by ID
- `GetUserByEmail(email string) (*User, error)` - Get user by email (includes password hash)
- `UpdateUser(userID int, req UpdateUserRequest) (*UserResponse, error)` - Update user information
- `UpdatePassword(userID int, newPassword string) error` - Update user password
- `DeleteUser(userID int) error` - Delete a user

##### Authentication

- `AuthenticateUser(req LoginRequest) (*UserResponse, error)` - Authenticate user with email/password

##### Admin Operations

- `GetAllUsers() ([]UserResponse, error)` - Get all users (for admin purposes)
- `GetUsersCount() (int, error)` - Get total number of users

## Data Models

### Band Models

#### Band

```go
type Band struct {
    ID          int       `db:"id" json:"id"`
    Name        string    `db:"name" json:"name"`
    Description string    `db:"description" json:"description"`
    UserID      int       `db:"user_id" json:"user_id"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
```

#### BandMember

```go
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
```

#### BandWithMembers

```go
type BandWithMembers struct {
    Band
    Members     []BandMember `json:"members,omitempty"`
    MemberCount int          `json:"member_count"`
}
```

### User Models

#### User

```go
type User struct {
    ID           int       `db:"id" json:"id"`
    FirstName    string    `db:"first_name" json:"first_name"`
    LastName     string    `db:"last_name" json:"last_name"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
```

#### UserResponse

```go
type UserResponse struct {
    ID        int       `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Request Types

### Band Request Types

#### CreateBandRequest

```go
type CreateBandRequest struct {
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Members     []BandMember `json:"members,omitempty"`
}
```

#### UpdateBandRequest

```go
type UpdateBandRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

#### AddMemberRequest

```go
type AddMemberRequest struct {
    Name  string `json:"name"`
    Role  string `json:"role"`
    Email string `json:"email,omitempty"`
    Phone string `json:"phone,omitempty"`
}
```

#### UpdateMemberRequest

```go
type UpdateMemberRequest struct {
    Name  string `json:"name"`
    Role  string `json:"role"`
    Email string `json:"email,omitempty"`
    Phone string `json:"phone,omitempty"`
}
```

### User Request Types

#### CreateUserRequest

```go
type CreateUserRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}
```

#### UpdateUserRequest

```go
type UpdateUserRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}
```

#### LoginRequest

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

## Security Features

### User Isolation

All repository methods ensure that users can only access their own data:

- Band operations verify `user_id` matches the authenticated user
- Member operations verify the band belongs to the authenticated user
- Unauthorized access returns `nil` instead of errors for better security

### Password Security

The user repository implements secure password handling:

- **bcrypt Hashing**: Uses bcrypt with default cost for password hashing
- **Secure Comparison**: Uses constant-time comparison for password verification
- **No Plain Text**: Passwords are never stored or returned in plain text
- **Automatic Hashing**: Password hashing is handled automatically during user creation

### Input Validation

The repositories rely on the application layer for input validation but include:

- SQL injection prevention through parameterized queries
- Proper error handling for database constraints
- Transaction rollback on errors
- Email uniqueness validation

## Testing

The repositories include comprehensive tests:

```bash
# Run all database tests
go test ./internal/database

# Run specific repository tests
go test ./internal/database -run TestBandRepository
go test ./internal/database -run TestUserRepository

# Run tests with verbose output
go test ./internal/database -v
```

### Test Setup

Tests use a separate test database (`postgres_test`) and include:

- Automatic table cleanup between tests
- Test user creation utilities
- Comprehensive CRUD operation testing
- Authentication testing
- Error condition testing

## Error Handling

The repositories use Go's error wrapping for meaningful error messages:

```go
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

Common error scenarios:
- Database connection issues
- Constraint violations
- Not found scenarios (returns `nil` instead of error)
- Authentication failures
- Duplicate email addresses

## Performance Considerations

- Uses `sqlx` for better performance and convenience
- Implements connection pooling
- Uses indexes on frequently queried columns
- Batches operations in transactions where appropriate
- Orders results for consistent pagination 

# Database Tests

This directory contains database tests for the Playlists application. The tests use a real PostgreSQL database and automatically run migrations before executing tests.

## Test Setup

### Automatic Migration

Tests automatically run database migrations before execution using the following components:

- **`test_setup.go`** - Contains the centralized test setup logic
- **`TestMain`** - Runs migrations once before all tests
- **`setupTestDB`** - Ensures database is migrated and returns a clean connection
- **`cleanupTables`** - Cleans up test data between tests

### Database Configuration

Tests use the following database configuration:
- **Host**: localhost
- **Port**: 5455
- **User**: postgres
- **Password**: postgres
- **Database**: postgres
- **SSL Mode**: disable

### Migration Process

1. **TestMain** runs once before all tests
2. Connects to the test database
3. Runs all migrations using Goose
4. Provides a clean database for all tests

### Test Isolation

Each test gets a clean database state:
- Tables are cleaned between tests
- No test data persists between test runs
- Foreign key constraints are respected during cleanup

## Running Tests

### Prerequisites

1. **PostgreSQL** running on port 5455
2. **Goose** installed and available in PATH
3. **Test database** accessible with the configured credentials

### Running All Database Tests

```bash
# From the project root
go test ./internal/database/...

# From the database directory
go test ./...
```

### Running Specific Tests

```bash
# Run only band repository tests
go test -run TestBandRepository

# Run only user repository tests
go test -run TestUserRepository

# Run migration verification test
go test -run TestMigrationsApplied
```

### Test Files

- **`test_setup.go`** - Centralized test setup and migration logic
- **`test_migrations.go`** - Verifies that migrations are applied correctly
- **`band_repository_test.go`** - Tests for band and band member operations
- **`user_repository_test.go`** - Tests for user operations and authentication

## Test Helper Functions

### `setupTestDB(t *testing.T) *sqlx.DB`

Returns a database connection with:
- Migrations applied
- Clean tables
- Proper error handling

### `cleanupTables(t *testing.T, db *sqlx.DB)`

Cleans all test data in the correct order:
1. `band_members` (due to foreign key constraints)
2. `bands`
3. `playlist_entries`
4. `users`

### `createTestUser(t *testing.T, db *sqlx.DB, email string) int`

Creates a test user and returns the user ID:
- Uses standard test data
- Handles errors properly
- Returns the created user's ID

## Migration Verification

The `TestMigrationsApplied` test verifies that:
- All expected tables exist
- Tables have the correct number of columns
- Database schema is properly set up

## Troubleshooting

### Migration Failures

If migrations fail, check:
1. PostgreSQL is running on port 5455
2. Database credentials are correct
3. Goose is installed and in PATH
4. Migration files are accessible

### Connection Issues

If database connection fails:
1. Verify PostgreSQL is running
2. Check port 5455 is accessible
3. Confirm database credentials
4. Ensure SSL mode is disabled

### Test Failures

If tests fail:
1. Check that migrations ran successfully
2. Verify table schemas match expectations
3. Ensure foreign key constraints are properly set up
4. Check that cleanup is working correctly

## Best Practices

1. **Always use `setupTestDB`** for database connections
2. **Clean up test data** between tests (handled automatically)
3. **Use descriptive test names** that indicate what is being tested
4. **Test both success and failure cases**
5. **Verify database state** after operations
6. **Use transactions** for complex test scenarios when needed 