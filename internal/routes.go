package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nahue/playlists/internal/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r chi.Router) {
	// Auth routes (public)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)
		r.Post("/logout", handlers.Logout)
	})

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		// User profile
		r.Get("/profile", handlers.GetProfile)

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

		// Band routes
		r.Route("/bands", func(r chi.Router) {
			r.Get("/", handlers.GetBands)
			r.Post("/", handlers.CreateBand)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetBand)
				r.Put("/", handlers.UpdateBand)
				r.Delete("/", handlers.DeleteBand)
			})
		})
	})

	// Serve built frontend
	r.Handle("/*", http.FileServer(http.Dir("./frontend/dist")))
}
