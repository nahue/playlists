package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nahue/playlists/internal/app"
	"github.com/nahue/playlists/internal/routes"
)

func main() {
	// Create new application instance with router
	application := app.NewApplication()
	// Create router and setup middleware
	r := routes.SetupRoutes(application)

	// Ensure graceful shutdown
	defer func() {
		if err := application.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

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
