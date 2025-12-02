package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fleetpass/internal/database"
	"fleetpass/internal/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type BulkUploadRequest struct {
	OrganizationID string `json:"organization_id"`
	LocationID     string `json:"location_id"`
}

type BulkUploadResult struct {
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	Total     int      `json:"total"`
	Errors    []string `json:"errors,omitempty"`
	VehicleIDs []string `json:"vehicle_ids,omitempty"`
}

func BulkUploadVehicles(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get organization and location IDs
	organizationID := r.FormValue("organization_id")
	locationID := r.FormValue("location_id")

	if organizationID == "" || locationID == "" {
		http.Error(w, "organization_id and location_id are required", http.StatusBadRequest)
		return
	}

	// Verify location exists and belongs to organization
	var location models.Location
	if err := database.DB.First(&location, "id = ? AND organization_id = ?", locationID, organizationID).Error; err != nil {
		http.Error(w, "Location not found or does not belong to organization", http.StatusBadRequest)
		return
	}

	// Get the CSV file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "Failed to read CSV headers", http.StatusBadRequest)
		return
	}

	// Map headers to indices
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(header))] = i
	}

	// Validate required headers
	requiredHeaders := []string{"vin", "make", "model", "year"}
	for _, required := range requiredHeaders {
		if _, exists := headerMap[required]; !exists {
			http.Error(w, fmt.Sprintf("Missing required header: %s", required), http.StatusBadRequest)
			return
		}
	}

	result := BulkUploadResult{
		Errors:     []string{},
		VehicleIDs: []string{},
	}

	rowNum := 1 // Start at 1 (header is row 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to parse CSV row", rowNum))
			rowNum++
			continue
		}

		rowNum++
		result.Total++

		// Parse vehicle from CSV row
		vehicle, err := parseVehicleFromCSV(record, headerMap, organizationID, locationID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}

		// Create vehicle in database
		if err := database.DB.Create(&vehicle).Error; err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d (VIN: %s): %s", rowNum, vehicle.VIN, err.Error()))
			continue
		}

		result.Success++
		result.VehicleIDs = append(result.VehicleIDs, vehicle.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	if result.Failed > 0 {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(result)
}

func parseVehicleFromCSV(record []string, headerMap map[string]int, organizationID, locationID string) (*models.Vehicle, error) {
	getValue := func(key string) string {
		if idx, exists := headerMap[key]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	getIntValue := func(key string) int {
		val := getValue(key)
		if val == "" {
			return 0
		}
		intVal, _ := strconv.Atoi(val)
		return intVal
	}

	getFloatValue := func(key string) float64 {
		val := getValue(key)
		if val == "" {
			return 0
		}
		floatVal, _ := strconv.ParseFloat(val, 64)
		return floatVal
	}

	// Required fields
	vin := getValue("vin")
	make := getValue("make")
	model := getValue("model")
	yearStr := getValue("year")

	if vin == "" || make == "" || model == "" || yearStr == "" {
		return nil, fmt.Errorf("missing required fields (vin, make, model, year)")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		return nil, fmt.Errorf("invalid year: %s", yearStr)
	}

	// Parse condition
	condition := models.VehicleCondition(getValue("condition"))
	if condition != "" && condition != models.VehicleConditionNew &&
		condition != models.VehicleConditionUsed &&
		condition != models.VehicleConditionCertifiedPreOwned {
		condition = models.VehicleConditionUsed // Default to used
	}

	// Parse features (pipe-separated)
	var features []string
	featuresStr := getValue("features")
	if featuresStr != "" {
		features = strings.Split(featuresStr, "|")
		for i := range features {
			features[i] = strings.TrimSpace(features[i])
		}
	}

	vehicle := &models.Vehicle{
		OrganizationID:       organizationID,
		LocationID:           locationID,
		VIN:                  vin,
		Make:                 make,
		Model:                model,
		Year:                 year,
		Trim:                 getValue("trim"),
		ColorExterior:        getValue("color_exterior"),
		ColorInterior:        getValue("color_interior"),
		Condition:            condition,
		Mileage:              getIntValue("mileage"),
		LicensePlate:         getValue("license_plate"),
		Status:               models.VehicleStatusAvailable,
		IsEligibleForService: true,
		BodyStyle:            getValue("body_style"),
		Transmission:         getValue("transmission"),
		Drivetrain:           getValue("drivetrain"),
		FuelType:             getValue("fuel_type"),
		Engine:               getValue("engine"),
		MPGCity:              getIntValue("mpg_city"),
		MPGHighway:           getIntValue("mpg_highway"),
		Seats:                getIntValue("seats"),
		Doors:                getIntValue("doors"),
		StockNumber:          getValue("stock_number"),
		Description:          getValue("description"),
		DailyRate:            getFloatValue("daily_rate"),
		WeeklyRate:           getFloatValue("weekly_rate"),
		MonthlyRate:          getFloatValue("monthly_rate"),
		Features:             models.StringArray(features),
		Images:               models.StringArray([]string{}),
	}

	return vehicle, nil
}
