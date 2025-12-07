package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateToken generates a random token for email verification or password reset
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a token for email verification
// Returns token and expiry time (24 hours from now)
func GenerateVerificationToken() (string, time.Time, error) {
	token, err := GenerateToken(32) // 64 character hex string
	if err != nil {
		return "", time.Time{}, err
	}
	expiry := time.Now().Add(24 * time.Hour)
	return token, expiry, nil
}

// GenerateResetToken generates a token for password reset
// Returns token and expiry time (1 hour from now)
func GenerateResetToken() (string, time.Time, error) {
	token, err := GenerateToken(32) // 64 character hex string
	if err != nil {
		return "", time.Time{}, err
	}
	expiry := time.Now().Add(1 * time.Hour)
	return token, expiry, nil
}

// IsTokenExpired checks if a token has expired
func IsTokenExpired(expiry *time.Time) bool {
	if expiry == nil {
		return true
	}
	return time.Now().After(*expiry)
}
