package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nahue/playlists/internal/database"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	userRepo *database.UserRepository
	logger   *log.Logger
}

// NewAuthHandler creates a new AuthHandler with the given repository and logger
func NewAuthHandler(userRepo *database.UserRepository, logger *log.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req database.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.userRepo.CreateUser(req)
	if err != nil {
		h.logger.Printf("Failed to create user: %v", err)
		if strings.Contains(err.Error(), "email already exists") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(*user)
	if err != nil {
		h.logger.Printf("Failed to generate JWT: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req database.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.userRepo.AuthenticateUser(req)
	if err != nil {
		h.logger.Printf("Authentication failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(*user)
	if err != nil {
		h.logger.Printf("Failed to generate JWT: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		h.logger.Printf("Failed to get user profile: %v", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Logout (client-side token removal, but we can blacklist tokens here if needed)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// For now, just return success - client should remove the token
	// In a production system, you might want to blacklist the token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// AuthMiddleware validates JWT tokens and adds user info to request context
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			h.logger.Printf("Token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Verify user still exists in database
		user, err := h.userRepo.GetUserByID(claims.UserID)
		if err != nil {
			h.logger.Printf("Failed to verify user: %v", err)
			http.Error(w, "User verification failed", http.StatusUnauthorized)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userEmail", claims.Email)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateJWT creates a JWT token for the given user
func (h *AuthHandler) generateJWT(user database.UserResponse) (string, error) {
	// Generate a random string for additional security
	randomString, err := h.generateRandomString(32)
	if err != nil {
		return "", err
	}

	// Create claims
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "playlists-app",
			Subject:   strconv.Itoa(user.ID),
			ID:        randomString,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateRandomString creates a random string of the specified length
func (h *AuthHandler) generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
