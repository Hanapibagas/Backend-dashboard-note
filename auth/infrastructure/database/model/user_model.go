package model

import (
	"auth/domain/entity"
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
func (um *UserModel) ToDomain() *entity.User {
	fullName := ""
	if um.FullName.Valid {
		fullName = um.FullName.String
	}

	user := &entity.User{
		ID:           um.ID,
		Email:        um.Email,
		PasswordHash: um.PasswordHash,
		FullName:     fullName,
		CreatedAt:    um.CreatedAt,
		UpdatedAt:    um.UpdatedAt,
	}

	return user
}

// FromDomain creates UserModel from User entity
func FromDomainUser(user *entity.User) UserModel {
	return UserModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName: sql.NullString{
			String: user.FullName,
			Valid:  user.FullName != "",
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
