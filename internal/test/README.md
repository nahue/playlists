# Test Package

This package contains comprehensive tests for the Playlists application database layer. All tests use a real PostgreSQL database and automatically run migrations before execution.

## Test Structure

### Test Files

- **`test_setup.go`** - Centralized test setup and migration logic
- **`test_migrations.go`** - Verifies that migrations are applied correctly
- **`band_repository_test.go`** - Tests for band and band member operations
- **`user_repository_test.go`** - Tests for user operations and authentication
- **`test.go`** - Database connection testing utilities

### Test Setup

#### Automatic Migration

Tests automatically run database migrations before execution using:

- **`TestMain`** - Runs migrations once before all tests
- **`setupTestDB`** - Ensures database is migrated and returns a clean connection
- **`cleanupTables`** - Cleans up test data between tests

#### Database Configuration

Tests use the following database configuration:
- **Host**: localhost
- **Port**: 5455
- **User**: postgres
- **Password**: postgres
- **Database**: postgres
- **SSL Mode**: disable

#### Migration Process

1. **TestMain** runs once before all tests
2. Connects to the test database
3. Runs all migrations using Goose
4. Provides a clean database for all tests

#### Test Isolation

Each test gets a clean database state:
- Tables are cleaned between tests
- No test data persists between test runs
- Foreign key constraints are respected during cleanup

## Running Tests

### Prerequisites

1. **PostgreSQL** running on port 5455
2. **Goose** installed and available in PATH
3. **Test database** accessible with the configured credentials

### Running All Tests

```bash
# From the project root
go test ./internal/test

# From the test directory
cd internal/test && go test
```

### Running Specific Tests

```bash
# Run only band repository tests
go test ./internal/test -run TestBandRepository

# Run only user repository tests
go test ./internal/test -run TestUserRepository

# Run migration verification test
go test ./internal/test -run TestMigrationsApplied

# Run tests with verbose output
go test ./internal/test -v
```

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

## Test Coverage

### Band Repository Tests

- **CreateBand** - Tests band creation with and without members
- **GetBandsByUserID** - Tests retrieving bands for a specific user
- **GetBandByID** - Tests retrieving a specific band with authorization
- **UpdateBand** - Tests band updates with authorization
- **DeleteBand** - Tests band deletion with authorization
- **AddBandMember** - Tests adding members to bands
- **UpdateBandMember** - Tests updating band member information
- **DeleteBandMember** - Tests removing members from bands
- **GetBandMemberByID** - Tests retrieving specific members

### User Repository Tests

- **CreateUser** - Tests user creation with validation
- **CreateUser_DuplicateEmail** - Tests email uniqueness validation
- **GetUserByID** - Tests user retrieval by ID
- **GetUserByEmail** - Tests user retrieval by email
- **AuthenticateUser** - Tests user authentication with password verification
- **UpdateUser** - Tests user information updates
- **UpdateUser_DuplicateEmail** - Tests email uniqueness during updates
- **UpdatePassword** - Tests password updates with verification
- **DeleteUser** - Tests user deletion
- **GetAllUsers** - Tests retrieving all users (admin function)
- **GetUsersCount** - Tests user count functionality

### Migration Tests

- **TestMigrationsApplied** - Verifies that all expected tables exist with correct schemas

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

## Integration with CI/CD

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions step
- name: Run Database Tests
  run: |
    # Start PostgreSQL container
    docker run -d --name test-db \
      -e POSTGRES_PASSWORD=postgres \
      -p 5455:5432 \
      postgres:15
    
    # Wait for database to be ready
    sleep 10
    
    # Run tests
    go test ./internal/test -v
```

## Performance Considerations

- Tests use a dedicated test database to avoid conflicts
- Database cleanup happens between tests for isolation
- Connection pooling is used for efficient database operations
- Tests are designed to be fast and reliable 