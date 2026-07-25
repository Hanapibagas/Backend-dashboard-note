package middleware

import (
	"strings"

	"auth/domain/service"
	"auth/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader is the header key for Authorization
	AuthorizationHeader = "Authorization"

	// BearerPrefix is the prefix for Bearer tokens
	BearerPrefix = "Bearer "
)

// AuthMiddleware represents the authentication middleware
type AuthMiddleware struct {
	tokenService service.ITokenService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(tokenService service.ITokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// RequireAuth validates JWT token and sets user context
// This middleware should be used for protected routes
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if it has Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			message := "Invalid token"
			if err.Error() == "token expired" {
				message = "Token has expired"
			}
			response.Unauthorized(c, message)
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr
		}
	}
	return ""
}

// GetEmail retrieves email from context
func GetEmail(c *gin.Context) string {
	if email, exists := c.Get("email"); exists {
		if emailStr, ok := email.(string); ok {
			return emailStr
		}
	}
	return ""
}

// OptionalAuth validates JWT token if present, but doesn't require it
// This middleware is useful for routes that work with or without authentication
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			// No auth header, continue without setting user context
			c.Next()
			return
		}

		// Try to validate token if present
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			c.Next()
			return
		}

		// Validate token and set context if valid
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID.String())
			c.Set("email", claims.Email)
		}

		c.Next()
	}
}
