# User Management & RBAC Design Document

## Overview
This document outlines the design for user registration, authentication, and role-based access control (RBAC) for FleetPass.

---

## 1. Database Schema

### Current User Model
```go
type User struct {
    ID        string    `json:"id" gorm:"type:uuid;primary_key"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Proposed Enhanced User Model
```go
type User struct {
    ID                   string         `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Email                string         `json:"email" gorm:"uniqueIndex;not null"`
    Password             string         `json:"-" gorm:"not null"` // bcrypt hashed
    FirstName            string         `json:"first_name" gorm:"type:varchar(100)"`
    LastName             string         `json:"last_name" gorm:"type:varchar(100)"`
    Phone                string         `json:"phone" gorm:"type:varchar(20)"`

    // Email verification
    EmailVerified        bool           `json:"email_verified" gorm:"default:false"`
    VerificationToken    string         `json:"-" gorm:"type:varchar(255)"`
    VerificationExpiry   *time.Time     `json:"-"`

    // Password reset
    ResetToken           string         `json:"-" gorm:"type:varchar(255)"`
    ResetTokenExpiry     *time.Time     `json:"-"`

    // Status
    IsActive             bool           `json:"is_active" gorm:"default:true"`
    LastLoginAt          *time.Time     `json:"last_login_at"`

    // Relationships
    Roles                []Role         `json:"roles" gorm:"many2many:user_roles;"`
    OrganizationID       *string        `json:"organization_id" gorm:"type:uuid;index"`
    Organization         *Organization  `json:"organization,omitempty"`

    CreatedAt            time.Time      `json:"created_at"`
    UpdatedAt            time.Time      `json:"updated_at"`
}
```

**Key Changes:**
- Added `FirstName`, `LastName`, `Phone` for profile
- Added email verification fields
- Added password reset fields
- Added `IsActive` for account status
- Added `LastLoginAt` for tracking
- Added `OrganizationID` to link users to organizations
- Added many-to-many relationship with `Roles`

---

### New: Role Model
```go
type Role struct {
    ID          string       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name        string       `json:"name" gorm:"uniqueIndex;not null"` // e.g., "admin", "manager", "staff", "customer"
    DisplayName string       `json:"display_name" gorm:"not null"` // e.g., "Administrator"
    Description string       `json:"description" gorm:"type:text"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}
```

**Predefined Roles:**
1. **Super Admin** - Full system access
2. **Admin** - Organization-level admin
3. **Manager** - Location manager
4. **Staff** - Day-to-day operations
5. **Customer** - End users who rent vehicles

---

### New: Permission Model
```go
type Permission struct {
    ID          string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name        string    `json:"name" gorm:"uniqueIndex;not null"` // e.g., "vehicles.create"
    Resource    string    `json:"resource" gorm:"not null;index"` // e.g., "vehicles"
    Action      string    `json:"action" gorm:"not null;index"` // e.g., "create", "read", "update", "delete"
    Description string    `json:"description" gorm:"type:text"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Permission Format:** `resource.action`

**Example Permissions:**
- `vehicles.create` - Create vehicles
- `vehicles.read` - View vehicles
- `vehicles.update` - Update vehicles
- `vehicles.delete` - Delete vehicles
- `rentals.create` - Create rentals
- `rentals.approve` - Approve rentals
- `users.manage` - Manage users
- `organizations.manage` - Manage organizations
- `reports.view` - View reports

---

### New: Join Tables

#### user_roles
```sql
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);
```

#### role_permissions
```sql
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
```

---

## 2. Data Relationships

```
Organization (1) ──────── (many) Users

Users (many) ──────── (many) Roles

Roles (many) ──────── (many) Permissions
```

### Example Data Structure:
```
Organization: "Acme Motors"
  ├─ User: john@acme.com
  │   └─ Roles: [Admin, Manager]
  │       └─ Permissions: [vehicles.*, rentals.*, users.manage]
  │
  └─ User: jane@acme.com
      └─ Roles: [Staff]
          └─ Permissions: [vehicles.read, rentals.create, rentals.read]
```

---

## 3. Role-Permission Matrix

| Role         | Permissions                                                                 |
|--------------|-----------------------------------------------------------------------------|
| Super Admin  | ALL (full system access)                                                    |
| Admin        | All within organization: vehicles.*, rentals.*, users.*, locations.*, reports.* |
| Manager      | vehicles.*, rentals.*, reports.view, users.read                             |
| Staff        | vehicles.read, rentals.create, rentals.read, rentals.update                 |
| Customer     | rentals.create (own), rentals.read (own), vehicles.read                     |

---

## 4. API Endpoints

### Authentication Endpoints

#### POST /api/register
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "user_id": "uuid"
}
```

**Flow:**
1. Validate input (email format, password strength)
2. Check if email already exists
3. Hash password with bcrypt
4. Generate verification token
5. Create user with `email_verified: false`
6. Send verification email
7. Return success message

---

#### POST /api/verify-email
**Request:**
```json
{
  "token": "verification-token-here"
}
```

**Response:**
```json
{
  "message": "Email verified successfully. You can now log in.",
  "token": "jwt-token"
}
```

**Flow:**
1. Find user by verification token
2. Check if token expired
3. Set `email_verified: true`
4. Clear verification token
5. Generate JWT token
6. Return success with JWT

---

#### POST /api/login (Enhanced)
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "token": "jwt-token",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "roles": ["admin", "manager"],
    "permissions": ["vehicles.create", "vehicles.read", ...],
    "organization_id": "org-uuid"
  }
}
```

