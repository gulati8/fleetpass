package models

import "time"

type Permission struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"` // e.g., "vehicles.create"
	Resource    string    `json:"resource" gorm:"not null;index"`   // e.g., "vehicles"
	Action      string    `json:"action" gorm:"not null;index"`     // e.g., "create", "read", "update", "delete"
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Predefined permissions
const (
	// Vehicle permissions
	PermissionVehiclesCreate = "vehicles.create"
	PermissionVehiclesRead   = "vehicles.read"
	PermissionVehiclesUpdate = "vehicles.update"
	PermissionVehiclesDelete = "vehicles.delete"

	// Rental permissions
	PermissionRentalsCreate  = "rentals.create"
	PermissionRentalsRead    = "rentals.read"
	PermissionRentalsUpdate  = "rentals.update"
	PermissionRentalsDelete  = "rentals.delete"
	PermissionRentalsApprove = "rentals.approve"

	// User permissions
	PermissionUsersManage = "users.manage"
	PermissionUsersRead   = "users.read"

	// Organization permissions
	PermissionOrganizationsManage = "organizations.manage"
	PermissionOrganizationsRead   = "organizations.read"

	// Location permissions
	PermissionLocationsCreate = "locations.create"
	PermissionLocationsRead   = "locations.read"
	PermissionLocationsUpdate = "locations.update"
	PermissionLocationsDelete = "locations.delete"

	// Report permissions
	PermissionReportsView = "reports.view"

	// System permissions
	PermissionSystemManage = "system.manage"
)
