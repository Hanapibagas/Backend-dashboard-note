package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidEmailFormat is returned when email format is invalid
	ErrInvalidEmailFormat = errors.New("invalid email format")

	// ErrEmailEmpty is returned when email is empty
	ErrEmailEmpty = errors.New("email cannot be empty")

	// ErrPasswordEmpty is returned when password is empty
	ErrPasswordEmpty = errors.New("password cannot be empty")

	// ErrPasswordTooShort is returned when password is less than 8 characters
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")

	// ErrPasswordMissingUppercase is returned when password has no uppercase letter
	ErrPasswordMissingUppercase = errors.New("password must contain at least one uppercase letter")

	// ErrPasswordMissingLowercase is returned when password has no lowercase letter
	ErrPasswordMissingLowercase = errors.New("password must contain at least one lowercase letter")

	// ErrPasswordMissingNumber is returned when password has no number
	ErrPasswordMissingNumber = errors.New("password must contain at least one number")
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
		return nil, ErrEmailEmpty
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmailFormat
	}

	// Convert to lowercase for consistency
	email = strings.ToLower(email)

	return &Email{value: email}, nil
}

// String returns the email string value
func (e *Email) String() string {
	return e.value
}

// Value returns the underlying email value
func (e *Email) Value() string {
	return e.value
}

// Password represents a password value object
type Password struct {
	value string
}

// NewPassword creates a new Password value object with validation
func NewPassword(password string) (*Password, error) {
	// Check if empty
	if password == "" {
		return nil, ErrPasswordEmpty
	}

	// Check minimum length
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// Check for uppercase
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return nil, ErrPasswordMissingUppercase
	}

	// Check for lowercase
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return nil, ErrPasswordMissingLowercase
	}

	// Check for number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return nil, ErrPasswordMissingNumber
	}

	return &Password{value: password}, nil
}

// String returns the password string value
func (p *Password) String() string {
	return p.value
}

// Value returns the underlying password value
func (p *Password) Value() string {
	return p.value
}
