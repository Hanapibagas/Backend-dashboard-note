package response

import (
	"auth/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ErrorResponse represents a detailed error response
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Success sends a successful response with data
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}

// ErrorWithData sends an error response with additional error data
func ErrorWithData(c *gin.Context, statusCode int, message string, errData interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errData,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errors []utils.ValidationError) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Validation failed",
		Error: ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: utils.FormatValidationErrors(errors),
			Details: errors,
		},
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// UnauthorizedWithMessage sends a 401 response with custom message
func UnauthorizedWithMessage(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// ValidationErrorWithMessage sends a validation error with custom message
func ValidationErrorWithMessage(c *gin.Context, message string, errors []utils.ValidationError) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error: ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: utils.FormatValidationErrors(errors),
			Details: errors,
		},
	})
}
