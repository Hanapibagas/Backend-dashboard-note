package model

import (
	"auth/domain/repository"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenModel represents the refresh_tokens table structure in MySQL database
// This is a database model that maps directly to the refresh_tokens table
type RefreshTokenModel struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// TableName returns the table name for this model
func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

// ToDomain converts RefreshTokenModel to RefreshToken entity
func (rtm *RefreshTokenModel) ToDomain() *repository.RefreshToken {
	return &repository.RefreshToken{
		ID:        rtm.ID,
		UserID:    rtm.UserID,
		Token:     rtm.Token,
		ExpiresAt: rtm.ExpiresAt.Unix(),
		CreatedAt: rtm.CreatedAt.Unix(),
	}
}

// FromDomainRefreshToken creates RefreshTokenModel from RefreshToken entity
func FromDomainRefreshToken(token *repository.RefreshToken) RefreshTokenModel {
	return RefreshTokenModel{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: time.Unix(token.ExpiresAt, 0),
		CreatedAt: time.Unix(token.CreatedAt, 0),
	}
}
