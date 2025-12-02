package handlers

import (
	"bytes"
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"fleetpass/internal/testutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBulkUploadVehicles_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")

	// Create CSV content
	csvContent := `vin,make,model,year,trim,color_exterior,condition,mileage
1HGBH41JXMN109186,Honda,Accord,2022,EX-L,Silver,used,15000
1FTFW1ET8EFA12345,Ford,F-150,2023,XLT,Blue,new,500`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	writer.WriteField("organization_id", org.ID)
	writer.WriteField("location_id", loc.ID)

	// Add file
	part, _ := writer.CreateFormFile("file", "vehicles.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles/bulk-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Execute handler
	BulkUploadVehicles(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result BulkUploadResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}

	if result.Success != 2 {
		t.Errorf("Expected success 2, got %d", result.Success)
	}

	if result.Failed != 0 {
		t.Errorf("Expected failed 0, got %d. Errors: %v", result.Failed, result.Errors)
	}

	// Verify vehicles were created
	var count int64
	db.Model(&models.Vehicle{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 vehicles in database, got %d", count)
	}
}

func TestBulkUploadVehicles_PartialSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")

	// Create CSV with one valid and one invalid row
	csvContent := `vin,make,model,year,trim,color_exterior,condition,mileage
1HGBH41JXMN109186,Honda,Accord,2022,EX-L,Silver,used,15000
INVALID,,Model,9999,,,invalid,abc`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("organization_id", org.ID)
	writer.WriteField("location_id", loc.ID)

	part, _ := writer.CreateFormFile("file", "vehicles.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles/bulk-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Execute handler
	BulkUploadVehicles(w, req)

	// Assertions
	if w.Code != http.StatusPartialContent {
		t.Errorf("Expected status %d, got %d", http.StatusPartialContent, w.Code)
	}

	var result BulkUploadResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Success != 1 {
		t.Errorf("Expected success 1, got %d", result.Success)
	}

	if result.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Failed)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error messages for failed rows")
	}
}

func TestBulkUploadVehicles_MissingFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create test data
	org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")
	loc := testutil.CreateTestLocation(t, db, org.ID, "Test Location", "San Francisco")

	// Create CSV without required headers
	csvContent := `vin,make,model
1HGBH41JXMN109186,Honda,Accord`

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("organization_id", org.ID)
	writer.WriteField("location_id", loc.ID)

	part, _ := writer.CreateFormFile("file", "vehicles.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles/bulk-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Execute handler
	BulkUploadVehicles(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBulkUploadVehicles_InvalidLocation(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	database.DB = db

	// Create CSV content
	csvContent := `vin,make,model,year
1HGBH41JXMN109186,Honda,Accord,2022`

	// Create multipart form with invalid location ID
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("organization_id", "invalid-org-id")
	writer.WriteField("location_id", "invalid-loc-id")

	part, _ := writer.CreateFormFile("file", "vehicles.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles/bulk-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Execute handler
	BulkUploadVehicles(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
