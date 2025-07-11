# Application Package

This package provides the main application structure and lifecycle management for the Playlists application.

## Overview

The `Application` struct encapsulates application configuration, database connection, repositories, handlers, and logging setup. It provides a centralized way to manage application state and dependencies using dependency injection.

## Structure

### Application
```go
type Application struct {
    Logger      *log.Logger
    Config      *Config
    DB          *sqlx.DB
    BandHandler *handlers.BandHandler
    AuthHandler *handlers.AuthHandler
}
```

### Config
```go
type Config struct {
    Port string
    Host string
}
```

## Methods

### NewApplication()
Creates a new Application instance with:
- Environment variable loading
- Database connection setup
- Database connection testing
- Repository initialization (Band, User)
- Handler initialization with dependency injection
- Logger configuration
- Configuration initialization

### NewConfig()
Creates a new Config instance from environment variables.

## Usage

```go
func main() {
    // Create new application instance with all dependencies
    application := app.NewApplication()

    // Create router and setup middleware
    r := routes.SetupRoutes(application)

    // Ensure graceful shutdown
    defer func() {
        if err := application.Shutdown(); err != nil {
            log.Printf("Error during shutdown: %v", err)
        }
    }()

    // Create HTTP server
    server := &http.Server{
        Addr:         fmt.Sprintf("%s:%s", application.Config.Host, application.Config.Port),
        Handler:      r,
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    application.Logger.Printf("Starting server on port %s", application.Config.Port)

    err := server.ListenAndServe()
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

## Dependency Injection

The Application struct uses dependency injection to provide:

### Repositories
- **BandRepository** - Manages bands and band members in PostgreSQL
- **UserRepository** - Manages users and authentication in PostgreSQL

### Handlers
- **BandHandler** - HTTP handlers for band operations with injected repository
- **AuthHandler** - HTTP handlers for authentication with injected repository and logger

### Benefits of Dependency Injection
- **Testability** - Easy to mock dependencies for testing
- **Separation of Concerns** - Clear boundaries between layers
- **Maintainability** - Easy to swap implementations
- **Configuration** - Dependencies are configured once and reused

## Configuration

The application uses environment variables for configuration:

- `SERVER_PORT` - Server port (default: "8080")
- `SERVER_HOST` - Server host (default: "localhost")

## Database Integration

The Application struct includes a database connection that is:
- Established during initialization
- Tested for connectivity
- Available throughout the application lifecycle
- Properly configured with connection pooling
- Used by repositories for data access

## Repository Pattern

The application implements the repository pattern:

```go
// Repositories are instantiated in NewApplication()
bandRepo := database.NewBandRepository(db)
userRepo := database.NewUserRepository(db)

// Handlers receive repositories via dependency injection
bandHandler := handlers.NewBandHandler(bandRepo, logger)
authHandler := handlers.NewAuthHandler(userRepo, logger)
```

## Handler Pattern

Handlers use dependency injection to receive their dependencies:

```go
// BandHandler receives repository and logger
type BandHandler struct {
    bandRepo *database.BandRepository
    logger   *log.Logger
}

// AuthHandler receives repository and logger
type AuthHandler struct {
    userRepo *database.UserRepository
    logger   *log.Logger
}
```

## Logging

The Application struct provides a configured logger that:
- Outputs to stdout
- Includes timestamps and file information
- Is injected into handlers for consistent logging
- Can be used throughout the application

## Benefits

1. **Centralized Configuration** - All app settings in one place
2. **Database Integration** - Direct access to database connection and repositories
3. **Dependency Injection** - Clean separation of concerns and testability
4. **Structured Logging** - Consistent logging across the application
5. **Environment Management** - Automatic environment variable loading
6. **Repository Pattern** - Clean data access layer
7. **Handler Pattern** - Organized HTTP request handling
8. **Testability** - Application can be tested with mock dependencies

## Testing

Run tests with:
```bash
go test ./internal/app/...
```

## Dependencies

- `github.com/jmoiron/sqlx` - Enhanced database operations
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/nahue/playlists/internal/database` - Database repositories
- `github.com/nahue/playlists/internal/handlers` - HTTP handlers 