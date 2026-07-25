package service

import (
	"time"

	"github.com/google/uuid"
)

// ITokenService defines token operations interface
// This is a domain service because token generation is a business concern in authentication
type ITokenService interface {
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GetAccessTokenExpiry() time.Duration
	GetRefreshTokenExpiry() time.Duration
}

// TokenClaims represents token claims in domain terms
type TokenClaims struct {
	UserID    uuid.UUID
	Email     string
	ExpiresAt time.Time
}
