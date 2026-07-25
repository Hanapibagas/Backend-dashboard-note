package usecase

import (
	"auth/domain/errors"
	"auth/domain/repository"
	"auth/domain/service"
	"time"

	"github.com/google/uuid"
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
	tokenService            service.ITokenService
	userRegistrationService service.IUserRegistrationService
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
	tokenService service.ITokenService,
	userRegistrationService service.IUserRegistrationService,
) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:                userRepo,
		refreshTokenRepo:        refreshTokenRepo,
		tokenService:            tokenService,
		userRegistrationService: userRegistrationService,
	}
}

// Register registers a new user
func (uc *AuthUsecaseImpl) Register(request *RegisterRequest) (*RegisterResponse, error) {
	// Delegate to domain service for business logic
	user, err := uc.userRegistrationService.RegisterUser(request.Email, request.Password, request.FullName)
	if err != nil {
		// Simply return the domain error - let the handler layer decide HTTP status
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
		// Simply return the domain error - let the handler layer decide HTTP status
		return nil, err
	}

	// Generate access token
	accessToken, err := uc.tokenService.GenerateAccessToken(user.GetID(), user.GetEmail())
	if err != nil {
		return nil, errors.NewDomainError("TOKEN_GENERATION_FAILED", "failed to generate access token", err)
	}

	// Generate refresh token
	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.GetID(), user.GetEmail())
	if err != nil {
		return nil, errors.NewDomainError("TOKEN_GENERATION_FAILED", "failed to generate refresh token", err)
	}

	// Calculate expiry time
	expiresIn := int64(uc.tokenService.GetAccessTokenExpiry().Seconds())

	// Save refresh token to database
	expiryTime := time.Now().Add(uc.tokenService.GetRefreshTokenExpiry())
	err = uc.refreshTokenRepo.Save(user.GetID(), refreshToken, expiryTime.Unix())
	if err != nil {
		return nil, errors.NewDomainError("REFRESH_TOKEN_SAVE_FAILED", "failed to save refresh token", err)
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
		return errors.NewDomainError("INVALID_USER_ID", "invalid user ID format", err)
	}

	// Delete all refresh tokens for this user from database
	// This invalidates all refresh tokens, forcing user to login again
	err = uc.refreshTokenRepo.DeleteByUserID(userID)
	if err != nil {
		// Return error even if it fails - this is a hard logout
		return errors.NewDomainError("REFRESH_TOKEN_DELETE_FAILED", "failed to delete refresh tokens", err)
	}

	return nil
}
