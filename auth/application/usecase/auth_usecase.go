package usecase

import (
	"auth/domain/repository"
	"auth/domain/service"
	"auth/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrEmailAlreadyExists is returned when email is already registered
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidCredentials is returned when email or password is incorrect
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidRefreshToken is returned when refresh token is invalid
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// AuthUsecase defines the interface for authentication business logic
type AuthUsecase interface {
	Register(request *RegisterRequest) (*RegisterResponse, error)
	Login(request *LoginRequest) (*LoginResponse, error)
	Logout(request *LogoutRequest) error
}

// AuthUsecaseImpl implements the AuthUsecase interface
type AuthUsecaseImpl struct {
	userRepo                repository.UserRepository
	refreshTokenRepo        repository.RefreshTokenRepository
	jwtManager              *utils.JWTManager
	userRegistrationService *service.UserRegistrationService
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string
	Password string
	FullName string
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents the login response
type LoginResponse struct {
	User   UserInfo  `json:"user"`
	Tokens TokenInfo `json:"tokens"`
}

// UserInfo represents user information
type UserInfo struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

// TokenInfo represents token information
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LogoutRequest represents the logout request
type LogoutRequest struct {
	UserID       string
	RefreshToken string
}

// NewAuthUsecase creates a new AuthUsecase implementation
func NewAuthUsecase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
	userRegistrationService *service.UserRegistrationService,
) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:                userRepo,
		refreshTokenRepo:        refreshTokenRepo,
		jwtManager:              jwtManager,
		userRegistrationService: userRegistrationService,
	}
}

// Register registers a new user
func (uc *AuthUsecaseImpl) Register(request *RegisterRequest) (*RegisterResponse, error) {
	// Delegate to domain service for business logic
	user, err := uc.userRegistrationService.RegisterUser(request.Email, request.Password, request.FullName)
	if err != nil {
		// Convert domain errors to usecase errors
		if err.Error() == "email already exists" {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	// Prepare response
	response := &RegisterResponse{
		UserID:    user.GetID().String(),
		Email:     user.GetEmail(),
		FullName:  user.GetFullName(),
		CreatedAt: user.GetCreatedAt(),
	}

	return response, nil
}

// Login authenticates a user and returns tokens
func (uc *AuthUsecaseImpl) Login(request *LoginRequest) (*LoginResponse, error) {
	// Delegate to domain service for authentication
	user, err := uc.userRegistrationService.AuthenticateUser(request.Email, request.Password)
	if err != nil {
		// Convert domain errors to usecase errors
		if err.Error() == "invalid credentials" {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Generate access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.GetID(), user.GetEmail())
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.GetID(), user.GetEmail())
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Calculate expiry time
	expiresIn := int64(uc.jwtManager.GetAccessTokenExpiry().Seconds())

	// Save refresh token to database
	expiryTime := time.Now().Add(uc.jwtManager.GetRefreshTokenExpiry())
	err = uc.refreshTokenRepo.Save(user.GetID(), refreshToken, expiryTime.Unix())
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	// Prepare response
	response := &LoginResponse{
		User: UserInfo{
			UserID:    user.GetID().String(),
			Email:     user.GetEmail(),
			FullName:  user.GetFullName(),
			CreatedAt: user.GetCreatedAt().Format(time.RFC3339),
		},
		Tokens: TokenInfo{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
		},
	}

	return response, nil
}

// Logout logs out a user by invalidating the refresh token
func (uc *AuthUsecaseImpl) Logout(request *LogoutRequest) error {
	// Get user ID from request (provided by auth middleware)
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Delete all refresh tokens for this user from database
	// This invalidates all refresh tokens, forcing user to login again
	err = uc.refreshTokenRepo.DeleteByUserID(userID)
	if err != nil {
		// Log error but don't fail - the access token will expire naturally
		// This is a soft logout approach
		return errors.New("failed to delete refresh tokens")
	}

	return nil
}