**Flow:**
1. Find user by email
2. Check if email verified
3. Check if account active
4. Verify password with bcrypt
5. Update `last_login_at`
6. Load user roles and permissions
7. Generate JWT with roles/permissions in claims
8. Return token and user data

---

#### POST /api/forgot-password
**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "If an account exists with this email, you will receive a password reset link."
}
```

**Flow:**
1. Find user by email
2. Generate reset token
3. Set token expiry (1 hour)
4. Send reset email
5. Return generic success message (security)

---

#### POST /api/reset-password
**Request:**
```json
{
  "token": "reset-token",
  "new_password": "NewSecurePass123!"
}
```

**Response:**
```json
{
  "message": "Password reset successful. You can now log in with your new password."
}
```

**Flow:**
1. Find user by reset token
2. Check if token expired
3. Validate new password
4. Hash new password
5. Clear reset token
6. Save user
7. Return success

---

### User Management Endpoints (Admin Only)

#### GET /api/users
**Query Params:** `?page=1&limit=20&role=admin&organization_id=uuid`

**Response:**
```json
{
  "users": [...],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

---

#### GET /api/users/:id
**Response:** Full user object with roles and permissions

---

#### PUT /api/users/:id
**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "is_active": true
}
```

---

#### POST /api/users/:id/roles
**Request:**
```json
{
  "role_ids": ["role-uuid-1", "role-uuid-2"]
}
```

**Flow:** Assign roles to user

---

#### DELETE /api/users/:id/roles/:role_id
**Flow:** Remove role from user

---

#### GET /api/users/me
**Response:** Current authenticated user with roles/permissions

---

### Role Management Endpoints (Super Admin Only)

#### GET /api/roles
#### GET /api/roles/:id
#### POST /api/roles
#### PUT /api/roles/:id
#### DELETE /api/roles/:id

---

### Permission Management Endpoints (Super Admin Only)

#### GET /api/permissions
#### GET /api/permissions/:id

---

## 5. JWT Token Structure

### Enhanced JWT Claims
```go
type Claims struct {
    UserID         string   `json:"user_id"`
    Email          string   `json:"email"`
    Roles          []string `json:"roles"`          // ["admin", "manager"]
    Permissions    []string `json:"permissions"`    // ["vehicles.create", ...]
    OrganizationID string   `json:"organization_id"`
    jwt.RegisteredClaims
}
```

**Token Example:**
```json
{
  "user_id": "uuid",
  "email": "john@example.com",
  "roles": ["admin"],
  "permissions": ["vehicles.create", "vehicles.read", "vehicles.update", "vehicles.delete"],
  "organization_id": "org-uuid",
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

## 6. Middleware Architecture

### Current: jwtauth.Authenticator
Validates JWT token exists and is valid.

### New: Permission Middleware
```go
func RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract claims from context
            // Check if user has required permission
            // If yes, call next
            // If no, return 403 Forbidden
        })
    }
}
```

### New: Role Middleware
```go
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract claims from context
            // Check if user has any of the required roles
            // If yes, call next
            // If no, return 403 Forbidden
        })
    }
}
```

### Usage Examples:
```go
// Require specific permission
r.With(RequirePermission("vehicles.create")).Post("/api/vehicles", CreateVehicle)

// Require specific role
r.With(RequireRole("admin", "manager")).Get("/api/reports", GetReports)

// Require organization context
r.With(RequireOrganization()).Get("/api/vehicles", GetVehicles)
```

---

## 7. Password Security

### Hashing
- **Algorithm:** bcrypt
- **Cost:** 12 (balance between security and performance)

### Password Requirements
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character
- Not in common password list

### Implementation:
```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## 8. Email Verification Flow

```
1. User registers
   ↓
2. System generates verification token (UUID)
   ↓
3. Email sent with link: /verify-email?token=xyz
   ↓
4. User clicks link
   ↓
5. Frontend calls: POST /api/verify-email
   ↓
6. Backend verifies token, updates user
   ↓
7. User redirected to login or auto-logged in
```

**Token Expiry:** 24 hours

---

## 9. Organization Context

### Multi-tenancy Considerations

**Approach:** Soft multi-tenancy
- Each user belongs to one organization
- Data is filtered by organization_id
- Super admins can access all organizations

### Data Isolation:
```go
// Example: Get vehicles for user's organization
func GetVehicles(userOrgID string) ([]Vehicle, error) {
    var vehicles []Vehicle
    db.Where("organization_id = ?", userOrgID).Find(&vehicles)
    return vehicles, nil
}
```

### Middleware:
```go
func RequireOrganization() func(http.Handler) http.Handler {
    // Ensures user has an organization context
    // Injects organization_id into request context
}
```

---

## 10. Frontend Changes

### New Pages
1. **Register Page** - `/register`
2. **Email Verification Page** - `/verify-email`
3. **Forgot Password Page** - `/forgot-password`
4. **Reset Password Page** - `/reset-password`
5. **User Management Page** - `/admin/users` (Admin only)
6. **User Profile Page** - `/profile`

### Updated Components
1. **AuthContext** - Add user roles/permissions
2. **ProtectedRoute** - Check permissions/roles
3. **Navbar** - Show user role, profile dropdown

### New Components
1. **RoleGuard** - Component wrapper for role-based rendering
2. **PermissionGuard** - Component wrapper for permission-based rendering

**Example Usage:**
```jsx
<RoleGuard roles={['admin', 'manager']}>
  <AdminPanel />
</RoleGuard>

<PermissionGuard permission="vehicles.create">
  <button>Create Vehicle</button>
</PermissionGuard>
```

---

## 11. Seeding Initial Data

### Seed Script Will Create:

**Permissions:**
```go
permissions := []Permission{
    {Name: "vehicles.create", Resource: "vehicles", Action: "create"},
    {Name: "vehicles.read", Resource: "vehicles", Action: "read"},
    {Name: "vehicles.update", Resource: "vehicles", Action: "update"},
    {Name: "vehicles.delete", Resource: "vehicles", Action: "delete"},
    {Name: "rentals.create", Resource: "rentals", Action: "create"},
    {Name: "rentals.read", Resource: "rentals", Action: "read"},
    {Name: "rentals.update", Resource: "rentals", Action: "update"},
    {Name: "rentals.delete", Resource: "rentals", Action: "delete"},
    {Name: "users.manage", Resource: "users", Action: "manage"},
    {Name: "organizations.manage", Resource: "organizations", Action: "manage"},
    // ... more permissions
}
```

**Roles:**
```go
roles := []Role{
    {
        Name: "super_admin",
        DisplayName: "Super Administrator",
        Description: "Full system access",
        Permissions: allPermissions,
    },
    {
        Name: "admin",
        DisplayName: "Administrator",
        Description: "Organization administrator",
        Permissions: [...],
    },
    // ... more roles
}
```

**Default Super Admin User:**
```go
User{
    Email: "admin@fleetpass.com",
    Password: bcrypt("changeme123"),
    FirstName: "System",
    LastName: "Administrator",
    EmailVerified: true,
    IsActive: true,
    Roles: [superAdminRole],
}
```

---

## 12. Migration Strategy

### Step 1: Update Database Schema
```sql
-- Add new columns to users table
ALTER TABLE users ADD COLUMN first_name VARCHAR(100);
ALTER TABLE users ADD COLUMN last_name VARCHAR(100);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN verification_token VARCHAR(255);
ALTER TABLE users ADD COLUMN verification_expiry TIMESTAMP;
ALTER TABLE users ADD COLUMN reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN reset_token_expiry TIMESTAMP;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP;
ALTER TABLE users ADD COLUMN organization_id UUID REFERENCES organizations(id);

-- Create roles table
CREATE TABLE roles (...);

-- Create permissions table
CREATE TABLE permissions (...);

-- Create join tables
CREATE TABLE user_roles (...);
CREATE TABLE role_permissions (...);
```

### Step 2: Run Seed Script
Creates permissions, roles, and default admin user.

### Step 3: Update Existing Users
Set `email_verified = true` for existing users.

---

## 13. Testing Strategy

### Backend Tests
- [ ] User registration validation
- [ ] Email verification flow
- [ ] Password reset flow
- [ ] Login with roles/permissions
- [ ] Permission checking middleware
- [ ] Role checking middleware
- [ ] Organization isolation

### Frontend Tests
- [ ] Registration form validation
- [ ] Login with role-based redirects
- [ ] RoleGuard component
- [ ] PermissionGuard component
- [ ] Protected routes with permissions

---

## 14. Security Considerations

### Email Verification
- Tokens are single-use
- Tokens expire after 24 hours
- Tokens are UUIDs (not sequential)

### Password Reset
- Tokens expire after 1 hour
- Tokens are single-use
- Generic messages (don't reveal if email exists)

### JWT
- Short expiry (1 hour recommended)
- Refresh token pattern (future)
- HTTPS only
- HttpOnly cookies (optional, more secure than localStorage)

### RBAC
- Permissions checked on backend (never trust frontend)
- Organization isolation enforced at database level
- Audit logging for sensitive actions

---

## 15. Future Enhancements

- [ ] Refresh tokens
- [ ] OAuth integration (Google, Microsoft)
- [ ] Two-factor authentication (2FA)
- [ ] Session management (view/revoke sessions)
- [ ] Audit logs
- [ ] User impersonation (admin feature)
- [ ] Custom roles (user-defined)
- [ ] Permission groups

---

## Questions to Consider

1. **Email Service:** Should we use SendGrid, AWS SES, or mock for now?
2. **Organization Assignment:** How do new users get assigned to organizations?
   - Self-select during registration?
   - Invitation-based?
   - Admin assigns after registration?
3. **Default Role:** What role should new users get by default? (Suggest: Customer)
4. **Multi-org Users:** Should users be able to belong to multiple organizations? (Suggest: No for now)
5. **Email Templates:** Should we create HTML email templates or plain text?

---

## Summary of Changes

### Database
- Update `users` table (10 new columns)
- Create `roles` table
- Create `permissions` table
- Create `user_roles` join table
- Create `role_permissions` join table

### Backend (Go)
- Update User model
- Create Role model
- Create Permission model
- Create registration handlers
- Create email verification handlers
- Create password reset handlers
- Create user management handlers
- Create role management handlers
- Create permission middleware
- Create role middleware
- Update login handler with roles/permissions
- Create seed script

### Frontend (React)
- Create registration page
- Create email verification page
- Create forgot password page
- Create reset password page
- Create user management page
- Update AuthContext with roles/permissions
- Create RoleGuard component
- Create PermissionGuard component
- Update ProtectedRoute component
- Create user profile page

### Testing
- Add 20+ new backend tests
- Add 15+ new frontend tests

---

**Estimated Effort:** 2-3 weeks for full implementation and testing.

Let me know if you have questions or want to adjust any part of this design!
