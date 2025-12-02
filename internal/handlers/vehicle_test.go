package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"fleetpass/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetVehicles(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Set the global database for handlers
	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")
	testutil.CreateTestVehicle(t, db, org.ID, loc.ID, "1HGBH41JXMN109186", "Honda", "Accord", 2022)
	testutil.CreateTestVehicle(t, db, org.ID, loc.ID, "1FTFW1ET8EFA12345", "Ford", "F-150", 2023)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	w := httptest.NewRecorder()

	// Execute handler
	GetVehicles(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var vehicles []models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&vehicles); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(vehicles) != 2 {
		t.Errorf("Expected 2 vehicles, got %d", len(vehicles))
	}
}

func TestGetVehicle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")
	vehicle := testutil.CreateTestVehicle(t, db, org.ID, loc.ID, "1HGBH41JXMN109186", "Honda", "Accord", 2022)

	// Create request with URL params
	req := httptest.NewRequest(http.MethodGet, "/api/vehicles/"+vehicle.ID, nil)
	w := httptest.NewRecorder()

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", vehicle.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute handler
	GetVehicle(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.VIN != vehicle.VIN {
		t.Errorf("Expected VIN %s, got %s", vehicle.VIN, result.VIN)
	}
}

func TestCreateVehicle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")

	// Create request
	reqBody := models.CreateVehicleRequest{
		LocationID:    loc.ID,
		VIN:           "1HGBH41JXMN109186",
		Make:          "Honda",
		Model:         "Accord",
		Year:          2022,
		Trim:          "EX-L",
		ColorExterior: "Silver",
		Condition:     models.VehicleConditionUsed,
		Mileage:       15000,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	CreateVehicle(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var vehicle models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&vehicle); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if vehicle.VIN != reqBody.VIN {
		t.Errorf("Expected VIN %s, got %s", reqBody.VIN, vehicle.VIN)
	}

	if vehicle.Make != reqBody.Make {
		t.Errorf("Expected make %s, got %s", reqBody.Make, vehicle.Make)
	}
}

func TestCreateVehicle_ValidationError(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create request with missing required fields
	reqBody := models.CreateVehicleRequest{
		Make:  "Honda",
		Model: "Accord",
		// Missing VIN, LocationID, Year
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	CreateVehicle(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateVehicle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")
	vehicle := testutil.CreateTestVehicle(t, db, org.ID, loc.ID, "1HGBH41JXMN109186", "Honda", "Accord", 2022)

	// Create update request
	reqBody := models.UpdateVehicleRequest{
		LocationID: loc.ID,
		Make:       "Honda",
		Model:      "Accord",
		Year:       2022,
		Mileage:    20000, // Updated mileage
		Status:     models.VehicleStatusMaintenance,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/vehicles/"+vehicle.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", vehicle.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute handler
	UpdateVehicle(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updated models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if updated.Mileage != reqBody.Mileage {
		t.Errorf("Expected mileage %d, got %d", reqBody.Mileage, updated.Mileage)
	}

	if updated.Status != reqBody.Status {
		t.Errorf("Expected status %s, got %s", reqBody.Status, updated.Status)
	}
}

func TestDeleteVehicle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")
	vehicle := testutil.CreateTestVehicle(t, db, org.ID, loc.ID, "1HGBH41JXMN109186", "Honda", "Accord", 2022)

	// Create delete request
	req := httptest.NewRequest(http.MethodDelete, "/api/vehicles/"+vehicle.ID, nil)
	w := httptest.NewRecorder()

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", vehicle.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute handler
	DeleteVehicle(w, req)

	// Assertions
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify vehicle is deleted
	var count int64
	db.Model(&models.Vehicle{}).Where("id = ?", vehicle.ID).Count(&count)
	if count != 0 {
		t.Error("Vehicle was not deleted from database")
	}
}
