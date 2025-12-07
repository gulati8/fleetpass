package models

import (
	"time"
)

type User struct {
	ID        string `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"` // bcrypt hashed, never send in JSON
	FirstName string `json:"first_name" gorm:"type:varchar(100)"`
	LastName  string `json:"last_name" gorm:"type:varchar(100)"`
	Phone     string `json:"phone" gorm:"type:varchar(20)"`

	// Email verification
	EmailVerified      bool       `json:"email_verified" gorm:"default:false"`
	VerificationToken  string     `json:"-" gorm:"type:varchar(255);index"`
	VerificationExpiry *time.Time `json:"-"`

	// Password reset
	ResetToken       string     `json:"-" gorm:"type:varchar(255);index"`
	ResetTokenExpiry *time.Time `json:"-"`

	// Status
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	LastLoginAt *time.Time `json:"last_login_at"`

	// Organization relationship
	OrganizationID *string        `json:"organization_id" gorm:"type:uuid;index"`
	Organization   *Organization  `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`

	// Roles relationship (many-to-many)
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// Helper method to get all permissions for a user
func (u *User) GetPermissions() []string {
	permissions := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permissions[perm.Name] = true
		}
	}

	result := make([]string, 0, len(permissions))
	for perm := range permissions {
		result = append(result, perm)
	}
	return result
}

// Helper method to check if user has a specific permission
func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission {
				return true
			}
		}
	}
	return false
}

// Helper method to check if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// Request/Response types

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Phone           string `json:"phone"`
	OrganizationID  string `json:"organization_id"` // Optional: for invitation-based registration
	InvitationToken string `json:"invitation_token"` // For invitation-based registration
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Phone          string   `json:"phone"`
	EmailVerified  bool     `json:"email_verified"`
	IsActive       bool     `json:"is_active"`
	Roles          []string `json:"roles"`
	Permissions    []string `json:"permissions"`
	OrganizationID *string  `json:"organization_id"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"is_active"`
}

type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required"`
}
