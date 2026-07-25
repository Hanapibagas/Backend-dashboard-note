package repository

import (
	"auth/domain/entity"
	"auth/domain/valueobject"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access operations
// This interface will be implemented in the infrastructure layer using native SQL
//
// DDD: Repository interface lives in domain layer
// Ubiquitous Language: Methods use domain concepts (Email VO)
type UserRepository interface {
	// Create saves a new user to the database
	Create(user *entity.User) error

	// FindByEmail finds a user by email address (accepts Email VO)
	// This is the DDD-correct approach: Repository works with domain types
	FindByEmail(email *valueobject.Email) (*entity.User, error)

	// FindByID finds a user by ID
	FindByID(id uuid.UUID) (*entity.User, error)

	// ExistsByEmail checks if a user with the given email exists (accepts Email VO)
	// This is the DDD-correct approach: Repository works with domain types
	ExistsByEmail(email *valueobject.Email) (bool, error)

	// Update updates an existing user
	Update(user *entity.User) error

	// Delete deletes a user by ID
	Delete(id uuid.UUID) error
}

// RefreshTokenRepository defines the interface for refresh token operations
// This is optional depending on whether you want to implement refresh token mechanism
type RefreshTokenRepository interface {
	// Save saves a new refresh token to the database
	Save(userID uuid.UUID, token string, expiresAt int64) error

	// Find finds a refresh token by token string
	Find(token string) (*RefreshToken, error)

	// Delete deletes a refresh token by token string
	Delete(token string) error

	// DeleteByUserID deletes all refresh tokens for a specific user
	DeleteByUserID(userID uuid.UUID) error

	// DeleteExpired deletes all expired refresh tokens
	DeleteExpired() error
}

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt int64
	CreatedAt int64
}
