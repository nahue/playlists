package routes

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nahue/playlists/internal/app"
	"github.com/nahue/playlists/internal/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware for frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:4321"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)

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
			// Band members routes
			r.Route("/{bandId}/members", func(r chi.Router) {
				r.Get("/", handlers.GetBandMembers)
				r.Post("/", handlers.AddBandMember)
				r.Route("/{memberId}", func(r chi.Router) {
					r.Put("/", handlers.UpdateBandMember)
					r.Delete("/", handlers.DeleteBandMember)
				})
			})
		})
	})

	// Serve built frontend
	r.Handle("/*", http.FileServer(http.Dir("./frontend/dist")))

	return r
}
