package jwt

import (
	"errors"
	"fmt"
	"time"

	"auth/domain/service"
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

// jwtClaims represents JWT custom claims
type jwtClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTTokenService implements ITokenService using JWT
type JWTTokenService struct {
	secret             string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTTokenService creates a new JWT token service
func NewJWTTokenService(cfg *config.Config) service.ITokenService {
	return &JWTTokenService{
		secret:             cfg.JWT.Secret,
		accessTokenExpiry:  time.Duration(cfg.JWT.AccessTokenExpiry) * time.Second,
		refreshTokenExpiry: time.Duration(cfg.JWT.RefreshTokenExpiry) * time.Second,
	}
}

// GenerateAccessToken generates a new access token for a user
func (j *JWTTokenService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := jwtClaims{
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
func (j *JWTTokenService) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	claims := jwtClaims{
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
func (j *JWTTokenService) ValidateToken(tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenMalformed
	}

	// Convert to domain TokenClaims
	return &service.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// GetAccessTokenExpiry returns the access token expiry duration
func (j *JWTTokenService) GetAccessTokenExpiry() time.Duration {
	return j.accessTokenExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (j *JWTTokenService) GetRefreshTokenExpiry() time.Duration {
	return j.refreshTokenExpiry
}
