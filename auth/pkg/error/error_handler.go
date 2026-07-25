package error

import (
	errorDomain "auth/domain/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler handles application errors and converts them to HTTP responses
// This is the SRP-correct approach: Single responsibility for error translation
type ErrorHandler struct{}

// NewErrorHandler creates a new ErrorHandler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError is a generic error handler that handles all domain errors
// This method uses errors.Is() for type-safe error checking (SRP & DDD)
func (eh *ErrorHandler) HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Handle domain errors using errors.Is() for type-safe checking
	// This is the SRP-correct approach: No string comparison

	// Check for specific error types
	switch {
	case errors.Is(err, errorDomain.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already exists",
		})
		return

	case errors.Is(err, errorDomain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return

	case errors.Is(err, errorDomain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return

	case errors.Is(err, errorDomain.ErrInvalidEmailFormat),
		errors.Is(err, errorDomain.ErrEmailEmpty):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return

	case errors.Is(err, errorDomain.ErrPasswordEmpty),
		errors.Is(err, errorDomain.ErrPasswordTooShort),
		errors.Is(err, errorDomain.ErrPasswordMissingUppercase),
		errors.Is(err, errorDomain.ErrPasswordMissingLowercase),
		errors.Is(err, errorDomain.ErrPasswordMissingNumber):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return

	case errors.Is(err, errorDomain.ErrPasswordMismatch):
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid password",
		})
		return

	default:
		// Generic error for unknown errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}
}

// HandleRegisterError handles errors from the registration usecase
// This method is kept for backward compatibility but delegates to HandleError
func (eh *ErrorHandler) HandleRegisterError(err error, c *gin.Context) {
	eh.HandleError(c, err)
}

// HandleLoginError handles errors from the login usecase
// This method is kept for backward compatibility but delegates to HandleError
func (eh *ErrorHandler) HandleLoginError(err error, c *gin.Context) {
	eh.HandleError(c, err)
}

// HandleLogoutError handles errors from the logout usecase
// This method is kept for backward compatibility but delegates to HandleError
func (eh *ErrorHandler) HandleLogoutError(err error, c *gin.Context) {
	eh.HandleError(c, err)
}
