package handler

import (
	"auth/application/usecase"
	errorHandler "auth/pkg/error"
	"auth/pkg/response"
	"auth/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUsecase  usecase.AuthUsecase
	errorHandler *errorHandler.ErrorHandler
	validator    *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authUsecase usecase.AuthUsecase, errHandler *errorHandler.ErrorHandler) *AuthHandler {
	return &AuthHandler{
		authUsecase:  authUsecase,
		errorHandler: errHandler,
		validator:    validator.New(),
	}
}

// RegisterRequest represents the HTTP request body for registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest represents the HTTP request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LogoutRequest represents the HTTP request body for logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate using validator
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		response.ValidationError(c, validationErrors)
		return
	}

	// Call usecase
	usecaseReq := &usecase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	usecaseResp, err := h.authUsecase.Register(usecaseReq)
	if err != nil {
		h.errorHandler.HandleRegisterError(err, c)
		return
	}

	// Return success response
	response.Created(c, "User registered successfully", gin.H{
		"user": gin.H{
			"user_id":    usecaseResp.UserID,
			"email":      usecaseResp.Email,
			"full_name":  usecaseResp.FullName,
			"created_at": usecaseResp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate using validator
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		response.ValidationError(c, validationErrors)
		return
	}

	// Call usecase
	usecaseReq := &usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	usecaseResp, err := h.authUsecase.Login(usecaseReq)
	if err != nil {
		h.errorHandler.HandleLoginError(err, c)
		return
	}

	// Return success response
	response.OK(c, "Login successful", gin.H{
		"user":   usecaseResp.User,
		"tokens": usecaseResp.Tokens,
	})
}

// LogoutHandler handles user logout
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	var req LogoutRequest

	// Bind JSON request (optional)
	c.ShouldBindJSON(&req)

	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")

	// Call usecase
	usecaseReq := &usecase.LogoutRequest{
		UserID:       userID,
		RefreshToken: req.RefreshToken,
	}

	err := h.authUsecase.Logout(usecaseReq)
	if err != nil {
		h.errorHandler.HandleLogoutError(err, c)
		return
	}

	// Return success response
	response.OK(c, "Logged out successfully", nil)
}

// formatValidationErrors converts validator errors to our ValidationError format
func (h *AuthHandler) formatValidationErrors(err error) []utils.ValidationError {
	var validationErrors []utils.ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := e.Field()
			var message string

			switch e.Tag() {
			case "required":
				message = field + " is required"
			case "email":
				message = field + " must be a valid email address"
			case "min":
				message = field + " must be at least " + e.Param() + " characters"
			case "max":
				message = field + " must be at most " + e.Param() + " characters"
			default:
				message = field + " is invalid"
			}

			validationErrors = append(validationErrors, utils.ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return validationErrors
}
