package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing
	// 10 is a good balance between security and performance
	BcryptCost = 10
)

var (
	// ErrEmptyPassword is returned when password is empty
	ErrEmptyPassword = fmt.Errorf("password cannot be empty")
)

// PasswordManager handles password hashing and verification
type PasswordManager struct{}

// NewPasswordManager creates a new password manager
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword hashes a plain text password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a plain text password with a hashed password
func (pm *PasswordManager) ComparePassword(hashedPassword, plainPassword string) (bool, error) {
	if plainPassword == "" {
		return false, ErrEmptyPassword
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to compare password: %w", err)
	}

	return true, nil
}

// ValidatePasswordStrength validates password strength before hashing
// Returns true if password meets requirements, false otherwise
func (pm *PasswordManager) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Add more validation rules if needed
	// This is a basic validation, for more comprehensive validation
	// use the valueobject.Password

	return nil
}
