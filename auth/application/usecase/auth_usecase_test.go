package usecase_test

import (
	"auth/application/usecase"
	"auth/domain/entity"
	errorDomain "auth/domain/errors"
	repoMocks "auth/domain/repository/mocks"
	serviceMocks "auth/domain/service/mocks"
	"auth/domain/valueobject"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestNewAuthUsecase tests the constructor
func TestNewAuthUsecase(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	assert.NotNil(t, uc)
}

// TestRegister_Success tests successful registration
func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Create test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, fullName)
	require.NoError(t, err)

	// Mock RegisterUser to return test user
	mockUserRegService.EXPECT().RegisterUser(email, password, fullName).
		Return(testUser, nil)

	request := &usecase.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}

	response, err := uc.Register(request)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, testUser.GetID().String(), response.UserID)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, fullName, response.FullName)
	assert.False(t, response.CreatedAt.IsZero())
}

// TestRegister_RegisterUserError tests registration with domain service error
func TestRegister_RegisterUserError(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Mock RegisterUser to return error
	mockUserRegService.EXPECT().RegisterUser(email, password, fullName).
		Return(nil, errorDomain.ErrEmailAlreadyExists)

	request := &usecase.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}

	response, err := uc.Register(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrEmailAlreadyExists))
}

// TestRegister_InvalidEmail tests registration with invalid email
func TestRegister_InvalidEmail(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "invalid-email"
	password := "Password123"
	fullName := "Test User"

	// Mock RegisterUser to return email validation error
	mockUserRegService.EXPECT().RegisterUser(email, password, fullName).
		Return(nil, errorDomain.ErrInvalidEmailFormat)

	request := &usecase.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}

	response, err := uc.Register(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrInvalidEmailFormat))
}

// TestRegister_InvalidPassword tests registration with invalid password
func TestRegister_InvalidPassword(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "weak"
	fullName := "Test User"

	// Mock RegisterUser to return password validation error
	mockUserRegService.EXPECT().RegisterUser(email, password, fullName).
		Return(nil, errorDomain.ErrPasswordTooShort)

	request := &usecase.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}

	response, err := uc.Register(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordTooShort))
}

// TestRegister_PasswordEmpty tests registration with empty password
func TestRegister_PasswordEmpty(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := ""
	fullName := "Test User"

	// Mock RegisterUser to return password empty error
	mockUserRegService.EXPECT().RegisterUser(email, password, fullName).
		Return(nil, errorDomain.ErrPasswordEmpty)

	request := &usecase.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}

	response, err := uc.Register(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordEmpty))
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"

	// Create test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	accessToken := "access.token.here"
	refreshToken := "refresh.token.here"
	accessExpiry := 15 * time.Minute

	// Mock AuthenticateUser
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(testUser, nil)

	// Mock GenerateAccessToken
	mockTokenService.EXPECT().GenerateAccessToken(testUser.GetID(), email).
		Return(accessToken, nil)

	// Mock GenerateRefreshToken
	mockTokenService.EXPECT().GenerateRefreshToken(testUser.GetID(), email).
		Return(refreshToken, nil)

	// Mock GetAccessTokenExpiry
	mockTokenService.EXPECT().GetAccessTokenExpiry().
		Return(accessExpiry)

	// Mock GetRefreshTokenExpiry
	mockTokenService.EXPECT().GetRefreshTokenExpiry().
		Return(7 * 24 * time.Hour)

	// Mock Save refresh token
	mockRefreshTokenRepo.EXPECT().Save(testUser.GetID(), refreshToken, mock.AnythingOfType("int64")).
		Return(nil)

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, testUser.GetID().String(), response.User.UserID)
	assert.Equal(t, email, response.User.Email)
	assert.Equal(t, "Test User", response.User.FullName)
	assert.Equal(t, accessToken, response.Tokens.AccessToken)
	assert.Equal(t, refreshToken, response.Tokens.RefreshToken)
	assert.Equal(t, int64(900), response.Tokens.ExpiresIn) // 15 minutes in seconds
}

// TestLogin_InvalidCredentials tests login with invalid credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "WrongPassword123"

	// Mock AuthenticateUser to return invalid credentials error
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(nil, errorDomain.ErrInvalidCredentials)

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrInvalidCredentials))
}

// TestLogin_UserNotFound tests login with non-existent user
func TestLogin_UserNotFound(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "nonexistent@example.com"
	password := "Password123"

	// Mock AuthenticateUser to return invalid credentials (user not found maps to this)
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(nil, errorDomain.ErrInvalidCredentials)

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrInvalidCredentials))
}

// TestLogin_AccessTokenGenerationError tests login with access token generation error
func TestLogin_AccessTokenGenerationError(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"

	// Create test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	// Mock AuthenticateUser
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(testUser, nil)

	// Mock GenerateAccessToken to return error
	mockTokenService.EXPECT().GenerateAccessToken(testUser.GetID(), email).
		Return("", errors.New("token generation failed"))

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "TOKEN_GENERATION_FAILED", domainErr.Code)
}

