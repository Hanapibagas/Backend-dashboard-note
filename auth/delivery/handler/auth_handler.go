package handler

import (
	"auth/application/usecase"
	errorHandler "auth/pkg/error"
	"auth/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
// SRP: Handler only handles HTTP concerns (binding, response formatting)
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	errHandler  *errorHandler.ErrorHandler
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authUsecase usecase.AuthUsecase, errHandler *errorHandler.ErrorHandler) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		errHandler:  errHandler,
	}
}

// RegisterRequest represents the HTTP request body for registration
// SRP: This struct only handles HTTP binding (basic format validation)
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`     // Basic: required field
	Password string `json:"password" binding:"required"`  // Basic: required field
	FullName string `json:"full_name" binding:"required"` // Basic: required field
}

// LoginRequest represents the HTTP request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LogoutRequest represents the HTTP request body for logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RegisterHandler handles user registration
// SRP: Only handles HTTP concerns, delegates all business logic to domain
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate JSON request (basic format validation only)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Call usecase (domain handles all business validation)
	usecaseReq := &usecase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	usecaseResp, err := h.authUsecase.Register(usecaseReq)
	if err != nil {
		h.errHandler.HandleError(c, err)
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
// SRP: Only handles HTTP concerns, delegates all business logic to domain
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Bind and validate JSON request (basic format validation only)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Call usecase (domain handles all business validation)
	usecaseReq := &usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	usecaseResp, err := h.authUsecase.Login(usecaseReq)
	if err != nil {
		h.errHandler.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, "Login successful", gin.H{
		"user":   usecaseResp.User,
		"tokens": usecaseResp.Tokens,
	})
}

// LogoutHandler handles user logout
// SRP: Only handles HTTP concerns, delegates all business logic to domain
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
		h.errHandler.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, "Logged out successfully", nil)
}
