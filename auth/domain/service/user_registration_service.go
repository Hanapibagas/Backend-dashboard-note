package service

import (
	"auth/domain/entity"
	"auth/domain/repository"
	"auth/domain/valueobject"
	"errors"
)

// UserRegistrationService handles domain logic for user registration
type UserRegistrationService struct {
	userRepo repository.UserRepository
}

// NewUserRegistrationService creates a new UserRegistrationService
func NewUserRegistrationService(userRepo repository.UserRepository) *UserRegistrationService {
	return &UserRegistrationService{
		userRepo: userRepo,
	}
}

// RegisterUser registers a new user with all business rules
func (s *UserRegistrationService) RegisterUser(email string, password string, fullName string) (*entity.User, error) {
	// Validate email using value object
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Validate password using value object
	_, err = valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Check if email already exists (domain business rule)
	exists, err := s.userRepo.ExistsByEmail(emailVO.String())
	if err != nil {
		return nil, errors.New("failed to check email existence")
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create user entity (password hashing happens inside)
	user, err := entity.NewUser(emailVO.String(), password, fullName)
	if err != nil {
		return nil, errors.New("failed to create user entity")
	}

	// Save user to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("failed to save user")
	}

	return user, nil
}

// AuthenticateUser authenticates a user with email and password
func (s *UserRegistrationService) AuthenticateUser(email string, password string) (*entity.User, error) {
	// Validate email
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(emailVO.String())
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("failed to find user")
	}

	// Verify password
	err = user.VerifyPassword(password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
