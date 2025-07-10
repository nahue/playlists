package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nahue/playlists/internal/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r chi.Router) {
	// Todo routes
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", handlers.GetTodos)
		r.Post("/", handlers.CreateTodo)
		r.Route("/{id}", func(r chi.Router) {
			r.Put("/", handlers.UpdateTodo)
			r.Delete("/", handlers.DeleteTodo)
		})
	})

	// Playlist routes
	r.Route("/playlist", func(r chi.Router) {
		r.Get("/", handlers.GetPlaylist)
		r.Post("/", handlers.AddToPlaylist)
		r.Get("/artists", handlers.GetArtists) // Artist autocomplete endpoint
		r.Get("/users", handlers.GetUserNames) // User name autocomplete endpoint
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetPlaylistEntry)
			r.Put("/", handlers.UpdatePlaylistEntry)
			r.Delete("/", handlers.DeletePlaylistEntry)
		})
	})

	// API routes only - frontend will be served by Nuxt
	r.Handle("/*", http.FileServer(http.Dir("./frontend")))
}
