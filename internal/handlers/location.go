package handlers

import (
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetLocations(w http.ResponseWriter, r *http.Request) {
	var locations []models.Location

	if err := database.DB.Order("created_at DESC").Find(&locations).Error; err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func GetLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var location models.Location
	if err := database.DB.First(&location, "id = ?", id).Error; err != nil {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.OrganizationID == "" || req.Name == "" {
		http.Error(w, "Organization ID and name are required", http.StatusBadRequest)
		return
	}

	location := models.Location{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		City:           req.City,
		State:          req.State,
		ZipCode:        req.ZipCode,
		Country:        req.Country,
		Phone:          req.Phone,
		Email:          req.Email,
		IsActive:       true,
	}

	if err := database.DB.Create(&location).Error; err != nil {
		http.Error(w, "Failed to create location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var location models.Location
	if err := database.DB.First(&location, "id = ?", id).Error; err != nil {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	// Update fields
	location.Name = req.Name
	location.AddressLine1 = req.AddressLine1
	location.AddressLine2 = req.AddressLine2
	location.City = req.City
	location.State = req.State
	location.ZipCode = req.ZipCode
	location.Country = req.Country
	location.Phone = req.Phone
	location.Email = req.Email
	location.IsActive = req.IsActive

	if err := database.DB.Save(&location).Error; err != nil {
		http.Error(w, "Failed to update location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := database.DB.Delete(&models.Location{}, "id = ?", id)
	if result.Error != nil {
		http.Error(w, "Failed to delete location", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
