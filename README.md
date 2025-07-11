# 🎵 Task & Playlist Manager

A modern web application for managing tasks and music playlists, built with Go on the backend and Alpine.js on the frontend.

## ✨ Features

### 📝 Task Management
- ✅ Create, edit and delete tasks
- ✅ Modern and responsive interface
- ✅ Real-time validation
- ✅ Deletion confirmations

### 🎵 Playlist Management
- ✅ Add songs with artist, title and username
- ✅ Smart autocomplete for artists and users
- ✅ Edit and delete songs
- ✅ Intuitive interface with smooth transitions

### 🎸 Band Management
- ✅ Create and manage bands
- ✅ Add and manage band members with roles
- ✅ Member contact information
- ✅ Band descriptions and details

### 🎨 User Interface
- ✅ Modern design with Tailwind CSS
- ✅ Tab navigation between modules
- ✅ Smooth animations and transitions
- ✅ Fully Spanish interface
- ✅ Responsive design for mobile and desktop

## 🚀 Technologies Used

### Backend
- **Go 1.24.4** - Programming language
- **Chi Router** - Lightweight and fast HTTP router
- **PostgreSQL** - Relational database
- **SQLx** - Enhanced database operations
- **JWT** - Authentication tokens
- **CORS** - Cross-Origin Resource Sharing support

### Frontend
- **Alpine.js** - Lightweight JavaScript framework
- **Tailwind CSS** - Utility-first CSS framework
- **HTML5** - Modern semantic markup

## 📋 Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- Modern web browser

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/nahue/playlists.git
   cd playlists
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure database**
   ```bash
   # Copy and edit the configuration file
   cp config.env.example config.env
   # Edit config.env with your database settings
   ```

4. **Run database migrations**
   ```bash
   ./scripts/migrate.sh
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Open in browser**
   ```
   http://localhost:8080
   ```

## 🏗️ Architecture

### Application Structure
The application uses a clean, modular architecture:

- **`main.go`** - Entry point with server configuration
- **`internal/app/`** - Application struct and lifecycle management
- **`internal/routes/`** - Route configuration and middleware setup
- **`internal/handlers/`** - HTTP request handlers
- **`internal/database/`** - Database connection and operations
- **`migrations/`** - Database schema migrations

### Key Components

#### Application Struct
```go
type Application struct {
    Logger *log.Logger
    Config *Config
    DB     *sqlx.DB
}
```

The Application struct encapsulates:
- Database connection management
- Configuration handling
- Logging setup
- Application lifecycle

#### Route Setup
Routes are configured in `internal/routes/routes.go` with:
- CORS middleware for frontend integration
- Authentication middleware for protected routes
- RESTful API endpoints
- Static file serving

## 📚 API Documentation

### Authentication Endpoints

#### POST /auth/register
Register a new user account.
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### POST /auth/login
Authenticate user and get JWT token.
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Task Endpoints

#### GET /api/todos
Get all tasks for authenticated user.
```bash
curl http://localhost:8080/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### POST /api/todos
Create a new task.
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"text": "New task"}'
```

### Playlist Endpoints

#### GET /api/playlist
Get all songs in the playlist.
```bash
curl http://localhost:8080/api/playlist \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### POST /api/playlist
Add a new song to the playlist.
```bash
curl -X POST http://localhost:8080/api/playlist \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "artist": "Artist",
    "song": "Song",
    "user_name": "User"
  }'
```

### Band Endpoints

#### GET /api/bands
Get all bands for authenticated user.
```bash
curl http://localhost:8080/api/bands \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### POST /api/bands
Create a new band.
```bash
curl -X POST http://localhost:8080/api/bands \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Band Name",
    "description": "Band description"
  }'
```

#### POST /api/bands/{bandId}/members
Add a member to a band.
```bash
curl -X POST http://localhost:8080/api/bands/1/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Member Name",
    "role": "Guitar",
    "email": "member@example.com"
  }'
```

## 🏗️ Project Structure

```
playlists/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── go.sum                  # Dependency checksums
├── config.env              # Environment configuration
├── README.md              # This file
├── DATABASE.md            # Database setup guide
├── migrations/            # Database migrations
├── scripts/               # Utility scripts
├── frontend/              # Frontend files
│   └── src/               # Source files
└── internal/              # Internal application code
    ├── app/               # Application struct and lifecycle
    ├── routes/            # Route configuration
    ├── handlers/          # HTTP handlers
    └── database/          # Database operations
```

## 🎯 Usage

### Task Management
1. Navigate to the "Task App" tab
2. Type a new task in the text field
3. Click "Add Task" or press Enter
4. To edit, click the edit icon
5. To delete, click the delete icon

### Playlist Management
1. Navigate to the "Playlist Manager" tab
2. Fill in artist, song, and username fields
3. Click "Add to Playlist"
4. Use autocomplete for existing artists and users
5. Edit or delete entries as needed

### Band Management
1. Navigate to the "Bands" section
2. Create a new band with name and description
3. Add members with roles and contact information
4. Edit band details and member information
5. Manage band membership

## 🔧 Development

### Running in Development Mode
```bash
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Database Management
```bash
# Run migrations
./scripts/migrate.sh

# Test database connection
./scripts/test-db.sh
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/app/...
```

## 🚀 Deployment

### Environment Variables
Set the following environment variables for production:

```env
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
DB_SSLMODE=require
JWT_SECRET=your-secure-jwt-secret
SERVER_PORT=8080
```

### Docker Deployment
```bash
# Build the application
docker build -t playlists .

# Run with environment variables
docker run -p 8080:8080 --env-file config.env playlists
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.