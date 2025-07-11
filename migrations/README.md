# Database Migrations

This directory contains database migrations for the Playlists application using Goose.

## Migration Files

1. **20250711135758_users_table.sql** - Creates the users table
2. **20250711135810_bands_table.sql** - Creates the bands table
3. **20250711135816_band_members_table.sql** - Creates the band members table
4. **20250711135822_playlist_table.sql** - Creates the playlist entries table

## Database Schema

### Users Table
- `id` - Primary key
- `first_name` - User's first name
- `last_name` - User's last name
- `email` - Unique email address
- `password_hash` - Hashed password using bcrypt
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Bands Table
- `id` - Primary key
- `name` - Band name
- `description` - Band description
- `user_id` - Foreign key to users table (band owner)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Band Members Table
- `id` - Primary key
- `band_id` - Foreign key to bands table
- `name` - Member name
- `role` - Member role/instrument
- `email` - Member email (optional)
- `phone` - Member phone (optional)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Playlist Entries Table
- `id` - Primary key
- `artist` - Artist name
- `song` - Song title
- `user_name` - User who added the song
- `user_id` - Foreign key to users table
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Usage

### Running Migrations

```bash
# Run all pending migrations
goose up

# Run migrations with custom database URL
goose -dir ./migrations postgres "host=localhost port=5454 user=postgres password=postgres dbname=postgres sslmode=disable" up

# Rollback last migration
goose down

# Check migration status
goose status
```

### Creating New Migrations

```bash
# Create a new migration
goose create migration_name sql
```

## Integration with Application

The database connection is integrated into the Application struct with repositories:

```go
type Application struct {
    Logger      *log.Logger
    Config      *Config
    DB          *sqlx.DB
    BandHandler *handlers.BandHandler
    AuthHandler *handlers.AuthHandler
}
```

Migrations are run during application setup, and the database connection is available throughout the application lifecycle via repositories.

## Features

- **Automatic timestamps** - `created_at` and `updated_at` fields are automatically managed
- **Foreign key constraints** - Proper relationships between tables
- **Indexes** - Performance optimizations for common queries
- **Cascade deletes** - When a user is deleted, their bands and playlist entries are also deleted
- **Unique constraints** - Email addresses must be unique
- **Connection pooling** - Optimized database connections
- **SQLx integration** - Enhanced database operations with struct scanning
- **Repository pattern** - Clean data access layer with dependency injection 