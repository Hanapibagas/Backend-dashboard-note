package service

import (
	"auth/domain/entity"
	errorDomain "auth/domain/errors"
	"auth/domain/repository"
	"auth/domain/valueobject"
	"errors"
)

// IUserRegistrationService defines the interface for user registration operations
// This interface allows for easy mocking in tests
type IUserRegistrationService interface {
	RegisterUser(email string, password string, fullName string) (*entity.User, error)
	AuthenticateUser(email string, password string) (*entity.User, error)
}

// UserRegistrationService handles domain logic for user registration
// This is a Domain Service - contains business logic that doesn't naturally fit in Entity or VO
type UserRegistrationService struct {
	userRepo repository.UserRepository
}

// NewUserRegistrationService creates a new UserRegistrationService
func NewUserRegistrationService(userRepo repository.UserRepository) IUserRegistrationService {
	return &UserRegistrationService{
		userRepo: userRepo,
	}
}

// RegisterUser registers a new user with all business rules
// This method uses Value Objects consistently (Ubiquitous Language)
func (s *UserRegistrationService) RegisterUser(email string, password string, fullName string) (*entity.User, error) {
	// Validate and create Email value object
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Validate and create Password value object
	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Hash the password (SRP: Password VO handles hashing)
	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return nil, errorDomain.WrapError(errorDomain.ErrFailedToCreateUser, err)
	}

	// Check if email already exists (domain business rule)
	// Use Email VO for domain operations (DDD & Ubiquitous Language)
	exists, err := s.userRepo.ExistsByEmail(emailVO)
	if err != nil {
		return nil, errorDomain.WrapError(errorDomain.ErrFailedToCheckEmail, err)
	}
	if exists {
		return nil, errorDomain.ErrEmailAlreadyExists
	}

	// Create user entity with value objects (DDD & Ubiquitous Language)
	user, err := entity.NewUser(emailVO, hashedPassword, fullName)
	if err != nil {
		return nil, errorDomain.WrapError(errorDomain.ErrFailedToCreateUser, err)
	}

	// Save user to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errorDomain.WrapError(errorDomain.ErrFailedToSaveUser, err)
	}

	return user, nil
}

// AuthenticateUser authenticates a user with email and password
// This method uses Value Objects consistently (Ubiquitous Language)
func (s *UserRegistrationService) AuthenticateUser(email string, password string) (*entity.User, error) {
	// Validate and create Email value object
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Find user by email (using Email VO - DDD & Ubiquitous Language)
	user, err := s.userRepo.FindByEmail(emailVO)
	if err != nil {
		if errors.Is(err, errorDomain.ErrUserNotFound) {
			return nil, errorDomain.ErrInvalidCredentials
		}
		return nil, errorDomain.WrapError(errorDomain.ErrFailedToFindUser, err)
	}

	// Verify password (delegates to HashedPassword VO)
	err = user.VerifyPassword(password)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil, errorDomain.ErrInvalidCredentials
	}

	return user, nil
}
