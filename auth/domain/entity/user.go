package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user entity in the domain layer
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new User entity with hashed password
func NewUser(email string, password string, fullName string) (*User, error) {
	// Hash password as part of entity creation
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return user, nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// UpdateFullName updates the user's full name
func (u *User) UpdateFullName(fullName string) {
	u.FullName = fullName
	u.UpdatedAt = time.Now()
}

// UpdatePassword updates the user's password with hashing
func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// VerifyPassword verifies if the provided password matches the stored hash
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// GetID returns the user ID
func (u *User) GetID() uuid.UUID {
	return u.ID
}

// GetEmail returns the user email
func (u *User) GetEmail() string {
	return u.Email
}

// GetFullName returns the user full name
func (u *User) GetFullName() string {
	return u.FullName
}

// GetCreatedAt returns the user creation time
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}
