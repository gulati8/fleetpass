package main

import (
	"fmt"
	"log"
	"net/http"

	"fleetpass/internal/database"
	"fleetpass/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	// Initialize JWT auth with secret key
	// TODO: Move this to environment variable
	tokenAuth = jwtauth.New("HS256", []byte("your-secret-key-change-this-in-production"), nil)
}

func main() {
	// Initialize database connection
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized successfully")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("FleetPass API v1.0"))
		})
		r.Post("/api/login", handlers.Login(tokenAuth))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/api/profile", handlers.GetProfile)
		r.Get("/api/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is a protected endpoint"))
		})

		// Organizations
		r.Get("/api/organizations", handlers.GetOrganizations)
		r.Post("/api/organizations", handlers.CreateOrganization)
		r.Get("/api/organizations/{id}", handlers.GetOrganization)
		r.Put("/api/organizations/{id}", handlers.UpdateOrganization)
		r.Delete("/api/organizations/{id}", handlers.DeleteOrganization)

		// Locations
		r.Get("/api/locations", handlers.GetLocations)
		r.Post("/api/locations", handlers.CreateLocation)
		r.Get("/api/locations/{id}", handlers.GetLocation)
		r.Put("/api/locations/{id}", handlers.UpdateLocation)
		r.Delete("/api/locations/{id}", handlers.DeleteLocation)

		// Vehicles
		r.Get("/api/vehicles", handlers.GetVehicles)
		r.Post("/api/vehicles", handlers.CreateVehicle)
		r.Post("/api/vehicles/bulk-upload", handlers.BulkUploadVehicles)
		r.Get("/api/vehicles/{id}", handlers.GetVehicle)
		r.Put("/api/vehicles/{id}", handlers.UpdateVehicle)
		r.Delete("/api/vehicles/{id}", handlers.DeleteVehicle)
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
