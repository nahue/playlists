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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.GetHead)

	// Auth routes (public)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.AuthHandler.Register)
		r.Post("/login", app.AuthHandler.Login)
		r.Post("/logout", app.AuthHandler.Logout)
	})

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		// User profile
		r.Get("/profile", handlers.GetProfile)

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
			r.Get("/", app.BandHandler.GetBands)
			r.Post("/", app.BandHandler.CreateBand)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", app.BandHandler.GetBand)
				r.Put("/", app.BandHandler.UpdateBand)
				r.Delete("/", app.BandHandler.DeleteBand)
			})
			// Band members routes
			r.Route("/{bandId}/members", func(r chi.Router) {
				r.Get("/", app.BandHandler.GetBandMembers)
				r.Post("/", app.BandHandler.AddBandMember)
				r.Route("/{memberId}", func(r chi.Router) {
					r.Put("/", app.BandHandler.UpdateBandMember)
					r.Delete("/", app.BandHandler.DeleteBandMember)
				})
			})
		})
	})

	// Serve built frontend
	r.Handle("/*", http.FileServer(http.Dir("./frontend/dist")))

	return r
}
