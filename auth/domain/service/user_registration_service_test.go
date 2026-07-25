package service_test

import (
	"auth/domain/entity"
	errorDomain "auth/domain/errors"
	"auth/domain/repository/mocks"
	"auth/domain/service"
	"auth/domain/valueobject"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestNewUserRegistrationService tests the constructor
func TestNewUserRegistrationService(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	assert.NotNil(t, svc)
}

// TestRegisterUser_Success tests successful user registration
func TestRegisterUser_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Mock ExistsByEmail to return false (email doesn't exist)
	mockRepo.EXPECT().ExistsByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(false, nil)

	// Mock Create to succeed
	mockRepo.EXPECT().Create(mock.AnythingOfType("*entity.User")).
		Return(nil)

	user, err := svc.RegisterUser(email, password, fullName)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, fullName, user.GetFullName())
	assert.Equal(t, email, user.GetEmail())
}

// TestRegisterUser_InvalidEmail_Empty tests registration with empty email
func TestRegisterUser_InvalidEmail_Empty(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("", "Password123", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrEmailEmpty))
}

// TestRegisterUser_InvalidEmailFormat tests registration with invalid email format
func TestRegisterUser_InvalidEmailFormat(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"test@",
		"test@example",
	}

	for _, email := range invalidEmails {
		user, err := svc.RegisterUser(email, "Password123", "Test User")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, errorDomain.ErrInvalidEmailFormat))
	}
}

// TestRegisterUser_InvalidPassword_Empty tests registration with empty password
func TestRegisterUser_InvalidPassword_Empty(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("test@example.com", "", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordEmpty))
}

// TestRegisterUser_PasswordTooShort tests registration with short password
func TestRegisterUser_PasswordTooShort(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("test@example.com", "Pass1", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordTooShort))
}

// TestRegisterUser_PasswordMissingUppercase tests registration without uppercase
func TestRegisterUser_PasswordMissingUppercase(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("test@example.com", "password123", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordMissingUppercase))
}

// TestRegisterUser_PasswordMissingLowercase tests registration without lowercase
func TestRegisterUser_PasswordMissingLowercase(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("test@example.com", "PASSWORD123", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordMissingLowercase))
}

// TestRegisterUser_PasswordMissingNumber tests registration without number
func TestRegisterUser_PasswordMissingNumber(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.RegisterUser("test@example.com", "Passwordabc", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrPasswordMissingNumber))
}

// TestRegisterUser_EmailAlreadyExists tests registration with existing email
func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Mock ExistsByEmail to return true (email already exists)
	mockRepo.EXPECT().ExistsByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(true, nil)

	user, err := svc.RegisterUser(email, password, fullName)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrEmailAlreadyExists))
}

// TestRegisterUser_ExistsByEmailError tests registration with repository error on email check
func TestRegisterUser_ExistsByEmailError(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Mock ExistsByEmail to return an error
	mockRepo.EXPECT().ExistsByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(false, errors.New("database connection failed"))

	user, err := svc.RegisterUser(email, password, fullName)

	assert.Error(t, err)
	assert.Nil(t, user)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "FAILED_TO_CHECK_EMAIL", domainErr.Code)
}

// TestRegisterUser_CreateUserError tests registration with user creation error
func TestRegisterUser_CreateUserError(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	fullName := "" // Empty full name might cause issues in some implementations

	// Mock ExistsByEmail to return false
	mockRepo.EXPECT().ExistsByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(false, nil)

	// Mock Create to return an error
	mockRepo.EXPECT().Create(mock.AnythingOfType("*entity.User")).
		Return(errors.New("database error"))

	user, err := svc.RegisterUser(email, password, fullName)

	assert.Error(t, err)
	assert.Nil(t, user)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "FAILED_TO_SAVE_USER", domainErr.Code)
}

// TestRegisterUser_SaveUserError tests registration with database save error
func TestRegisterUser_SaveUserError(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	fullName := "Test User"

	// Mock ExistsByEmail to return false
	mockRepo.EXPECT().ExistsByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(false, nil)

	// Mock Create to return an error
	mockRepo.EXPECT().Create(mock.AnythingOfType("*entity.User")).
		Return(errors.New("database connection error"))

	user, err := svc.RegisterUser(email, password, fullName)

	assert.Error(t, err)
	assert.Nil(t, user)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "FAILED_TO_SAVE_USER", domainErr.Code)
}

// TestAuthenticateUser_Success tests successful authentication
func TestAuthenticateUser_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"

	// Create a hashed password for testing
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)

	// Create a test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	// Mock FindByEmail to return the test user
	mockRepo.EXPECT().FindByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(testUser, nil)

	user, err := svc.AuthenticateUser(email, password)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.GetEmail())
}

// TestAuthenticateUser_InvalidEmailFormat tests authentication with invalid email format
func TestAuthenticateUser_InvalidEmailFormat(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"test@",
	}

	for _, email := range invalidEmails {
		user, err := svc.AuthenticateUser(email, "Password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, errorDomain.ErrInvalidEmailFormat))
	}
}

// TestAuthenticateUser_InvalidEmail_Empty tests authentication with empty email
func TestAuthenticateUser_InvalidEmail_Empty(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	user, err := svc.AuthenticateUser("", "Password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrEmailEmpty))
}

// TestAuthenticateUser_UserNotFound tests authentication with non-existent user
func TestAuthenticateUser_UserNotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "nonexistent@example.com"
	password := "Password123"

	// Mock FindByEmail to return ErrUserNotFound
	mockRepo.EXPECT().FindByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(nil, errorDomain.ErrUserNotFound)

	user, err := svc.AuthenticateUser(email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrInvalidCredentials))
}

// TestAuthenticateUser_RepositoryError tests authentication with repository error
func TestAuthenticateUser_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"

	// Mock FindByEmail to return a database error
	mockRepo.EXPECT().FindByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(nil, errors.New("database connection failed"))

	user, err := svc.AuthenticateUser(email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "FAILED_TO_FIND_USER", domainErr.Code)
}

// TestAuthenticateUser_InvalidPassword tests authentication with wrong password
func TestAuthenticateUser_InvalidPassword(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"
	wrongPassword := "WrongPassword123"

	// Create a hashed password for testing
	hashedPassword, err := valueobject.NewPassword(password)
	require.NoError(t, err)
	hashedPasswordVO, err := hashedPassword.Hash()
	require.NoError(t, err)

	// Create a test user
	emailVO, err := valueobject.NewEmail(email)
	require.NoError(t, err)
	testUser, err := entity.NewUser(emailVO, hashedPasswordVO, "Test User")
	require.NoError(t, err)

	// Mock FindByEmail to return the test user
	mockRepo.EXPECT().FindByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(testUser, nil)

	user, err := svc.AuthenticateUser(email, wrongPassword)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, errorDomain.ErrInvalidCredentials))
}

// TestAuthenticateUser_FindByEmailError tests authentication with error from FindByEmail
func TestAuthenticateUser_FindByEmailError(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	svc := service.NewUserRegistrationService(mockRepo)

	email := "test@example.com"
	password := "Password123"

	// Mock FindByEmail to return a non-ErrUserNotFound error
	mockRepo.EXPECT().FindByEmail(mock.AnythingOfType("*valueobject.Email")).
		Return(nil, errors.New("unexpected error"))

	user, err := svc.AuthenticateUser(email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	var domainErr *errorDomain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, "FAILED_TO_FIND_USER", domainErr.Code)
}
