package entity

import (
	"auth/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity in the domain layer
// This is a Rich Domain Model with behavior and value objects
type User struct {
	ID           uuid.UUID
	Email        *valueobject.Email          // Value Object - Ubiquitous Language
	PasswordHash *valueobject.HashedPassword // Value Object - SRP: Hashing in VO, not Entity
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new User entity with value objects
// This is the SRP-correct approach: Entity focuses on User data and behavior
// Password hashing is handled by Password value object
func NewUser(email *valueobject.Email, hashedPassword *valueobject.HashedPassword, fullName string) (*User, error) {
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

// UpdateFullName updates the user's full name
func (u *User) UpdateFullName(fullName string) {
	u.FullName = fullName
	u.UpdatedAt = time.Now()
}

// UpdatePassword updates the user's password with new hashed password
// This method accepts a HashedPassword VO (hashing already done by Password VO)
func (u *User) UpdatePassword(hashedPassword *valueobject.HashedPassword) {
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
}

// VerifyPassword verifies if the provided password matches the stored hash
// Delegates to HashedPassword VO (SRP: VO handles password verification)
func (u *User) VerifyPassword(plainPassword string) error {
	return u.PasswordHash.VerifyPassword(plainPassword)
}

// GetID returns the user ID
func (u *User) GetID() uuid.UUID {
	return u.ID
}

// GetEmail returns the user email as string
// Convenience method for infrastructure layer
func (u *User) GetEmail() string {
	return u.Email.String()
}

// GetEmailVO returns the user Email value object
// Preferred method for domain layer operations
func (u *User) GetEmailVO() *valueobject.Email {
	return u.Email
}

// GetPasswordHash returns the hashed password as string
// Convenience method for infrastructure layer (e.g., saving to DB)
func (u *User) GetPasswordHash() string {
	return u.PasswordHash.String()
}

// GetPasswordHashVO returns the HashedPassword value object
// Preferred method for domain layer operations
func (u *User) GetPasswordHashVO() *valueobject.HashedPassword {
	return u.PasswordHash
}

// GetFullName returns the user full name
func (u *User) GetFullName() string {
	return u.FullName
}

// GetCreatedAt returns the user creation time
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// GetUpdatedAt returns the last update time
func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
