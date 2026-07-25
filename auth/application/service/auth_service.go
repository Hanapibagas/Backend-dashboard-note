package service

import (
	"auth/domain/repository"
	"auth/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrRefreshTokenNotFound is returned when refresh token is not found
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	// ErrRefreshTokenExpired is returned when refresh token is expired
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

// AuthService defines the interface for additional authentication services
type AuthService interface {
	RefreshAccessToken(refreshToken string) (*RefreshTokenResponse, error)
	CleanupExpiredTokens() error
	ValidateRefreshToken(refreshToken string) (*TokenValidationResponse, error)
}

// RefreshTokenResponse represents the response for refresh token operation
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// TokenValidationResponse represents the response for token validation
type TokenValidationResponse struct {
	IsValid  bool   `json:"is_valid"`
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	ExpiresAt int64 `json:"expires_at,omitempty"`
	Message  string `json:"message,omitempty"`
}

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	userRepo        repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager      *utils.JWTManager
}

// NewAuthService creates a new AuthService implementation
func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
) AuthService {
	return &AuthServiceImpl{
		userRepo:        userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:      jwtManager,
	}
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *AuthServiceImpl) RefreshAccessToken(refreshToken string) (*RefreshTokenResponse, error) {
	// Validate the refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if refresh token exists in database (if using refresh token storage)
	// tokenData, err := s.refreshTokenRepo.Find(refreshToken)
	// if err != nil {
	//     return nil, ErrRefreshTokenNotFound
	// }

	// Check if token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, ErrRefreshTokenExpired
	}

	// Verify user still exists
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrRefreshTokenNotFound
	}

	// Generate new access token
	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.GetID(), user.GetEmail())
	if err != nil {
		return nil, errors.New("failed to generate new access token")
	}

	// Calculate expiry time
	expiresIn := int64(s.jwtManager.GetAccessTokenExpiry().Seconds())

	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
		ExpiresIn:   expiresIn,
	}, nil
}

// CleanupExpiredTokens deletes all expired refresh tokens from the database
// This should be called periodically (e.g., via cron job)
func (s *AuthServiceImpl) CleanupExpiredTokens() error {
	err := s.refreshTokenRepo.DeleteExpired()
	if err != nil {
		return errors.New("failed to cleanup expired tokens")
	}
	return nil
}

// ValidateRefreshToken validates a refresh token and returns its information
func (s *AuthServiceImpl) ValidateRefreshToken(refreshToken string) (*TokenValidationResponse, error) {
	// Validate token structure and signature
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return &TokenValidationResponse{
			IsValid: false,
			Message: "Invalid or malformed token",
		}, nil
	}

	// Check if expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return &TokenValidationResponse{
			IsValid: false,
			Message: "Token has expired",
		}, nil
	}

	// Check if token exists in database (if using refresh token storage)
	// tokenData, err := s.refreshTokenRepo.Find(refreshToken)
	// if err != nil {
	//     return &TokenValidationResponse{
	//         IsValid:  false,
	//         Message:  "Token not found in database",
	//     }, nil
	// }

	// Verify user exists
	_, err = s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return &TokenValidationResponse{
			IsValid: false,
			Message: "User not found",
		}, nil
	}

	return &TokenValidationResponse{
		IsValid:   true,
		UserID:    claims.UserID.String(),
		Email:     claims.Email,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

// RevokeUserTokens revokes all refresh tokens for a specific user
// This is useful for security events like password change or suspicious activity
func (s *AuthServiceImpl) RevokeUserTokens(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	err = s.refreshTokenRepo.DeleteByUserID(userUUID)
	if err != nil {
		return errors.New("failed to revoke user tokens")
	}

	return nil
}
