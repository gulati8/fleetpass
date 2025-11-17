package handlers

import (
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	var organizations []models.Organization

	if err := database.DB.Order("created_at DESC").Find(&organizations).Error; err != nil {
		http.Error(w, "Failed to fetch organizations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organizations)
}

func GetOrganization(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var org models.Organization
	if err := database.DB.First(&org, "id = ?", id).Error; err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Name == "" || req.Slug == "" {
		http.Error(w, "Name and slug are required", http.StatusBadRequest)
		return
	}

	org := models.Organization{
		Name:     req.Name,
		Slug:     req.Slug,
		IsActive: true,
	}

	if err := database.DB.Create(&org).Error; err != nil {
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(org)
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var org models.Organization
	if err := database.DB.First(&org, "id = ?", id).Error; err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	// Update fields
	org.Name = req.Name
	org.Slug = req.Slug
	org.IsActive = req.IsActive

	if err := database.DB.Save(&org).Error; err != nil {
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := database.DB.Delete(&models.Organization{}, "id = ?", id)
	if result.Error != nil {
		http.Error(w, "Failed to delete organization", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
