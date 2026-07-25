package errors

import (
	"fmt"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is allows error comparison with errors.Is()
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Predefined domain errors
var (
	// User-related errors
	ErrEmailAlreadyExists = &DomainError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "email already exists",
	}

	ErrUserNotFound = &DomainError{
		Code:    "USER_NOT_FOUND",
		Message: "user not found",
	}

	ErrInvalidCredentials = &DomainError{
		Code:    "INVALID_CREDENTIALS",
		Message: "invalid credentials",
	}

	ErrInvalidRefreshToken = &DomainError{
		Code:    "INVALID_REFRESH_TOKEN",
		Message: "invalid refresh token",
	}

	// Value Object validation errors
	ErrInvalidEmailFormat = &DomainError{
		Code:    "INVALID_EMAIL_FORMAT",
		Message: "invalid email format",
	}

	ErrEmailEmpty = &DomainError{
		Code:    "EMAIL_EMPTY",
		Message: "email cannot be empty",
	}

	ErrPasswordEmpty = &DomainError{
		Code:    "PASSWORD_EMPTY",
		Message: "password cannot be empty",
	}

	ErrPasswordTooShort = &DomainError{
		Code:    "PASSWORD_TOO_SHORT",
		Message: "password must be at least 8 characters",
	}

	ErrPasswordMissingUppercase = &DomainError{
		Code:    "PASSWORD_MISSING_UPPERCASE",
		Message: "password must contain at least one uppercase letter",
	}

	ErrPasswordMissingLowercase = &DomainError{
		Code:    "PASSWORD_MISSING_LOWERCASE",
		Message: "password must contain at least one lowercase letter",
	}

	ErrPasswordMissingNumber = &DomainError{
		Code:    "PASSWORD_MISSING_NUMBER",
		Message: "password must contain at least one number",
	}

	ErrPasswordMismatch = &DomainError{
		Code:    "PASSWORD_MISMATCH",
		Message: "password does not match",
	}

	// Repository errors
	ErrFailedToCheckEmail = &DomainError{
		Code:    "FAILED_TO_CHECK_EMAIL",
		Message: "failed to check email existence",
	}

	ErrFailedToCreateUser = &DomainError{
		Code:    "FAILED_TO_CREATE_USER",
		Message: "failed to create user",
	}

	ErrFailedToFindUser = &DomainError{
		Code:    "FAILED_TO_FIND_USER",
		Message: "failed to find user",
	}

	ErrFailedToSaveUser = &DomainError{
		Code:    "FAILED_TO_SAVE_USER",
		Message: "failed to save user",
	}
)

// NewDomainError creates a new domain error with underlying error
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError wraps an error with a domain error
func WrapError(domainErr *DomainError, err error) error {
	if err == nil {
		return nil
	}
	return &DomainError{
		Code:    domainErr.Code,
		Message: domainErr.Message,
		Err:     err,
	}
}
