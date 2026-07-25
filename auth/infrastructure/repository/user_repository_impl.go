package repository

import (
	"auth/domain/entity"
	"auth/domain/errors"
	"auth/domain/repository"
	"auth/domain/valueobject"
	"auth/infrastructure/database"
	"auth/infrastructure/database/model"
	"database/sql"
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

	// Use FromDomainUser to extract primitives from value objects
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
		return errors.WrapError(errors.ErrFailedToSaveUser, err)
	}

	return nil
}

// FindByEmail finds a user by email address (accepts Email VO)
func (r *userRepositoryImpl) FindByEmail(email *valueobject.Email) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	// Extract string from Email VO for database query
	emailStr := email.String()

	var userModel model.UserModel
	err := r.db.QueryRow(query, emailStr).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.FullName,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(errors.ErrFailedToFindUser, err)
	}

	// Convert UserModel to User entity with value objects
	user, err := userModel.ToDomain()
	if err != nil {
		return nil, errors.WrapError(errors.ErrFailedToCreateUser, err)
	}

	return user, nil
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
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(errors.ErrFailedToFindUser, err)
	}

	// Convert UserModel to User entity with value objects
	user, err := userModel.ToDomain()
	if err != nil {
		return nil, errors.WrapError(errors.ErrFailedToCreateUser, err)
	}

	return user, nil
}

// ExistsByEmail checks if a user with the given email exists (accepts Email VO)
func (r *userRepositoryImpl) ExistsByEmail(email *valueobject.Email) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	// Extract string from Email VO for database query
	emailStr := email.String()

	var count int
	err := r.db.QueryRow(query, emailStr).Scan(&count)
	if err != nil {
		return false, errors.WrapError(errors.ErrFailedToCheckEmail, err)
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

	// Use FromDomainUser to extract primitives from value objects
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
		return errors.WrapError(errors.ErrFailedToSaveUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(errors.ErrFailedToSaveUser, err)
	}

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return errors.WrapError(errors.ErrFailedToSaveUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(errors.ErrFailedToSaveUser, err)
	}

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
