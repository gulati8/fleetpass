package testutil

import (
	"fleetpass/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use environment variable or default test database
	dsn := "host=localhost user=fleetpass_user password=fleetpass_password dbname=fleetpass_test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v. Run tests with a test database available.", err)
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Organization{},
		&models.Location{},
		&models.Vehicle{},
		&models.User{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestDB removes all test data
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Delete in reverse order of dependencies
	db.Exec("TRUNCATE TABLE vehicles CASCADE")
	db.Exec("TRUNCATE TABLE locations CASCADE")
	db.Exec("TRUNCATE TABLE organizations CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")
}

// CreateTestOrganization creates a test organization
func CreateTestOrganization(t *testing.T, db *gorm.DB, name, slug string) *models.Organization {
	t.Helper()

	org := &models.Organization{
		Name:     name,
		Slug:     slug,
		IsActive: true,
	}

	if err := db.Create(org).Error; err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	return org
}

// CreateTestLocation creates a test location
func CreateTestLocation(t *testing.T, db *gorm.DB, orgID, name, city string) *models.Location {
	t.Helper()

	loc := &models.Location{
		OrganizationID: orgID,
		Name:           name,
		City:           city,
		State:          "CA",
		Country:        "USA",
		IsActive:       true,
	}

	if err := db.Create(loc).Error; err != nil {
		t.Fatalf("Failed to create test location: %v", err)
	}

	return loc
}

// CreateTestVehicle creates a test vehicle
func CreateTestVehicle(t *testing.T, db *gorm.DB, orgID, locationID, vin, make, model string, year int) *models.Vehicle {
	t.Helper()

	vehicle := &models.Vehicle{
		OrganizationID:       orgID,
		LocationID:           locationID,
		VIN:                  vin,
		Make:                 make,
		Model:                model,
		Year:                 year,
		Status:               models.VehicleStatusAvailable,
		IsEligibleForService: true,
	}

	if err := db.Create(vehicle).Error; err != nil {
		t.Fatalf("Failed to create test vehicle: %v", err)
	}

	return vehicle
}

// CreateTestUser creates a test user
func CreateTestUser(t *testing.T, db *gorm.DB, email, password string) *models.User {
	t.Helper()

	user := &models.User{
		Email:    email,
		Password: password, // In real tests, this should be hashed
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}
