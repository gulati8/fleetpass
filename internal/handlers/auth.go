package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Hardcoded users for now (will replace with database later)
var users = map[string]struct {
	ID       string
	Password string
	Role     string
}{
	"admin@fleetpass.com": {
		ID:       "1",
		Password: "admin123",
		Role:     "admin",
	},
	"user@fleetpass.com": {
		ID:       "2",
		Password: "user123",
		Role:     "user",
	},
}

func Login(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate credentials
		user, exists := users[req.Email]
		if !exists || user.Password != req.Password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		claims := map[string]interface{}{
			"user_id": user.ID,
			"email":   req.Email,
			"role":    user.Role,
		}
		jwtauth.SetExpiryIn(claims, 24*time.Hour)
		_, tokenString, err := tokenAuth.Encode(claims)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Return token and user info
		response := LoginResponse{
			Token: tokenString,
			User: UserProfile{
				ID:    user.ID,
				Email: req.Email,
				Role:  user.Role,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	profile := UserProfile{
		ID:    claims["user_id"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
