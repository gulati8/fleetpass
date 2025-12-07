package auth

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Common weak passwords to reject
var commonPasswords = map[string]bool{
	"password":   true,
	"password1":  true,
	"123456":     true,
	"12345678":   true,
	"qwerty":     true,
	"abc123":     true,
	"monkey":     true,
	"1234567":    true,
	"letmein":    true,
	"trustno1":   true,
	"dragon":     true,
	"baseball":   true,
	"iloveyou":   true,
	"master":     true,
	"sunshine":   true,
	"ashley":     true,
	"bailey":     true,
	"passw0rd":   true,
	"shadow":     true,
	"123123":     true,
	"654321":     true,
	"superman":   true,
	"qazwsx":     true,
	"michael":    true,
	"football":   true,
}

// PasswordRequirements holds the password validation rules
type PasswordRequirements struct {
	MinLength          int
	RequireUppercase   bool
	RequireLowercase   bool
	RequireNumber      bool
	RequireSpecialChar bool
	RejectCommon       bool
}

// DefaultPasswordRequirements returns the default password requirements
func DefaultPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:          8,
		RequireUppercase:   true,
		RequireLowercase:   true,
		RequireNumber:      true,
		RequireSpecialChar: true,
		RejectCommon:       true,
	}
}

// ValidatePassword checks if a password meets all requirements
func ValidatePassword(password string, requirements PasswordRequirements) error {
	// Check minimum length
	if len(password) < requirements.MinLength {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for common passwords
	if requirements.RejectCommon && commonPasswords[strings.ToLower(password)] {
		return errors.New("password is too common, please choose a stronger password")
	}

	// Check for uppercase
	if requirements.RequireUppercase {
		matched, _ := regexp.MatchString(`[A-Z]`, password)
		if !matched {
			return errors.New("password must contain at least one uppercase letter")
		}
	}

	// Check for lowercase
	if requirements.RequireLowercase {
		matched, _ := regexp.MatchString(`[a-z]`, password)
		if !matched {
			return errors.New("password must contain at least one lowercase letter")
		}
	}

	// Check for number
	if requirements.RequireNumber {
		matched, _ := regexp.MatchString(`[0-9]`, password)
		if !matched {
			return errors.New("password must contain at least one number")
		}
	}

	// Check for special character
	if requirements.RequireSpecialChar {
		matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
		if !matched {
			return errors.New("password must contain at least one special character")
		}
	}

	return nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
