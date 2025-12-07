package handlers

import (
	"encoding/json"
	"fleetpass/internal/auth"
	"fleetpass/internal/database"
	"fleetpass/internal/email"
	"fleetpass/internal/models"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth
var emailService = email.GetEmailService()

// InitTokenAuth initializes the JWT auth instance
func InitTokenAuth(ta *jwtauth.JWTAuth) {
	tokenAuth = ta
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Email, password, first name, and last name are required", http.StatusBadRequest)
		return
	}

	// Validate password
	if err := auth.ValidatePassword(req.Password, auth.DefaultPasswordRequirements()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Generate verification token
	verificationToken, verificationExpiry, err := auth.GenerateVerificationToken()
	if err != nil {
		http.Error(w, "Error generating verification token", http.StatusInternalServerError)
		return
	}

	// Get default customer role
	var customerRole models.Role
	if err := database.DB.Where("name = ?", models.RoleCustomer).First(&customerRole).Error; err != nil {
		http.Error(w, "Error assigning default role", http.StatusInternalServerError)
		return
	}

	// Create user
	user := models.User{
		Email:              req.Email,
		Password:           hashedPassword,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Phone:              req.Phone,
		EmailVerified:      false,
		VerificationToken:  verificationToken,
		VerificationExpiry: &verificationExpiry,
		IsActive:           true,
		Roles:              []models.Role{customerRole},
	}

	// If organization ID provided (invitation-based), assign it
	if req.OrganizationID != "" {
		user.OrganizationID = &req.OrganizationID
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Send verification email
	if err := emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
		// Log error but don't fail registration
		println("Error sending verification email:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registration successful. Please check your email to verify your account.",
		"user_id": user.ID,
	})
}

// VerifyEmail handles email verification
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by verification token
	var user models.User
	if err := database.DB.Where("verification_token = ?", req.Token).Preload("Roles.Permissions").First(&user).Error; err != nil {
		http.Error(w, "Invalid or expired verification token", http.StatusBadRequest)
		return
	}

	// Check if token expired
	if auth.IsTokenExpired(user.VerificationExpiry) {
		http.Error(w, "Verification token has expired. Please request a new one.", http.StatusBadRequest)
		return
	}

	// Update user
	user.EmailVerified = true
	user.VerificationToken = ""
	user.VerificationExpiry = nil

	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	// Send welcome email
	emailService.SendWelcomeEmail(user.Email, user.FirstName)

	// Generate JWT token and log the user in
	token, err := generateJWTToken(&user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Email verified successfully",
		"token":   token,
		"user":    buildUserProfile(&user),
	})
}

// ForgotPassword handles password reset requests
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Always return success (don't reveal if email exists)
	successMessage := "If an account exists with this email, you will receive a password reset link."

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// User not found, but return success message for security
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": successMessage})
		return
	}

	// Generate reset token
	resetToken, resetExpiry, err := auth.GenerateResetToken()
	if err != nil {
		http.Error(w, "Error generating reset token", http.StatusInternalServerError)
		return
	}

	// Update user with reset token
	user.ResetToken = resetToken
	user.ResetTokenExpiry = &resetExpiry

	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	// Send reset email
	if err := emailService.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		println("Error sending reset email:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": successMessage})
}

// ResetPassword handles password reset
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by reset token
	var user models.User
	if err := database.DB.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	// Check if token expired
	if auth.IsTokenExpired(user.ResetTokenExpiry) {
		http.Error(w, "Reset token has expired. Please request a new one.", http.StatusBadRequest)
		return
	}

	// Validate new password
	if err := auth.ValidatePassword(req.NewPassword, auth.DefaultPasswordRequirements()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Update user
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error resetting password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successful. You can now log in with your new password.",
	})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by email and preload roles/permissions
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).Preload("Roles.Permissions").Preload("Organization").First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if email verified
	if !user.EmailVerified {
		http.Error(w, "Please verify your email address before logging in", http.StatusUnauthorized)
		return
	}

	// Check if account active
	if !user.IsActive {
		http.Error(w, "Your account has been deactivated. Please contact support.", http.StatusUnauthorized)
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	// Generate JWT token
	token, err := generateJWTToken(&user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.LoginResponse{
		Token: token,
		User:  buildUserProfile(&user),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile returns the current user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	// Get user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Fetch user from database with roles/permissions
	var user models.User
	if err := database.DB.Where("id = ?", userID).Preload("Roles.Permissions").Preload("Organization").First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	profile := buildUserProfile(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// Helper functions

func generateJWTToken(user *models.User) (string, error) {
	// Get role names
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Build claims
	claims := map[string]interface{}{
		"user_id":         user.ID,
		"email":           user.Email,
		"roles":           roleNames,
		"permissions":     user.GetPermissions(),
		"organization_id": user.OrganizationID,
	}

	jwtauth.SetExpiryIn(claims, 24*time.Hour)
	_, tokenString, err := tokenAuth.Encode(claims)
	return tokenString, err
}

func buildUserProfile(user *models.User) models.UserProfile {
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	return models.UserProfile{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		EmailVerified:  user.EmailVerified,
		IsActive:       user.IsActive,
		Roles:          roleNames,
		Permissions:    user.GetPermissions(),
		OrganizationID: user.OrganizationID,
	}
}
