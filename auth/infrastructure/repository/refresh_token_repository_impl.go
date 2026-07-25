package repository

import (
	"auth/domain/repository"
	"auth/infrastructure/database"
	"auth/infrastructure/database/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// refreshTokenRepositoryImpl implements the RefreshTokenRepository interface using native SQL
type refreshTokenRepositoryImpl struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new RefreshToken repository implementation
func NewRefreshTokenRepository() repository.RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{
		db: database.GetDB(),
	}
}

// Save saves a new refresh token to the database
func (r *refreshTokenRepositoryImpl) Save(userID uuid.UUID, token string, expiresAt int64) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	tokenID := uuid.New()
	createdAt := time.Now()

	_, err := r.db.Exec(query, tokenID, userID, token, time.Unix(expiresAt, 0), createdAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// Find finds a refresh token by token string
func (r *refreshTokenRepositoryImpl) Find(token string) (*repository.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = ?
		LIMIT 1
	`

	var tokenModel model.RefreshTokenModel
	err := r.db.QueryRow(query, token).Scan(
		&tokenModel.ID,
		&tokenModel.UserID,
		&tokenModel.Token,
		&tokenModel.ExpiresAt,
		&tokenModel.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	return tokenModel.ToDomain(), nil
}

// Delete deletes a refresh token by token string
func (r *refreshTokenRepositoryImpl) Delete(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}

// DeleteByUserID deletes all refresh tokens for a specific user
func (r *refreshTokenRepositoryImpl) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens by user ID: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *refreshTokenRepositoryImpl) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}
