package repository

import (
	"auth/domain/entity"
	"auth/domain/repository"
	"auth/infrastructure/database"
	"auth/infrastructure/database/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// userRepositoryImpl implements the UserRepository interface using native SQL
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new User repository implementation
func NewUserRepository() repository.UserRepository {
	return &userRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create saves a new user to the database
func (r *userRepositoryImpl) Create(user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	userModel := model.FromDomainUser(user)
	userModel.CreatedAt = time.Now()
	userModel.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		userModel.ID,
		userModel.Email,
		userModel.PasswordHash,
		userModel.FullName,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByEmail finds a user by email address
func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var userModel model.UserModel
	err := r.db.QueryRow(query, email).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.FullName,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return userModel.ToDomain(), nil
}

// FindByID finds a user by ID
func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var userModel model.UserModel
	err := r.db.QueryRow(query, id).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.FullName,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return userModel.ToDomain(), nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// Update updates an existing user
func (r *userRepositoryImpl) Update(user *entity.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, full_name = ?, updated_at = ?
		WHERE id = ?
	`

	userModel := model.FromDomainUser(user)
	userModel.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		userModel.Email,
		userModel.PasswordHash,
		userModel.FullName,
		userModel.UpdatedAt,
		userModel.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
