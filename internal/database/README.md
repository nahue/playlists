# Database Package

This package provides database functionality for the Playlists application using PostgreSQL and SQLx.

## Features

- **PostgreSQL Integration**: Full PostgreSQL support with connection pooling
- **SQLx Enhanced**: Enhanced database operations with struct scanning
- **Repository Pattern**: Clean data access layer with dependency injection
- **Transaction Support**: Full transaction support for complex operations
- **Connection Management**: Proper connection lifecycle management
- **Error Handling**: Comprehensive error handling and wrapping

## Usage

```go
import "github.com/nahue/playlists/internal/database"

// Create database configuration
config := &database.Config{
    Host:     "localhost",
    Port:     "5454",
    User:     "postgres",
    Password: "postgres",
    DBName:   "postgres",
    SSLMode:  "disable",
}

// Open database connection
db, err := database.Open(config)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Create repositories
bandRepo := database.NewBandRepository(db)
userRepo := database.NewUserRepository(db)
```

## Repositories

### Band Repository

The `BandRepository` provides a clean interface for managing bands and band members in PostgreSQL.

#### Features

- **Band Management**: Create, read, update, and delete bands
- **Member Management**: Add, update, and delete band members
- **User Isolation**: Bands are scoped to their owners
- **Transaction Support**: Complex operations use transactions
- **Comprehensive CRUD**: Full band lifecycle management

#### Usage

```go
import "github.com/nahue/playlists/internal/database"

// Create repository instance
repo := database.NewBandRepository(db)

// Create a new band
req := database.CreateBandRequest{
    Name:        "My Band",
    Description: "A great band",
    Members: []database.BandMember{
        {Name: "John", Role: "Guitarist"},
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

Tests are located in the `internal/test` package and include comprehensive tests for all repository methods:

```bash
# Run all database tests
go test ./internal/test

# Run specific repository tests
go test ./internal/test -run TestBandRepository
go test ./internal/test -run TestUserRepository

# Run tests with verbose output
go test ./internal/test -v
```

### Test Setup

Tests use a separate test database and include:

- Automatic table cleanup between tests
- Test user creation utilities
- Comprehensive CRUD operation testing
- Authentication testing
- Error condition testing

## Error Handling

The repositories use Go's error wrapping for meaningful error messages:

```go
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

## Database Operations

### Using SQLx with Repositories
The application uses SQLx for enhanced database operations within repositories:

```go
// Query with struct scanning
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE active = $1", true)

// Query single row
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)

// Execute with result
result, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)

// Transaction
tx, err := db.Beginx()
if err != nil {
    return err
}
defer tx.Rollback()

// Use transaction
err = tx.Exec("INSERT INTO users (name) VALUES ($1)", name)
if err != nil {
    return err
}

err = tx.Commit()
```

## Security Features

### Password Security
- **bcrypt Hashing**: Passwords are hashed using bcrypt with default cost
- **Secure Comparison**: Constant-time password verification
- **No Plain Text**: Passwords are never stored or returned in plain text

### User Isolation
- **Scoped Operations**: All band operations are scoped to the authenticated user
- **Authorization**: Users can only access their own data
- **Database Constraints**: Foreign key constraints ensure data integrity

## Performance Features

### Connection Pooling
- **Optimized Connections**: Efficient connection pool management
- **Connection Limits**: Configurable connection pool limits
- **Connection Timeouts**: Proper timeout handling

### Query Optimization
- **Prepared Statements**: Uses prepared statements for better performance
- **Batch Operations**: Batches operations in transactions where appropriate
- **Orders results for consistent pagination** 