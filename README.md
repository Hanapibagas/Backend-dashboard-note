# Backend Dashboard Note - Auth Service

## 🔑 Authentication Service

Production-ready authentication microservice built with **Go** using **Clean Architecture** and **Domain-Driven Design (DDD)** principles.

### ✨ Key Highlights

- **Clean Architecture** - Clear separation of concerns with 4 distinct layers
- **Domain-Driven Design** - Rich domain model with entities, value objects, and domain services
- **Test-Driven Development** - 90.3%+ test coverage with 120+ test cases
- **JWT Authentication** - Secure token-based auth with access & refresh tokens
- **High Test Coverage** - Domain (94.1%), Application (100%)

### 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│                    (HTTP Handlers, Middleware)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                          │
│                   (Use Cases, DTOs)                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                             │
│              (Entities, Value Objects, Services)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                         │
│              (Database, JWT, Repository Impls)               │
└─────────────────────────────────────────────────────────────┘
```

### 📊 Test Coverage Quality

| Layer | Coverage | Quality |
|-------|----------|---------|
| Application (Use Cases) | 100.0% | ✅ Excellent |
| Domain (Services) | 94.1% | ✅ Excellent |
| Domain (Value Objects) | 94.1% | ✅ Excellent |
| Domain (Errors) | 83.3% | ✅ Good |
| Domain (Entities) | 80.0% | ✅ Good |

**Overall: 90.3% coverage** across all business-critical layers

### 🚀 Core Features

- **User Registration** - Secure signup with email & password validation
- **User Login** - JWT-based authentication with dual-token system
- **User Logout** - Token invalidation for secure logout
- **Password Security** - Bcrypt hashing with strength requirements
- **Protected Routes** - Middleware for route protection
- **Comprehensive Testing** - TDD approach with high coverage

### 📡 Quick Start

```bash
# Navigate to auth service
cd auth

# Install dependencies
go mod download

# Configure environment
cp cmd/.env.example cmd/.env
# Edit cmd/.env with your database & JWT settings

# Run database migrations
mysql -u root -p < database/schema.sql

# Start the service
go run cmd/main.go
```

Service runs on `http://localhost:8080`

### 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/auth/register` | User registration |
| POST | `/api/auth/login` | User login |
| POST | `/api/auth/logout` | User logout (protected) |
| GET | `/api/protected/profile` | Protected route example |

### 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific layer tests
go test -v ./domain/...
go test -v ./application/...
```

### 🔒 Security Features

- ✅ Bcrypt password hashing with automatic salt
- ✅ JWT token-based authentication
- ✅ Configurable token expiration
- ✅ Password strength validation (8+ chars, uppercase, lowercase, number)
- ✅ Email format validation (RFC-compliant)
- ✅ CORS protection
- ✅ Panic recovery middleware

### 📦 Tech Stack

| Package | Purpose |
|---------|---------|
| gin-gonic/gin | HTTP web framework |
| golang-jwt/jwt | JWT implementation |
| google/uuid | UUID generation |
| golang.org/x/crypto | Bcrypt password hashing |
| stretchr/testify | Testing framework |
| vektra/mockery | Mock generation |

### 🏛️ Design Patterns

- **Clean Architecture** - Dependency inversion & separation of concerns
- **Domain-Driven Design** - Entities, value objects, domain services
- **SOLID Principles** - Single responsibility, open/closed, Liskov substitution, interface segregation, dependency inversion
- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Loose coupling & testability

### 📁 Project Structure

```
auth/
├── cmd/                    # Application entry point
├── delivery/               # HTTP handlers & middleware
├── application/            # Use cases & business logic
├── domain/                 # Core domain (entities, VOs, services)
├── infrastructure/         # Database, JWT, repository impls
└── pkg/                    # Shared utilities
```

### 🔑 Password Requirements

- Minimum 8 characters
- At least 1 uppercase letter (A-Z)
- At least 1 lowercase letter (a-z)
- At least 1 number (0-9)

Valid examples: `Password123`, `SecurePass456`

### 📝 Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DB_HOST | ✅ | - | Database host |
| DB_PORT | ✅ | - | Database port |
| DB_NAME | ✅ | - | Database name |
| DB_USER | ✅ | - | Database user |
| DB_PASSWORD | ✅ | - | Database password |
| JWT_SECRET | ✅ | - | JWT secret key |
| SERVER_PORT | ❌ | 8080 | Server port |
| GIN_MODE | ❌ | debug | Gin mode (debug/release/test) |
| ACCESS_TOKEN_EXPIRY | ❌ | 3600 | Access token expiry (seconds) |
| REFRESH_TOKEN_EXPIRY | ❌ | 604800 | Refresh token expiry (seconds) |

### 🚢 Deployment

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o auth-service cmd/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o auth-service.exe cmd/main.go

# Docker build
docker build -t auth-service .
docker run -p 8080:8080 --env-file cmd/.env auth-service
```

