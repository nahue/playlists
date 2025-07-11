# Database Setup Guide

This guide explains how to set up and use the PostgreSQL database for the Playlists application.

## Prerequisites

- PostgreSQL installed and running
- Goose migration tool installed (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

## Database Configuration

The application uses environment variables for database configuration. You can set these in the `config.env` file:

```env
DB_HOST=localhost
DB_PORT=5454
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
```

## Database Integration

The database is integrated into the Application struct:

```go
type Application struct {
    Logger *log.Logger
    Config *Config
    DB     *sqlx.DB
}
```

The database connection is:
- Established during application initialization
- Tested for connectivity
- Available throughout the application lifecycle
- Properly configured with connection pooling

## Running Migrations

### Option 1: Using the migration script
```bash
./scripts/migrate.sh
```

### Option 2: Using Goose directly
```bash
# Load environment variables
source config.env

# Run migrations
goose -dir ./migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE" up
```

### Option 3: Using environment variables
```bash
export DB_HOST=localhost
export DB_PORT=5454
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=postgres
export DB_SSLMODE=disable

goose -dir ./migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE" up
```

## Migration Commands

```bash
# Check migration status
goose -dir ./migrations postgres "$DB_URL" status

# Rollback last migration
goose -dir ./migrations postgres "$DB_URL" down

# Rollback to specific version
goose -dir ./migrations postgres "$DB_URL" down-to 20250711135758

# Reset all migrations
goose -dir ./migrations postgres "$DB_URL" reset
```

## Database Schema

The application creates the following tables:

1. **users** - User accounts and authentication
2. **bands** - Band information and ownership
3. **band_members** - Band member details and roles
4. **todos** - Todo items with completion status
5. **playlist_entries** - Music playlist entries

## Connection Details

- **Driver**: PostgreSQL
- **Connection Pool**: Max 25 open connections, 5 idle connections
- **SSL Mode**: Disabled (for local development)
- **Timezone**: Uses database default
- **Integration**: Direct access via `application.DB`

## Using Database in Handlers

Handlers can access the database through the Application struct:

```go
func GetTodos(w http.ResponseWriter, r *http.Request) {
    // Get application from context or pass directly
    app := r.Context().Value("app").(*app.Application)
    
    // Use database connection
    todos := []Todo{}
    err := app.DB.Select(&todos, "SELECT * FROM todos WHERE user_id = $1", userID)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    // Return response
    json.NewEncoder(w).Encode(todos)
}
```

## Database Operations

### Using SQLx
The application uses SQLx for enhanced database operations:

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

## Troubleshooting

### Connection Issues
1. Ensure PostgreSQL is running
2. Check database credentials in `config.env`
3. Verify database exists and is accessible
4. Check firewall settings

### Migration Issues
1. Ensure Goose is installed: `go install github.com/pressly/goose/v3/cmd/goose@latest`
2. Check database permissions
3. Verify migration files are in the correct directory

### Performance
- Indexes are created for common query patterns
- Connection pooling is configured for optimal performance
- Foreign key constraints ensure data integrity

## Development vs Production

### Development
- Use local PostgreSQL instance
- SSL mode disabled
- Default credentials (change in production)

### Production
- Use managed PostgreSQL service
- Enable SSL mode
- Use strong, unique passwords
- Set appropriate connection limits
- Enable connection pooling

## Testing Database Connection

Use the provided test script:
```bash
./scripts/test-db.sh
```

This script will:
- Test connection with psql (if available)
- Test connection with Goose
- Verify database accessibility 