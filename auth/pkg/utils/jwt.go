package utils

import (
	"errors"
	"fmt"
	"time"

	"auth/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// ErrInvalidToken is returned when token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken is returned when token is expired
	ErrExpiredToken = errors.New("token expired")

	// ErrTokenMalformed is returned when token format is malformed
	ErrTokenMalformed = errors.New("token malformed")
)

// Claims represents JWT custom claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secret             string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:             cfg.JWT.Secret,
		accessTokenExpiry:  time.Duration(cfg.JWT.AccessTokenExpiry) * time.Second,
		refreshTokenExpiry: time.Duration(cfg.JWT.RefreshTokenExpiry) * time.Second,
	}
}

// GenerateAccessToken generates a new access token for a user
func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token for a user
func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenMalformed
	}

	return claims, nil
}

// ExtractUserIDFromToken extracts user ID from token string without full validation
// This is useful for logging or debugging purposes
func (j *JWTManager) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

// GetAccessTokenExpiry returns the access token expiry duration
func (j *JWTManager) GetAccessTokenExpiry() time.Duration {
	return j.accessTokenExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (j *JWTManager) GetRefreshTokenExpiry() time.Duration {
	return j.refreshTokenExpiry
}
