package models

import "time"

type Location struct {
	ID             string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrganizationID string    `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`
	AddressLine1   string    `json:"address_line1" gorm:"type:varchar(255)"`
	AddressLine2   string    `json:"address_line2" gorm:"type:varchar(255)"`
	City           string    `json:"city" gorm:"type:varchar(100)"`
	State          string    `json:"state" gorm:"type:varchar(50)"`
	ZipCode        string    `json:"zip_code" gorm:"type:varchar(20)"`
	Country        string    `json:"country" gorm:"type:varchar(100)"`
	Phone          string    `json:"phone" gorm:"type:varchar(50)"`
	Email          string    `json:"email" gorm:"type:varchar(255)"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Location) TableName() string {
	return "locations"
}

type CreateLocationRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	City           string `json:"city"`
	State          string `json:"state"`
	ZipCode        string `json:"zip_code"`
	Country        string `json:"country"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

type UpdateLocationRequest struct {
	Name         string `json:"name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsActive     bool   `json:"is_active"`
}
