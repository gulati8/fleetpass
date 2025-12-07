package models

import "time"

type Role struct {
	ID          string       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"` // e.g., "super_admin", "admin", "manager", "staff", "customer"
	DisplayName string       `json:"display_name" gorm:"not null"`     // e.g., "Super Administrator"
	Description string       `json:"description" gorm:"type:text"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Role) TableName() string {
	return "roles"
}

// Predefined role names
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleStaff      = "staff"
	RoleCustomer   = "customer"
)

// Request/Response types

type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required"`
	DisplayName   string   `json:"display_name" validate:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	DisplayName   string   `json:"display_name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}