### 🔄 CI/CD Pipeline

This project uses **GitHub Actions** with **reusable workflows** for continuous integration and continuous deployment.

#### Pipeline Architecture

```
.github/workflows/
├── templates/                    # Reusable workflow templates
│   ├── ci-template.yml          # CI/quality checks template
│   └── docker-template.yml      # Docker build & push template
├── services/                     # Service-specific workflows
│   └── auth/                    # Auth service workflows
│       ├── ci.yml              # Calls ci-template
│       └── docker.yml          # Calls docker-template
└── _examples/                   # Templates for new services
    ├── ci.yml.example
    └── docker.yml.example
```

**CI Template** - Runs for each service:
- ✅ Automated testing on every push & PR
- ✅ Code quality checks (go vet, go fmt, go mod tidy)
- ✅ Build verification
- ✅ Security scanning with Gosec
- ✅ Coverage report generation

**Docker Template** - Runs for each service:
- ✅ Multi-architecture Docker builds (linux/amd64, linux/arm64)
- ✅ Automated pushes to Docker Hub
- ✅ Container security scanning with Trivy
- ✅ Layer caching for faster builds

**Reusable Design** - Templates sekali, gunakan untuk semua services:
- Setiap service punya workflow sendiri yang call template
- Independent pipeline execution (service A tidak affect service B)
- Mudah tambah service baru

#### Setup Instructions

**1. Docker Hub Secrets**

Add these secrets to your GitHub repository (`Settings` → `Secrets and variables` → `Actions`):

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `DOCKER_HUB_USERNAME` | Your Docker Hub username | `hanapi` |
| `DOCKER_HUB_TOKEN` | Docker Hub access token | `dckr_xxxxx...` |

To create a Docker Hub token:
- Go to [Docker Hub Settings](https://hub.docker.com/settings/security)
- Click "New Access Token"
- Give it a descriptive name (e.g., `github-actions`)
- Select "Read & Write" permissions
- Copy the token and add it to GitHub secrets

**2. Workflow Triggers**

The CI/CD pipelines run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Manual trigger via GitHub Actions UI

**3. Adding New Services**

For complete instructions on adding a new service, see [**ADDING_NEW_SERVICE.md**](ADDING_NEW_SERVICE.md).

Quick summary for adding a new service (e.g., `product`):

```bash
# 1. Create service workflow directory
mkdir -p .github/workflows/services/product

# 2. Copy and customize example workflows
cp .github/workflows/_examples/ci.yml.example \
   .github/workflows/services/product/ci.yml
cp .github/workflows/_examples/docker.yml.example \
   .github/workflows/services/product/docker.yml

# 3. Edit the files and replace placeholders:
#    - <SERVICE_NAME> → product
#    - <SERVICE_PATH> → product
#    - <DOCKER_IMAGE_NAME> → hanapi/product-service
```

See [`.github/workflows/README.md`](.github/workflows/README.md) for workflow architecture details.

#### Docker Image Naming

Images are tagged as:
- `hanapi/auth-service:latest` (main branch only)
- `hanapi/auth-service:develop` (develop branch)
- `hanapi/auth-service:main-<sha>` (specific commits)
- `hanapi/auth-service:pr-<number>` (pull requests)

#### Pipeline Status

Check pipeline status in GitHub Actions tab:
- 📊 Test results and coverage reports
- 🔍 Security scan findings
- 🐳 Docker image build status
- 📦 Available image tags

#### Local Testing

```bash
# Test CI pipeline locally
cd auth
go test -v -race ./...
go vet ./...
go fmt ./...

# Build Docker image locally
docker build -t hanapi/auth-service:latest .
docker run -p 8080:8080 --env-file cmd/.env hanapi/auth-service:latest
```

### 📖 Documentation

- [Auth Service Documentation](auth/README.md) - Complete auth service guide
- [Adding New Service Guide](ADDING_NEW_SERVICE.md) - How to add new microservices
- [GitHub Workflows Documentation](.github/workflows/README.md) - CI/CD workflows explained

---

**Built with ❤️ using Clean Architecture & DDD principles**
