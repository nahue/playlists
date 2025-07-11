# Routes Package

This package provides route configuration and middleware setup for the Playlists application.

## Overview

The routes package is responsible for:
- Setting up the Chi router with middleware
- Configuring CORS for frontend integration
- Defining all API endpoints
- Setting up authentication middleware
- Serving static files

## Structure

### SetupRoutes Function
```go
func SetupRoutes(app *app.Application) *chi.Mux
```

This function creates and configures the router with:
- CORS middleware
- Standard HTTP middleware (logging, recovery, etc.)
- Authentication routes
- Protected API routes
- Static file serving

## Middleware Configuration

### CORS Middleware
Configured to allow requests from:
- `http://localhost:3000`
- `http://localhost:3001`
- `http://localhost:4321`

Supports all common HTTP methods and includes necessary headers for authentication.

### Standard Middleware
- **Logger** - Request logging
- **Recoverer** - Panic recovery
- **CleanPath** - URL path cleaning
- **GetHead** - GET/HEAD method handling

## Route Structure

### Public Routes (`/auth`)
- `POST /auth/register` - User registration
- `POST /auth/login` - User authentication
- `POST /auth/logout` - User logout

### Protected Routes (`/api`)
All routes under `/api` require authentication via JWT token.

#### User Management
- `GET /api/profile` - Get user profile

#### Task Management (`/api/todos`)
- `GET /api/todos` - Get all tasks
- `POST /api/todos` - Create new task
- `PUT /api/todos/{id}` - Update task
- `DELETE /api/todos/{id}` - Delete task

#### Playlist Management (`/api/playlist`)
- `GET /api/playlist` - Get all playlist entries
- `POST /api/playlist` - Add new playlist entry
- `GET /api/playlist/artists` - Artist autocomplete
- `GET /api/playlist/users` - User autocomplete
- `GET /api/playlist/{id}` - Get specific entry
- `PUT /api/playlist/{id}` - Update entry
- `DELETE /api/playlist/{id}` - Delete entry

#### Band Management (`/api/bands`)
- `GET /api/bands` - Get all bands
- `POST /api/bands` - Create new band
- `GET /api/bands/{id}` - Get specific band
- `PUT /api/bands/{id}` - Update band
- `DELETE /api/bands/{id}` - Delete band

#### Band Members (`/api/bands/{bandId}/members`)
- `GET /api/bands/{bandId}/members` - Get band members
- `POST /api/bands/{bandId}/members` - Add member
- `PUT /api/bands/{bandId}/members/{memberId}` - Update member
- `DELETE /api/bands/{bandId}/members/{memberId}` - Remove member

### Static Files
- `/*` - Serves frontend files from `./frontend/dist`

## Authentication

Protected routes use the `handlers.AuthMiddleware` which:
- Validates JWT tokens from Authorization header
- Extracts user information from token
- Adds user context to request
- Returns 401 for invalid/missing tokens

## Usage

```go
func main() {
    // Create application instance
    application := app.NewApplication()

    // Setup routes with application context
    router := routes.SetupRoutes(application)

    // Create server with router
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Start server
    server.ListenAndServe()
}
```

## Benefits

1. **Centralized Routing** - All routes defined in one place
2. **Middleware Management** - Consistent middleware across all routes
3. **Authentication Integration** - Seamless JWT authentication
4. **CORS Support** - Frontend integration ready
5. **Static File Serving** - Built-in frontend support
6. **Application Context** - Access to database and configuration

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/go-chi/chi/middleware` - Chi middleware
- `github.com/go-chi/cors` - CORS middleware
- `github.com/nahue/playlists/internal/app` - Application struct
- `github.com/nahue/playlists/internal/handlers` - Request handlers 