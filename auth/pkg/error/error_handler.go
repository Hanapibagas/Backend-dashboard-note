package error

import (
	"auth/application/usecase"
	"auth/domain/valueobject"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler handles application errors and converts them to HTTP responses
type ErrorHandler struct{}

// NewErrorHandler creates a new ErrorHandler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleRegisterError handles errors from the registration usecase
func (eh *ErrorHandler) HandleRegisterError(err error, c *gin.Context) {
	if err == nil {
		return
	}

	// Handle specific usecase errors
	if err == usecase.ErrEmailAlreadyExists {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already exists",
		})
		return
	}

	// Handle email validation errors
	if err == valueobject.ErrInvalidEmailFormat || err == valueobject.ErrEmailEmpty {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Handle password validation errors
	if err == valueobject.ErrPasswordTooShort || err == valueobject.ErrPasswordMissingUppercase ||
		err == valueobject.ErrPasswordMissingLowercase || err == valueobject.ErrPasswordMissingNumber ||
		err == valueobject.ErrPasswordEmpty {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generic error
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Failed to register user",
	})
}

// HandleLoginError handles errors from the login usecase
func (eh *ErrorHandler) HandleLoginError(err error, c *gin.Context) {
	if err == nil {
		return
	}

	// Handle specific usecase errors
	if err == usecase.ErrInvalidCredentials {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// Handle email validation errors
	if err == valueobject.ErrInvalidEmailFormat || err == valueobject.ErrEmailEmpty {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generic error
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Failed to login",
	})
}

// HandleLogoutError handles errors from the logout usecase
func (eh *ErrorHandler) HandleLogoutError(err error, c *gin.Context) {
	if err == nil {
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Failed to logout",
	})
}

// IsValidationError checks if an error is a validation error
func (eh *ErrorHandler) IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	if err == valueobject.ErrInvalidEmailFormat || err == valueobject.ErrEmailEmpty ||
		err == valueobject.ErrPasswordTooShort || err == valueobject.ErrPasswordMissingUppercase ||
		err == valueobject.ErrPasswordMissingLowercase || err == valueobject.ErrPasswordMissingNumber ||
		err == valueobject.ErrPasswordEmpty {
		return true
	}

	return false
}

// GetErrorMessage extracts the error message from an error
func (eh *ErrorHandler) GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

// IsConflictError checks if an error is a conflict error
func (eh *ErrorHandler) IsConflictError(err error) bool {
	return err == usecase.ErrEmailAlreadyExists
}

// IsUnauthorizedError checks if an error is an unauthorized error
func (eh *ErrorHandler) IsUnauthorizedError(err error) bool {
	return err == usecase.ErrInvalidCredentials
}

// IsNotFoundError checks if an error is a not found error
func (eh *ErrorHandler) IsNotFoundError(err error) bool {
	return err == usecase.ErrUserNotFound
}

// WrapError wraps an error with additional context
func (eh *ErrorHandler) WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.New(message + ": " + err.Error())
}
