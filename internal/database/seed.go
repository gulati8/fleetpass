package database

import (
	"fleetpass/internal/auth"
	"fleetpass/internal/models"
	"log"

	"gorm.io/gorm"
)

// SeedDatabase populates the database with initial data
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed permissions
	if err := seedPermissions(db); err != nil {
		return err
	}

	// Seed roles
	if err := seedRoles(db); err != nil {
		return err
	}

	// Seed super admin user
	if err := seedSuperAdmin(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedPermissions(db *gorm.DB) error {
	log.Println("Seeding permissions...")

	permissions := []models.Permission{
		// Vehicle permissions
		{Name: models.PermissionVehiclesCreate, Resource: "vehicles", Action: "create", Description: "Create new vehicles"},
		{Name: models.PermissionVehiclesRead, Resource: "vehicles", Action: "read", Description: "View vehicles"},
		{Name: models.PermissionVehiclesUpdate, Resource: "vehicles", Action: "update", Description: "Update vehicles"},
		{Name: models.PermissionVehiclesDelete, Resource: "vehicles", Action: "delete", Description: "Delete vehicles"},

		// Rental permissions
		{Name: models.PermissionRentalsCreate, Resource: "rentals", Action: "create", Description: "Create rentals"},
		{Name: models.PermissionRentalsRead, Resource: "rentals", Action: "read", Description: "View rentals"},
		{Name: models.PermissionRentalsUpdate, Resource: "rentals", Action: "update", Description: "Update rentals"},
		{Name: models.PermissionRentalsDelete, Resource: "rentals", Action: "delete", Description: "Delete rentals"},
		{Name: models.PermissionRentalsApprove, Resource: "rentals", Action: "approve", Description: "Approve rentals"},

		// User permissions
		{Name: models.PermissionUsersManage, Resource: "users", Action: "manage", Description: "Manage users"},
		{Name: models.PermissionUsersRead, Resource: "users", Action: "read", Description: "View users"},

		// Organization permissions
		{Name: models.PermissionOrganizationsManage, Resource: "organizations", Action: "manage", Description: "Manage organizations"},
		{Name: models.PermissionOrganizationsRead, Resource: "organizations", Action: "read", Description: "View organizations"},

		// Location permissions
		{Name: models.PermissionLocationsCreate, Resource: "locations", Action: "create", Description: "Create locations"},
		{Name: models.PermissionLocationsRead, Resource: "locations", Action: "read", Description: "View locations"},
		{Name: models.PermissionLocationsUpdate, Resource: "locations", Action: "update", Description: "Update locations"},
		{Name: models.PermissionLocationsDelete, Resource: "locations", Action: "delete", Description: "Delete locations"},

		// Report permissions
		{Name: models.PermissionReportsView, Resource: "reports", Action: "view", Description: "View reports"},

		// System permissions
		{Name: models.PermissionSystemManage, Resource: "system", Action: "manage", Description: "Manage system settings"},
	}

	for _, perm := range permissions {
		if err := db.Where(models.Permission{Name: perm.Name}).FirstOrCreate(&perm).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d permissions", len(permissions))
	return nil
}

func seedRoles(db *gorm.DB) error {
	log.Println("Seeding roles...")

	// Get all permissions for super admin
	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// Define roles with their permissions
	roles := []struct {
		role            models.Role
		permissionNames []string
	}{
		{
			role: models.Role{
				Name:        models.RoleSuperAdmin,
				DisplayName: "Super Administrator",
				Description: "Full system access - can manage everything",
			},
			permissionNames: []string{}, // Will get all permissions
		},
		{
			role: models.Role{
				Name:        models.RoleAdmin,
				DisplayName: "Administrator",
				Description: "Organization administrator - can manage organization, locations, vehicles, and users",
			},
			permissionNames: []string{
				models.PermissionVehiclesCreate, models.PermissionVehiclesRead, models.PermissionVehiclesUpdate, models.PermissionVehiclesDelete,
				models.PermissionRentalsCreate, models.PermissionRentalsRead, models.PermissionRentalsUpdate, models.PermissionRentalsDelete, models.PermissionRentalsApprove,
				models.PermissionUsersManage, models.PermissionUsersRead,
				models.PermissionOrganizationsRead,
				models.PermissionLocationsCreate, models.PermissionLocationsRead, models.PermissionLocationsUpdate, models.PermissionLocationsDelete,
				models.PermissionReportsView,
			},
		},
		{
			role: models.Role{
				Name:        models.RoleManager,
				DisplayName: "Manager",
				Description: "Location manager - can manage vehicles and rentals",
			},
			permissionNames: []string{
				models.PermissionVehiclesCreate, models.PermissionVehiclesRead, models.PermissionVehiclesUpdate, models.PermissionVehiclesDelete,
				models.PermissionRentalsCreate, models.PermissionRentalsRead, models.PermissionRentalsUpdate, models.PermissionRentalsApprove,
				models.PermissionUsersRead,
				models.PermissionLocationsRead,
				models.PermissionReportsView,
			},
		},
		{
			role: models.Role{
				Name:        models.RoleStaff,
				DisplayName: "Staff",
				Description: "Staff member - can create and manage rentals",
			},
			permissionNames: []string{
				models.PermissionVehiclesRead,
				models.PermissionRentalsCreate, models.PermissionRentalsRead, models.PermissionRentalsUpdate,
				models.PermissionLocationsRead,
			},
		},
		{
			role: models.Role{
				Name:        models.RoleCustomer,
				DisplayName: "Customer",
				Description: "Customer - can view vehicles and create own rentals",
			},
			permissionNames: []string{
				models.PermissionVehiclesRead,
				models.PermissionRentalsCreate, models.PermissionRentalsRead,
			},
		},
	}

	for _, roleData := range roles {
		// Check if role exists
		var existingRole models.Role
		result := db.Where("name = ?", roleData.role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new role
			role := roleData.role

			// Assign permissions
			if roleData.role.Name == models.RoleSuperAdmin {
				// Super admin gets all permissions
				role.Permissions = allPermissions
			} else {
				// Other roles get specific permissions
				var permissions []models.Permission
				for _, permName := range roleData.permissionNames {
					var perm models.Permission
					if err := db.Where("name = ?", permName).First(&perm).Error; err == nil {
						permissions = append(permissions, perm)
					}
				}
				role.Permissions = permissions
			}

			if err := db.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created role: %s", role.Name)
		} else {
			log.Printf("Role already exists: %s", roleData.role.Name)
		}
	}

	log.Println("Roles seeded successfully")
	return nil
}

func seedSuperAdmin(db *gorm.DB) error {
	log.Println("Seeding super admin user...")

	// Check if super admin already exists
	var existingUser models.User
	result := db.Where("email = ?", "admin@fleetpass.com").First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		// Hash password
		hashedPassword, err := auth.HashPassword("Admin123!")
		if err != nil {
			return err
		}

		// Get super admin role
		var superAdminRole models.Role
		if err := db.Where("name = ?", models.RoleSuperAdmin).First(&superAdminRole).Error; err != nil {
			return err
		}

		// Create super admin user
		user := models.User{
			Email:         "admin@fleetpass.com",
			Password:      hashedPassword,
			FirstName:     "System",
			LastName:      "Administrator",
			EmailVerified: true,
			IsActive:      true,
			Roles:         []models.Role{superAdminRole},
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		log.Println("✅ Super admin user created!")
		log.Println("   Email: admin@fleetpass.com")
		log.Println("   Password: Admin123!")
		log.Println("   ⚠️  CHANGE THIS PASSWORD IMMEDIATELY IN PRODUCTION!")
	} else {
		log.Println("Super admin user already exists")
	}

	return nil
}
