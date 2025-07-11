package app

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.Port == "" {
		t.Error("Expected Port to have a default value")
	}

	if config.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", config.Port)
	}
}

func TestNew(t *testing.T) {
	router := chi.NewRouter()
	app := New(router)

	if app == nil {
		t.Error("Expected New() to return a non-nil Application")
	}

	if app.config == nil {
		t.Error("Expected Application to have a config")
	}

	if app.router == nil {
		t.Error("Expected Application to have a router")
	}

	if app.server != nil {
		t.Error("Expected server to be nil before initialization")
	}
}

func TestGetEnv(t *testing.T) {
	// Test with default value
	result := getEnv("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}
