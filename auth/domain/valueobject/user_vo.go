package valueobject

import (
	"auth/domain/errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (*Email, error) {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check if empty
	if email == "" {
		return nil, errors.ErrEmailEmpty
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, errors.ErrInvalidEmailFormat
	}

	// Convert to lowercase for consistency
	email = strings.ToLower(email)

	return &Email{value: email}, nil
}

// String returns the email string value
func (e *Email) String() string {
	return e.value
}

// Password represents a plain text password value object
// Responsibility: Validate password strength
type Password struct {
	value string
}

// NewPassword creates a new Password value object with validation
func NewPassword(password string) (*Password, error) {
	// Check if empty
	if password == "" {
		return nil, errors.ErrPasswordEmpty
	}

	// Check minimum length
	if len(password) < 8 {
		return nil, errors.ErrPasswordTooShort
	}

	// Check for uppercase
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return nil, errors.ErrPasswordMissingUppercase
	}

	// Check for lowercase
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return nil, errors.ErrPasswordMissingLowercase
	}

	// Check for number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return nil, errors.ErrPasswordMissingNumber
	}

	return &Password{value: password}, nil
}

// String returns the password string value
// WARNING: Use with caution, only for hashing purposes
func (p *Password) String() string {
	return p.value
}

// Hash converts the plain password to a hashed password
// This is the SRP-correct approach: Password VO handles its own hashing
func (p *Password) Hash() (*HashedPassword, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &HashedPassword{value: string(hashedBytes)}, nil
}

// HashedPassword represents a hashed password value object
// Responsibility: Encapsulate hashed password with verification capability
type HashedPassword struct {
	value string
}

// NewHashedPassword creates a HashedPassword from an existing hash
// Use this when loading from database
func NewHashedPassword(hash string) *HashedPassword {
	return &HashedPassword{value: hash}
}

// String returns the hashed password string
func (hp *HashedPassword) String() string {
	return hp.value
}

// VerifyPassword checks if the given plain password matches the hashed password
func (hp *HashedPassword) VerifyPassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hp.value), []byte(plainPassword))
	if err != nil {
		return errors.ErrPasswordMismatch
	}
	return nil
}
