# Auth Service

[![Go Version](https://img.shields.io/badge/Go-1.26.4-%2300ADD8?style=flat&logo=go)](https://golang.org)
[![Test-Driven Development](https://img.shields.io/badge/TDD-Test%20Driven%20Development-%2345B64D?style=flat)](https://en.wikipedia.org/wiki/Test-driven_development)
[![Domain Coverage](https://img.shields.io/badge/Domain%20Coverage-94.1%-%2328A745?style=flat)]()
[![Application Coverage](https://img.shields.io/badge/Application%20Coverage-100%-%2328A745?style=flat)]()

A robust, production-ready authentication service built with Go using **Clean Architecture** and **Domain-Driven Design (DDD)** principles. This service provides secure user authentication with JWT tokens, registration, login, and logout functionality.

## 🎯 Test-Driven Development (TDD)

This project is developed using **Test-Driven Development (TDD)** methodology, ensuring:

- ✅ **High Code Quality** - All critical business logic is thoroughly tested
- ✅ **Regression Prevention** - Comprehensive test coverage prevents future bugs
- ✅ **Living Documentation** - Tests serve as executable documentation
- ✅ **Refactoring Confidence** - Safe refactoring with test safety net
- ✅ **Domain & Application Coverage** - 90.3% overall coverage with all layers above 80%

### TDD Coverage Quality Standards

| Coverage Range | Quality Level | Project Status |
|----------------|---------------|----------------|
| 90%+ | Excellent ✅ | Domain Service, Value Objects, Application |
| 80-89% | Good ✅ | Domain Errors, Entity |
| <80% | Needs Improvement ⚠️ | None |

**All business-critical layers exceed 80% coverage threshold!**

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Handlers   │  │ Middleware   │  │   HTTP/JSON  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   UseCases   │  │   DTOs/Req   │  │  Business     │     │
│  │              │  │              │  │  Logic        │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Entities   │  │ Value Objects│  │ Domain Errors│     │
│  │              │  │              │  │              │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │   Repos      │  │Domain Services│                        │
│  │  (Interfaces)│  │               │                        │
│  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Database   │  │     JWT      │  │   Repository │     │
│  │              │  │              │  │  Impls       │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Project Structure

```
auth/
├── cmd/                          # Application entry point
│   ├── main.go                   # Main application file
│   └── .env                      # Environment variables
├── delivery/                     # External systems interface
│   ├── handler/                  # HTTP handlers
│   │   └── auth_handler.go       # Authentication HTTP handlers
│   └── middleware/               # Gin middlewares
│       └── auth_middleware.go    # JWT authentication middleware
├── application/                  # Application business logic
│   └── usecase/                  # Use cases (application services)
│       ├── auth_usecase.go       # Authentication use case
│       └── auth_usecase_test.go  # Use case tests (100% coverage)
├── domain/                       # Core business logic
│   ├── entity/                   # Domain entities
│   │   ├── user.go               # User entity
│   │   └── user_test.go          # Entity tests
│   ├── valueobject/              # Value objects
│   │   ├── user_vo.go            # Email, Password, HashedPassword VOs
│   │   └── user_vo_test.go       # VO tests
│   ├── repository/               # Repository interfaces
│   │   ├── user_repository.go
│   │   ├── refresh_token_repository.go
│   │   └── mocks/                # Generated mocks
│   ├── service/                  # Domain services
│   │   ├── user_registration_service.go
│   │   ├── token_service.go
│   │   ├── user_registration_service_test.go  # Service tests (94.1% coverage)
│   │   └── mocks/                # Generated mocks
│   └── errors/                   # Domain errors
│       ├── errors.go             # Domain error definitions
│       └── errors_test.go
├── infrastructure/               # External concerns
│   ├── database/                 # Database setup
│   │   ├── connection.go         # Database connection
│   │   └── model/                 # Database models
│   ├── repository/               # Repository implementations
│   │   ├── user_repository_impl.go
│   │   └── refresh_token_repository_impl.go
│   └── jwt/                      # JWT token service
│       └── jwt_token_service.go
└── pkg/                          # Shared utilities
    ├── config/                   # Configuration
    │   └── config.go
    ├── error/                    # Error handling
    │   └── error_handler.go
    ├── response/                 # HTTP responses
    │   └── response.go
    └── utils/                    # Utilities
        └── validator.go

```

## ✨ Features

- **User Registration** - Secure user registration with email validation
- **User Login** - JWT-based authentication with access & refresh tokens
- **User Logout** - Secure logout with token invalidation
- **Password Security** - Bcrypt hashing with strength validation
- **Email Validation** - Comprehensive email format validation
- **JWT Authentication** - Secure token-based authentication
- **Protected Routes** - Middleware for protecting routes
- **CORS Support** - Cross-origin resource sharing configured
- **Comprehensive Testing** - High test coverage (domain: 94.1%, usecase: 100%)

## 🚀 Getting Started

### Prerequisites

- Go 1.26.4 or higher
- MySQL 5.7+ or MySQL 8.0+
- Git

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd auth
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
Create `cmd/.env` file:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=auth_db
DB_USER=root
DB_PASSWORD=your_password

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug  # debug, release, test
APP_ENV=development

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ACCESS_TOKEN_EXPIRY=3600    # 1 hour in seconds
REFRESH_TOKEN_EXPIRY=604800 # 7 days in seconds
```

4. **Create database tables**
```sql
CREATE DATABASE auth_db;
USE auth_db;

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
);

CREATE TABLE refresh_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    token VARCHAR(500) NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_token (token(255)),
    INDEX idx_user_id (user_id)
);
```

5. **Run the application**
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "auth-service",
  "version": "1.0.0"
}
```

### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123",
  "full_name": "John Doe"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Response (400/409):**
```json
{
  "status": "error",
  "message": "EMAIL_ALREADY_EXISTS: email already exists",
  "code": "EMAIL_ALREADY_EXISTS"
}
```

### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user": {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "full_name": "John Doe",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "expires_in": 3600
    }
  }
}
```

**Error Response (401):**
```json
{
  "status": "error",
  "message": "INVALID_CREDENTIALS: invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

### Logout (Protected)
```http
POST /api/auth/logout
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Logout successful"
}
```

### Protected Route Example
```http
GET /api/protected/profile
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "message": "This is a protected route",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com"
}
```

## 🔒 Password Requirements

Passwords must meet the following security requirements:
- **Minimum length**: 8 characters
- **Uppercase letter**: At least one (A-Z)
- **Lowercase letter**: At least one (a-z)
- **Number**: At least one (0-9)

Example valid passwords:
- `Password123`
- `SecurePass456`
- `MyPassword99`

## 🧪 Testing

The project has **comprehensive test coverage** following TDD methodology with unit tests for all critical business logic layers.

### 📊 Test Coverage Report

| Layer | Package | Coverage | Status |
|-------|---------|----------|--------|
| **Application** | `application/usecase` | **100.0%** | ✅ Excellent |
| **Domain** | `domain/service` | **94.1%** | ✅ Excellent |
| **Domain** | `domain/valueobject` | **94.1%** | ✅ Excellent |
| **Domain** | `domain/errors` | **83.3%** | ✅ Good |
| **Domain** | `domain/entity` | **80.0%** | ✅ Good |

**Overall Domain & Application Coverage: 90.3%** ✅

> **Note**: Infrastructure and delivery layers (handlers, middleware, repositories) are excluded from unit testing as they are thin wrappers around frameworks and external dependencies. These are tested through integration tests.

### Test Statistics

- **Total Test Cases**: 120+ tests
- **Test Execution Time**: ~30 seconds
- **Test Framework**: testify + mock
- **Mock Generation**: mockery

### Run all tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run tests for specific layer
```bash
# Domain layer tests
go test -v ./domain/...

# Application layer tests
go test -v ./application/...

# Specific package tests
go test -v ./domain/service/...
go test -v ./domain/valueobject/...
go test -v ./application/usecase/...
```

### View coverage report
```bash
# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage in terminal
go tool cover -func=coverage.out
```

### Test-Driven Development Workflow

This project follows TDD workflow:

1. **Red** 🔴 - Write a failing test for new functionality
2. **Green** 🟢 - Write minimum code to make test pass
3. **Refactor** 🔵 - Improve code while tests stay green

**Example**:
```bash
# TDD cycle for new feature
1. Write test:      vim domain/service/user_service_test.go
2. Run test (fail): go test ./domain/service/...
3. Write code:      vim domain/service/user_service.go
4. Run test (pass): go test ./domain/service/...
5. Refactor:        vim domain/service/user_service.go
6. Run test (pass): go test ./domain/service/...
```

### Generate Mocks
The project uses [mockery](https://github.com/vektra/mockery) for generating mocks:

```bash
# Install mockery
go install github.com/vektra/mockery/v2@latest

# Generate all mocks
mockery --all --dir=./domain/repository --output=./domain/repository/mocks
mockery --all --dir=./domain/service --output=./domain/service/mocks
```

## 🔧 Configuration

Configuration is managed through environment variables. The service loads configuration from:

1. `cmd/.env` file (development)
2. System environment variables (production)

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_NAME` | Database name | `auth_db` |
| `DB_USER` | Database user | `root` |
| `DB_PASSWORD` | Database password | `secret` |
| `JWT_SECRET` | JWT secret key | `your-secret-key` |

### Optional Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode | `debug` |
| `APP_ENV` | Application environment | `development` |
| `ACCESS_TOKEN_EXPIRY` | Access token expiry (seconds) | `3600` |
| `REFRESH_TOKEN_EXPIRY` | Refresh token expiry (seconds) | `604800` |

## 🏛️ Design Patterns & Principles

### Clean Architecture
- **Dependency Inversion**: Dependencies point inward toward the domain
- **Separation of Concerns**: Each layer has a specific responsibility
- **Testability**: Business logic is independent of frameworks

### Domain-Driven Design (DDD)
- **Entities**: Core domain objects with identity (`User`)
- **Value Objects**: Immutable objects without identity (`Email`, `Password`, `HashedPassword`)
- **Domain Services**: Business logic that doesn't fit in entities
- **Repositories**: Data access abstractions
- **Ubiquitous Language**: Consistent terminology across all layers

### SOLID Principles
- **S**RP: Single Responsibility - Each component has one reason to change
- **O**CP: Open/Closed - Open for extension, closed for modification
- **L**SP: Liskov Substitution - Subtypes are substitutable for base types
- **I**SP: Interface Segregation - Small, focused interfaces
- **D**IP: Dependency Inversion - Depend on abstractions, not concretions

## 🔐 Security Features

1. **Password Hashing** - Bcrypt with automatic salt generation
2. **JWT Tokens** - Secure token-based authentication
3. **Token Expiration** - Configurable access & refresh token expiry
4. **Password Strength Validation** - Enforces strong password requirements
5. **Email Validation** - RFC-compliant email format validation
6. **CORS Protection** - Configurable CORS middleware
7. **Panic Recovery** - Graceful error handling

## 📦 Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `gin-gonic/gin` | v1.12.0 | HTTP web framework |
| `go-sql-driver/mysql` | v1.10.0 | MySQL driver |
| `golang-jwt/jwt` | v5.3.1 | JWT implementation |
| `google/uuid` | v1.6.0 | UUID generation |
| `joho/godotenv` | v1.5.1 | Environment configuration |
| `stretchr/testify` | v1.11.1 | Testing framework |
| `golang.org/x/crypto` | v0.54.0 | Bcrypt password hashing |

## 🚀 Deployment

### Production Build
```bash
# Build for current platform
go build -o auth-service cmd/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o auth-service-linux cmd/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o auth-service.exe cmd/main.go
```

### Docker Deployment
Create a `Dockerfile`:
```dockerfile
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-service .
COPY --from=builder /app/cmd/.env .env
EXPOSE 8080
CMD ["./auth-service"]
```

Build and run:
```bash
docker build -t auth-service .
docker run -p 8080:8080 --env-file cmd/.env auth-service
```

## 📝 Error Codes

| Code | Description |
|------|-------------|
| `EMAIL_EMPTY` | Email cannot be empty |
| `INVALID_EMAIL_FORMAT` | Invalid email format |
| `EMAIL_ALREADY_EXISTS` | Email already exists |
| `PASSWORD_EMPTY` | Password cannot be empty |
| `PASSWORD_TOO_SHORT` | Password must be at least 8 characters |
| `PASSWORD_MISSING_UPPERCASE` | Password must contain uppercase letter |
| `PASSWORD_MISSING_LOWERCASE` | Password must contain lowercase letter |
| `PASSWORD_MISSING_NUMBER` | Password must contain number |
| `INVALID_CREDENTIALS` | Invalid email or password |
| `USER_NOT_FOUND` | User not found |
| `TOKEN_GENERATION_FAILED` | Failed to generate token |
| `REFRESH_TOKEN_SAVE_FAILED` | Failed to save refresh token |
| `REFRESH_TOKEN_DELETE_FAILED` | Failed to delete refresh tokens |
| `INVALID_USER_ID` | Invalid user ID format |
| `INVALID_REFRESH_TOKEN` | Invalid refresh token |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Author

**Hanapi**

## 🙏 Acknowledgments

- Clean Architecture principles by Robert C. Martin
- Domain-Driven Design by Eric Evans
- Gin web framework
- JWT authentication best practices
