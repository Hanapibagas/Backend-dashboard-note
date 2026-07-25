package main

import (
	appUsecase "auth/application/usecase"
	"auth/delivery/handler"
	"auth/delivery/middleware"
	domainService "auth/domain/service"
	"auth/infrastructure/database"
	infraJWT "auth/infrastructure/jwt"
	repoImpl "auth/infrastructure/repository"
	"auth/pkg/config"
	errorHandler "auth/pkg/error"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connection
	err = database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	log.Println("Database connected successfully")

	// Initialize Repositories
	userRepo := repoImpl.NewUserRepository()
	refreshTokenRepo := repoImpl.NewRefreshTokenRepository()
	log.Println("Repositories initialized")

	// Initialize Domain Services
	userRegistrationService := domainService.NewUserRegistrationService(userRepo)
	log.Println("Domain services initialized")

	// Initialize Infrastructure Services
	tokenService := infraJWT.NewJWTTokenService(cfg)
	log.Println("Infrastructure services initialized")

	// Initialize Error Handler
	errHandler := errorHandler.NewErrorHandler()
	log.Println("Error handler initialized")

	// Initialize UseCase
	authUsecase := appUsecase.NewAuthUsecase(userRepo, refreshTokenRepo, tokenService, userRegistrationService)
	log.Println("UseCase initialized")

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authUsecase, errHandler)
	log.Println("Handlers initialized")

	// Setup Gin router
	router := setupRouter(cfg, tokenService, authHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Environment: %s", cfg.App.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter configures the Gin router with middleware and routes
func setupRouter(cfg *config.Config, tokenService domainService.ITokenService, authHandler *handler.AuthHandler) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())   // Recovery from panics
	router.Use(gin.Logger())     // Logger
	router.Use(corsMiddleware()) // CORS middleware

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterHandler)
			auth.POST("/login", authHandler.LoginHandler)
			auth.POST("/logout", middleware.NewAuthMiddleware(tokenService).RequireAuth(), authHandler.LogoutHandler) // Protected route - requires valid JWT
		}

		// Protected routes example (optional - for future use)
		protected := api.Group("/protected")
		protected.Use(middleware.NewAuthMiddleware(tokenService).RequireAuth())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID := middleware.GetUserID(c)
				email := middleware.GetEmail(c)
				c.JSON(200, gin.H{
					"message": "This is a protected route",
					"user_id": userID,
					"email":   email,
				})
			})
		}
	}

	return router
}

// corsMiddleware handles CORS for the API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed origins from environment or use default
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:5173", // Vite default port
		}

		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// init function to check for required environment variables
func init() {
	requiredEnvVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"JWT_SECRET",
	}

	missing := false
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Printf("Warning: Environment variable %s is not set", envVar)
			missing = true
		}
	}

	if missing {
		log.Println("Warning: Some environment variables are not set. Please check your .env file")
	}
}
