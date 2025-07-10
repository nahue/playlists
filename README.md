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
- **CORS** - Cross-Origin Resource Sharing support

### Frontend
- **Alpine.js** - Lightweight JavaScript framework
- **Tailwind CSS** - Utility-first CSS framework
- **HTML5** - Modern semantic markup

## 📋 Prerequisites

- Go 1.24.4 or higher
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

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Open in browser**
   ```
   http://localhost:8080
   ```

## 📚 API Documentation

### Task Endpoints

#### GET /todos
Get all tasks.
```bash
curl http://localhost:8080/todos
```

#### POST /todos
Create a new task.
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"text": "New task"}'
```

#### PUT /todos/{id}
Update an existing task.
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"text": "Updated task"}'
```

#### DELETE /todos/{id}
Delete a task.
```bash
curl -X DELETE http://localhost:8080/todos/1
```

### Playlist Endpoints

#### GET /playlist
Get all songs in the playlist.
```bash
curl http://localhost:8080/playlist
```

#### POST /playlist
Add a new song to the playlist.
```bash
curl -X POST http://localhost:8080/playlist \
  -H "Content-Type: application/json" \
  -d '{
    "artist": "Artist",
    "song": "Song",
    "user_name": "User"
  }'
```

#### PUT /playlist/{id}
Update an existing song.
```bash
curl -X PUT http://localhost:8080/playlist/1 \
  -H "Content-Type: application/json" \
  -d '{
    "artist": "Updated Artist",
    "song": "Updated Song",
    "user_name": "User"
  }'
```

#### DELETE /playlist/{id}
Delete a song from the playlist.
```bash
curl -X DELETE http://localhost:8080/playlist/1
```

#### GET /playlist/artists?q={query}
Get artist suggestions for autocomplete.
```bash
curl "http://localhost:8080/playlist/artists?q=artist"
```

#### GET /playlist/users?q={query}
Get username suggestions for autocomplete.
```bash
curl "http://localhost:8080/playlist/users?q=user"
```

## 🏗️ Project Structure

```
playlists/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── go.sum                  # Dependency checksums
├── README.md              # This file
├── frontend/              # Frontend files
│   └── index.html         # Main interface
└── internal/              # Internal application code
    ├── routes.go          # Route configuration
    └── handlers/          # Endpoint handlers
        ├── todo_handlers.go      # Task logic
        └── playlist_handlers.go  # Playlist logic
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
2. Fill out the form with:
   - **Artist**: Artist name (with autocomplete)
   - **Song**: Song title
   - **Your Name**: Your username (with autocomplete)
3. Click "Add Song"
4. To edit, click the edit icon in the corresponding row
5. To delete, click the delete icon

## 🔧 Configuration

### Environment Variables
The application runs on port 8080 by default. To change the port, modify the line in `main.go`:

```go
log.Fatal(http.ListenAndServe(":8080", r))
```

### CORS
The application is configured to allow requests from:
- `http://localhost:3000`
- `http://localhost:3001`

To modify the allowed origins, edit the CORS configuration in `main.go`.

## 🧪 Development

### Run in development mode
```bash
go run main.go
```

### Build for production
```bash
go build -o playlists main.go
./playlists
```

### Run tests (when added)
```bash
go test ./...
```

## 📝 Technical Notes

- **Storage**: Data is stored in memory (lost when server restarts)
- **Validation**: Required fields are validated to not be empty
- **Autocomplete**: Suggestions are based on existing data in the list
- **Responsive**: Interface adapts to different screen sizes

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is under the MIT License. See the `LICENSE` file for more details.

## 👨‍💻 Author

**Nahuel** - [GitHub](https://github.com/nahue)

---

⭐ If you like this project, give it a star on GitHub! 