// TestLogin_RefreshTokenGenerationError tests login with refresh token generation error
func TestLogin_RefreshTokenGenerationError(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"

	// Create test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	accessToken := "access.token.here"

	// Mock AuthenticateUser
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(testUser, nil)

	// Mock GenerateAccessToken
	mockTokenService.EXPECT().GenerateAccessToken(testUser.GetID(), email).
		Return(accessToken, nil)

	// Mock GenerateRefreshToken to return error
	mockTokenService.EXPECT().GenerateRefreshToken(testUser.GetID(), email).
		Return("", errors.New("refresh token generation failed"))

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "TOKEN_GENERATION_FAILED", domainErr.Code)
}

// TestLogin_SaveRefreshTokenError tests login with save refresh token error
func TestLogin_SaveRefreshTokenError(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := "Password123"

	// Create test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	accessToken := "access.token.here"
	refreshToken := "refresh.token.here"
	accessExpiry := 15 * time.Minute

	// Mock AuthenticateUser
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(testUser, nil)

	// Mock GenerateAccessToken
	mockTokenService.EXPECT().GenerateAccessToken(testUser.GetID(), email).
		Return(accessToken, nil)

	// Mock GenerateRefreshToken
	mockTokenService.EXPECT().GenerateRefreshToken(testUser.GetID(), email).
		Return(refreshToken, nil)

	// Mock GetAccessTokenExpiry
	mockTokenService.EXPECT().GetAccessTokenExpiry().
		Return(accessExpiry)

	// Mock GetRefreshTokenExpiry
	mockTokenService.EXPECT().GetRefreshTokenExpiry().
		Return(7 * 24 * time.Hour)

	// Mock Save refresh token to return error
	mockRefreshTokenRepo.EXPECT().Save(testUser.GetID(), refreshToken, mock.AnythingOfType("int64")).
		Return(errors.New("database error"))

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "REFRESH_TOKEN_SAVE_FAILED", domainErr.Code)
}

// TestLogin_EmptyEmail tests login with empty email
func TestLogin_EmptyEmail(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := ""
	password := "Password123"

	// Mock AuthenticateUser to return email empty error
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(nil, errorDomain.ErrEmailEmpty)

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrEmailEmpty))
}

// TestLogin_EmptyPassword tests login with empty password
func TestLogin_EmptyPassword(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	email := "test@example.com"
	password := ""

	// Mock AuthenticateUser to return password empty error
	mockUserRegService.EXPECT().AuthenticateUser(email, password).
		Return(nil, errorDomain.ErrPasswordEmpty)

	request := &usecase.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := uc.Login(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordEmpty))
}

// TestLogout_Success tests successful logout
func TestLogout_Success(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	userID := uuid.New()
	refreshToken := "refresh.token.here"

	// Mock DeleteByUserID
	mockRefreshTokenRepo.EXPECT().DeleteByUserID(userID).
		Return(nil)

	request := &usecase.LogoutRequest{
		UserID:       userID.String(),
		RefreshToken: refreshToken,
	}

	err := uc.Logout(request)

	assert.NoError(t, err)
}

// TestLogout_InvalidUserID tests logout with invalid user ID
func TestLogout_InvalidUserID(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	invalidUserID := "invalid-uuid"
	refreshToken := "refresh.token.here"

	request := &usecase.LogoutRequest{
		UserID:       invalidUserID,
		RefreshToken: refreshToken,
	}

	err := uc.Logout(request)

	assert.Error(t, err)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "INVALID_USER_ID", domainErr.Code)
}

// TestLogout_EmptyUserID tests logout with empty user ID
func TestLogout_EmptyUserID(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	emptyUserID := ""
	refreshToken := "refresh.token.here"

	request := &usecase.LogoutRequest{
		UserID:       emptyUserID,
		RefreshToken: refreshToken,
	}

	err := uc.Logout(request)

	assert.Error(t, err)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "INVALID_USER_ID", domainErr.Code)
}

// TestLogout_DeleteError tests logout with delete refresh token error
func TestLogout_DeleteError(t *testing.T) {
	mockUserRepo := new(repoMocks.UserRepositoryMock)
	mockRefreshTokenRepo := new(repoMocks.RefreshTokenRepositoryMock)
	mockTokenService := new(serviceMocks.TokenServiceMock)
	mockUserRegService := new(serviceMocks.UserRegistrationServiceMock)

	uc := usecase.NewAuthUsecase(
		mockUserRepo,
		mockRefreshTokenRepo,
		mockTokenService,
		mockUserRegService,
	)

	userID := uuid.New()
	refreshToken := "refresh.token.here"

	// Mock DeleteByUserID to return error
	mockRefreshTokenRepo.EXPECT().DeleteByUserID(userID).
		Return(errors.New("database connection failed"))

	request := &usecase.LogoutRequest{
		UserID:       userID.String(),
		RefreshToken: refreshToken,
	}

	err := uc.Logout(request)

	assert.Error(t, err)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "REFRESH_TOKEN_DELETE_FAILED", domainErr.Code)
}
