package model

import (
	"auth/domain/entity"
	"auth/domain/errors"
	"auth/domain/valueobject"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// UserModel represents the user table structure in MySQL database
// This is a database model that maps directly to the users table
type UserModel struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName returns the table name for this model
func (UserModel) TableName() string {
	return "users"
}

// ToDomain converts UserModel to User entity
// This method reconstructs the entity with value objects from database primitives
func (um *UserModel) ToDomain() (*entity.User, error) {
	fullName := ""
	if um.FullName.Valid {
		fullName = um.FullName.String
	}

	// Reconstruct Email VO from database string
	emailVO, err := valueobject.NewEmail(um.Email)
	if err != nil {
		// If email from DB is invalid, wrap in domain error
		return nil, errors.WrapError(errors.ErrFailedToCreateUser, err)
	}

	// Reconstruct HashedPassword VO from database string
	hashedPasswordVO := valueobject.NewHashedPassword(um.PasswordHash)

	// Create User entity with value objects
	user, err := entity.NewUser(emailVO, hashedPasswordVO, fullName)
	if err != nil {
		return nil, errors.WrapError(errors.ErrFailedToCreateUser, err)
	}

	// Set ID and timestamps from database (NewUser creates new ones)
	// This is acceptable because we're reconstructing from persistence
	user = &entity.User{
		ID:           um.ID,
		Email:        emailVO,
		PasswordHash: hashedPasswordVO,
		FullName:     fullName,
		CreatedAt:    um.CreatedAt,
		UpdatedAt:    um.UpdatedAt,
	}

	return user, nil
}

// ToDomainSimple converts UserModel to User entity without error checking
// This is a convenience method for when you know the data is valid
// WARNING: Use only when you're certain the data is valid (e.g., from tests)
func (um *UserModel) ToDomainSimple() *entity.User {
	fullName := ""
	if um.FullName.Valid {
		fullName = um.FullName.String
	}

	// Reconstruct Email VO (assume valid from DB)
	emailVO, _ := valueobject.NewEmail(um.Email)
	hashedPasswordVO := valueobject.NewHashedPassword(um.PasswordHash)

	// Create User entity with value objects
	user := &entity.User{
		ID:           um.ID,
		Email:        emailVO,
		PasswordHash: hashedPasswordVO,
		FullName:     fullName,
		CreatedAt:    um.CreatedAt,
		UpdatedAt:    um.UpdatedAt,
	}

	return user
}

// FromDomain creates UserModel from User entity
// This method extracts primitives from entity and value objects for database storage
func FromDomainUser(user *entity.User) UserModel {
	// Extract primitives from value objects for database storage
	// This is the infrastructure layer's responsibility: adapt domain to persistence
	emailStr := user.GetEmail()              // Email VO → string
	passwordHashStr := user.GetPasswordHash() // HashedPassword VO → string

	return UserModel{
		ID:           user.ID,
		Email:        emailStr,
		PasswordHash: passwordHashStr,
		FullName: sql.NullString{
			String: user.FullName,
			Valid:  user.FullName != "",
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
