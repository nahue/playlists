# Application Package

This package provides the main application structure and lifecycle management for the Playlists application.

## Overview

The `Application` struct encapsulates application configuration, database connection, and logging setup. It provides a centralized way to manage application state and dependencies.

## Structure

### Application
```go
type Application struct {
    Logger *log.Logger
    Config *Config
    DB     *sqlx.DB
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
- Logger configuration
- Configuration initialization

### NewConfig()
Creates a new Config instance from environment variables.

## Usage

```go
func main() {
    // Create new application instance
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

## Logging

The Application struct provides a configured logger that:
- Outputs to stdout
- Includes timestamps and file information
- Can be used throughout the application for consistent logging

## Benefits

1. **Centralized Configuration** - All app settings in one place
2. **Database Integration** - Direct access to database connection
3. **Structured Logging** - Consistent logging across the application
4. **Environment Management** - Automatic environment variable loading
5. **Dependency Injection** - Easy to pass application context to handlers
6. **Testability** - Application can be tested with mock dependencies

## Testing

Run tests with:
```bash
go test ./internal/app/...
```

## Dependencies

- `github.com/jmoiron/sqlx` - Enhanced database operations
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/nahue/playlists/internal/database` - Database connection management 