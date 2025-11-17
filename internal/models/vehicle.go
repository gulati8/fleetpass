package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringArray is a custom type for handling PostgreSQL arrays
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(a)
}

type VehicleCondition string
type VehicleStatus string

const (
	VehicleConditionNew             VehicleCondition = "new"
	VehicleConditionUsed            VehicleCondition = "used"
	VehicleConditionCertifiedPreOwned VehicleCondition = "certified_pre_owned"

	VehicleStatusAvailable   VehicleStatus = "available"
	VehicleStatusRented      VehicleStatus = "rented"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusInactive    VehicleStatus = "inactive"
)

type Vehicle struct {
	ID                   string           `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrganizationID       string           `json:"organization_id" gorm:"type:uuid;not null;index"`
	LocationID           string           `json:"location_id" gorm:"type:uuid;not null;index"`
	VIN                  string           `json:"vin" gorm:"type:varchar(17);uniqueIndex;not null"`
	Make                 string           `json:"make" gorm:"type:varchar(100);not null;index:idx_make_model"`
	Model                string           `json:"model" gorm:"type:varchar(100);not null;index:idx_make_model"`
	Year                 int              `json:"year" gorm:"not null;index"`
	Trim                 string           `json:"trim" gorm:"type:varchar(100)"`
	ColorExterior        string           `json:"color_exterior" gorm:"type:varchar(100)"`
	ColorInterior        string           `json:"color_interior" gorm:"type:varchar(100)"`
	Condition            VehicleCondition `json:"condition" gorm:"type:varchar(50)"`
	Mileage              int              `json:"mileage" gorm:"default:0"`
	LicensePlate         string           `json:"license_plate" gorm:"type:varchar(20)"`
	Status               VehicleStatus    `json:"status" gorm:"type:varchar(50);default:'available';index"`
	IsEligibleForService bool             `json:"is_eligible_for_service" gorm:"default:true"`

	// Warranty
	HasWarranty            bool       `json:"has_warranty" gorm:"default:false"`
	WarrantyExpirationDate *time.Time `json:"warranty_expiration_date,omitempty"`
	WarrantyType           string     `json:"warranty_type" gorm:"type:varchar(100)"`
	WarrantyDetails        string     `json:"warranty_details" gorm:"type:text"`

	// Pricing
	DailyRate   float64 `json:"daily_rate" gorm:"type:decimal(10,2);default:0"`
	WeeklyRate  float64 `json:"weekly_rate" gorm:"type:decimal(10,2);default:0"`
	MonthlyRate float64 `json:"monthly_rate" gorm:"type:decimal(10,2);default:0"`

	// Additional details
	BodyStyle    string `json:"body_style" gorm:"type:varchar(50)"`
	Transmission string `json:"transmission" gorm:"type:varchar(100)"`
	Drivetrain   string `json:"drivetrain" gorm:"type:varchar(50)"`
	FuelType     string `json:"fuel_type" gorm:"type:varchar(50)"`
	Engine       string `json:"engine" gorm:"type:varchar(100)"`
	MPGCity      int    `json:"mpg_city"`
	MPGHighway   int    `json:"mpg_highway"`
	Seats        int    `json:"seats"`
	Doors        int    `json:"doors"`
	StockNumber  string `json:"stock_number" gorm:"type:varchar(50)"`
	Description  string `json:"description" gorm:"type:text"`
	Features     StringArray `json:"features" gorm:"type:jsonb"`
	Images       StringArray `json:"images" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

type CreateVehicleRequest struct {
	LocationID           string           `json:"location_id"`
	VIN                  string           `json:"vin"`
	Make                 string           `json:"make"`
	Model                string           `json:"model"`
	Year                 int              `json:"year"`
	Trim                 string           `json:"trim"`
	ColorExterior        string           `json:"color_exterior"`
	ColorInterior        string           `json:"color_interior"`
	Condition            VehicleCondition `json:"condition"`
	Mileage              int              `json:"mileage"`
	LicensePlate         string           `json:"license_plate"`
	BodyStyle            string           `json:"body_style"`
	Transmission         string           `json:"transmission"`
	Drivetrain           string           `json:"drivetrain"`
	FuelType             string           `json:"fuel_type"`
	Engine               string           `json:"engine"`
	MPGCity              int              `json:"mpg_city"`
	MPGHighway           int              `json:"mpg_highway"`
	Seats                int              `json:"seats"`
	Doors                int              `json:"doors"`
	StockNumber          string           `json:"stock_number"`
	Description          string           `json:"description"`
	DailyRate            float64          `json:"daily_rate"`
	WeeklyRate           float64          `json:"weekly_rate"`
	MonthlyRate          float64          `json:"monthly_rate"`
	Features             []string         `json:"features"`
	Images               []string         `json:"images"`
}

type UpdateVehicleRequest struct {
	LocationID           string           `json:"location_id"`
	Make                 string           `json:"make"`
	Model                string           `json:"model"`
	Year                 int              `json:"year"`
	Trim                 string           `json:"trim"`
	ColorExterior        string           `json:"color_exterior"`
	ColorInterior        string           `json:"color_interior"`
	Condition            VehicleCondition `json:"condition"`
	Mileage              int              `json:"mileage"`
	LicensePlate         string           `json:"license_plate"`
	Status               VehicleStatus    `json:"status"`
	IsEligibleForService bool             `json:"is_eligible_for_service"`
	BodyStyle            string           `json:"body_style"`
	Transmission         string           `json:"transmission"`
	Drivetrain           string           `json:"drivetrain"`
	FuelType             string           `json:"fuel_type"`
	Engine               string           `json:"engine"`
	MPGCity              int              `json:"mpg_city"`
	MPGHighway           int              `json:"mpg_highway"`
	Seats                int              `json:"seats"`
	Doors                int              `json:"doors"`
	StockNumber          string           `json:"stock_number"`
	Description          string           `json:"description"`
	DailyRate            float64          `json:"daily_rate"`
	WeeklyRate           float64          `json:"weekly_rate"`
	MonthlyRate          float64          `json:"monthly_rate"`
	Features             []string         `json:"features"`
	Images               []string         `json:"images"`
}
