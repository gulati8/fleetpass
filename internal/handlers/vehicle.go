package handlers

import (
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetVehicles(w http.ResponseWriter, r *http.Request) {
	var vehicles []models.Vehicle

	if err := database.DB.Order("created_at DESC").Find(&vehicles).Error; err != nil {
		http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func GetVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, "id = ?", id).Error; err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.LocationID == "" || req.VIN == "" || req.Make == "" || req.Model == "" || req.Year == 0 {
		http.Error(w, "Location ID, VIN, make, model, and year are required", http.StatusBadRequest)
		return
	}

	// Get location to extract organization_id
	var location models.Location
	if err := database.DB.First(&location, "id = ?", req.LocationID).Error; err != nil {
		http.Error(w, "Location not found", http.StatusBadRequest)
		return
	}

	vehicle := models.Vehicle{
		OrganizationID:       location.OrganizationID,
		LocationID:           req.LocationID,
		VIN:                  req.VIN,
		Make:                 req.Make,
		Model:                req.Model,
		Year:                 req.Year,
		Trim:                 req.Trim,
		ColorExterior:        req.ColorExterior,
		ColorInterior:        req.ColorInterior,
		Condition:            req.Condition,
		Mileage:              req.Mileage,
		LicensePlate:         req.LicensePlate,
		Status:               models.VehicleStatusAvailable,
		IsEligibleForService: true,
		BodyStyle:            req.BodyStyle,
		Transmission:         req.Transmission,
		Drivetrain:           req.Drivetrain,
		FuelType:             req.FuelType,
		Engine:               req.Engine,
		MPGCity:              req.MPGCity,
		MPGHighway:           req.MPGHighway,
		Seats:                req.Seats,
		Doors:                req.Doors,
		StockNumber:          req.StockNumber,
		Description:          req.Description,
		DailyRate:            req.DailyRate,
		WeeklyRate:           req.WeeklyRate,
		MonthlyRate:          req.MonthlyRate,
		Features:             models.StringArray(req.Features),
		Images:               models.StringArray(req.Images),
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		http.Error(w, "Failed to create vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, "id = ?", id).Error; err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	// Update fields
	vehicle.LocationID = req.LocationID
	vehicle.Make = req.Make
	vehicle.Model = req.Model
	vehicle.Year = req.Year
	vehicle.Trim = req.Trim
	vehicle.ColorExterior = req.ColorExterior
	vehicle.ColorInterior = req.ColorInterior
	vehicle.Condition = req.Condition
	vehicle.Mileage = req.Mileage
	vehicle.LicensePlate = req.LicensePlate
	vehicle.Status = req.Status
	vehicle.IsEligibleForService = req.IsEligibleForService
	vehicle.BodyStyle = req.BodyStyle
	vehicle.Transmission = req.Transmission
	vehicle.Drivetrain = req.Drivetrain
	vehicle.FuelType = req.FuelType
	vehicle.Engine = req.Engine
	vehicle.MPGCity = req.MPGCity
	vehicle.MPGHighway = req.MPGHighway
	vehicle.Seats = req.Seats
	vehicle.Doors = req.Doors
	vehicle.StockNumber = req.StockNumber
	vehicle.Description = req.Description
	vehicle.DailyRate = req.DailyRate
	vehicle.WeeklyRate = req.WeeklyRate
	vehicle.MonthlyRate = req.MonthlyRate
	vehicle.Features = models.StringArray(req.Features)
	vehicle.Images = models.StringArray(req.Images)

	if err := database.DB.Save(&vehicle).Error; err != nil {
		http.Error(w, "Failed to update vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := database.DB.Delete(&models.Vehicle{}, "id = ?", id)
	if result.Error != nil {
		http.Error(w, "Failed to delete vehicle", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
