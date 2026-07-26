# Adding a New Service Guide

This guide explains how to add a new microservice to this monorepo with complete CI/CD setup.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Step-by-Step Guide](#step-by-step-guide)
  - [Step 1: Create Service Directory](#step-1-create-service-directory)
  - [Step 2: Setup Service Code](#step-2-setup-service-code)
  - [Step 3: Create Docker Configuration](#step-3-create-docker-configuration)
  - [Step 4: Setup CI/CD Workflows](#step-4-setup-cicd-workflows)
  - [Step 5: Test Locally](#step-5-test-locally)
  - [Step 6: Commit and Push](#step-6-commit-and-push)
- [Complete Example: Product Service](#complete-example-product-service)
- [Troubleshooting](#troubleshooting)

---

## Overview

Adding a new service involves:

1. **Creating service directory** with proper structure
2. **Setting up Go modules** and dependencies
3. **Adding Dockerfile** for containerization
4. **Configuring CI/CD workflows** using reusable templates
5. **Testing** everything works

The CI/CD setup uses **reusable GitHub Actions workflows**, so you only need to configure service-specific parameters.

---

## Prerequisites

Before adding a new service, ensure:

- ✅ You have write access to the repository
- ✅ GitHub secrets are configured (`DOCKER_HUB_USERNAME`, `DOCKER_HUB_TOKEN`)
- ✅ You understand the service's purpose and requirements
- ✅ You've decided on the service name (e.g., `product`, `user`, `order`)

---

## Step-by-Step Guide

### Step 1: Create Service Directory

Create a new directory for your service at the root of the monorepo:

```bash
mkdir -p <service-name>
cd <service-name>
```

Example:
```bash
mkdir -p product
cd product
```

---

### Step 2: Setup Service Code

#### 2.1 Initialize Go Module

```bash
go mod init github.com/hanapi/<service-name>
```

Example:
```bash
go mod init github.com/hanapi/product
```

#### 2.2 Create Service Structure

Follow the same Clean Architecture pattern as the auth service:

```
<service-name>/
├── cmd/
│   └── main.go              # Application entry point
├── domain/                  # Core business logic
│   ├── entity/             # Domain entities
│   ├── valueobject/        # Value objects
│   ├── repository/         # Repository interfaces
│   └── service/            # Domain services
├── application/            # Use cases
│   └── usecase/
├── infrastructure/         # External dependencies
│   ├── database/
│   ├── repository/
│   └── config/
├── delivery/              # HTTP handlers
│   ├── handler/
│   └── middleware/
├── pkg/                   # Shared utilities
│   ├── config/
│   ├── error/
│   └── response/
├── database/              # Database schemas
│   └── schema.sql
├── Dockerfile
├── .dockerignore
├── go.mod
├── go.sum
└── README.md
```

#### 2.3 Create Main Entry Point (`cmd/main.go`)

```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize configuration
    // Initialize database
    // Setup routes
    // Start server

    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "service": "<service-name>",
        })
    })

    log.Println("Starting <service-name> service on :8080")
    r.Run(":8080")
}
```

#### 2.4 Create Basic Dockerfile

Create `Dockerfile` in your service directory:

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy database schema if exists
COPY --from=builder /app/database ./database

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

#### 2.5 Create .dockerignore

Create `.dockerignore` in your service directory:

```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
main
<service-name>

# Test files
*_test.go
*.test
*.out
coverage.html
coverage.out

# Git
.git
.gitignore

# Documentation
README.md
*.md

# IDE
.vscode
.idea
*.swp
*.swo
*~

# Environment
.env
.env.local
.env.*.local

# CI/CD
.github
.gitlab-ci.yml

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Misc
.mockery.yml
```

---

### Step 3: Create Docker Configuration

#### 3.1 Copy Dockerfile from Auth (Optional)

If your service has similar structure to auth, you can copy and modify:

```bash
cp ../auth/Dockerfile ./Dockerfile
cp ../auth/.dockerignore ./.dockerignore
```

#### 3.2 Customize for Your Service

Edit the Dockerfile if your service has specific requirements.

---

### Step 4: Setup CI/CD Workflows

#### 4.1 Naming Convention (no subdirectory needed)

> ⚠️ GitHub Actions **only detects workflows placed DIRECTLY in `.github/workflows/`**.
> Do NOT create subdirectories like `.github/workflows/services/<service-name>/` —
> workflows there will be silently ignored.
> Ref: [github/community#18055](https://github.com/orgs/community/discussions/18055)

Workflow files use a **flat naming convention**:
- `<service-name>-ci.yml` — CI workflow
- `<service-name>-docker.yml` — Docker workflow

#### 4.2 Create CI Workflow

Copy and customize the example:

```bash
cp .github/workflows/_examples/ci.yml.example \
   .github/workflows/<service-name>-ci.yml
```

Edit the file and replace placeholders:

```yaml
name: CI - <SERVICE_NAME>

on:
  push:
    branches:
      - main
      - develop
    paths:
      - '<SERVICE_PATH>/**'
      - '.github/workflows/<SERVICE_NAME>-ci.yml'
  pull_request:
    branches:
      - main
      - develop
    paths:
      - '<SERVICE_PATH>/**'
      - '.github/workflows/<SERVICE_NAME>-ci.yml'

jobs:
  call-ci-workflow:
    uses: ./.github/workflows/ci-template.yml
    with:
      service_name: <SERVICE_NAME>
      service_path: <SERVICE_PATH>
      go_version: '1.23'
```

**Replace:**
- `<SERVICE_NAME>`: Your service name (e.g., `Product Service`)
- `<SERVICE_PATH>`: Path to service (e.g., `product`)

#### 4.3 Create Docker Workflow

```bash
cp .github/workflows/_examples/docker.yml.example \
   .github/workflows/<service-name>-docker.yml
```

Edit the file and replace placeholders:

```yaml
name: Docker - <SERVICE_NAME>

on:
  push:
    branches:
      - main
      - develop
    paths:
      - '<SERVICE_PATH>/**'
      - '.github/workflows/<SERVICE_NAME>-docker.yml'
  pull_request:
    branches:
      - main
      - develop
    paths:
      - '<SERVICE_PATH>/**'
      - '.github/workflows/<SERVICE_NAME>-docker.yml'
  workflow_dispatch:

jobs:
  call-docker-workflow:
    uses: ./.github/workflows/docker-template.yml
    secrets:
      docker_hub_username: ${{ secrets.DOCKER_HUB_USERNAME }}
      docker_hub_token: ${{ secrets.DOCKER_HUB_TOKEN }}
    with:
      service_name: <SERVICE_NAME>
      service_path: <SERVICE_PATH>
      docker_image: <DOCKER_IMAGE_NAME>
      docker_platforms: linux/amd64,linux/arm64
```

**Replace:**
- `<SERVICE_NAME>`: Your service name (e.g., `Product Service`)
- `<SERVICE_PATH>`: Path to service (e.g., `product`)
- `<DOCKER_IMAGE_NAME>`: Docker Hub image name (e.g., `hanapi/product-service`)

---

### Step 5: Test Locally

#### 5.1 Test Go Application

```bash
cd <service-name>

# Run tests
go test ./...

# Run application
go run cmd/main.go

# Build binary
go build -o <service-name> cmd/main.go
```

#### 5.2 Test Docker Build

```bash
# Build Docker image
docker build -t <docker-image-name>:latest .

# Run container
docker run -p 8080:8080 <docker-image-name>:latest

# Test health endpoint
curl http://localhost:8080/health
```

#### 5.3 Test CI Checks Locally

```bash
cd <service-name>

# Format check
gofmt -l .

# Vet check
go vet ./...

# Mod tidy check
go mod tidy
git diff go.mod go.sum  # Should be empty
```

---

### Step 6: Commit and Push

#### 6.1 Review Changes

```bash
# Check all files
git status

# Review workflow files
cat .github/workflows/<service-name>-ci.yml
cat .github/workflows/<service-name>-docker.yml
```

#### 6.2 Commit Changes

```bash
git add .
git commit -m "feat: add <service-name> service

- Add service structure with Clean Architecture
- Setup CI/CD workflows using reusable templates
- Add Dockerfile and docker configuration
- Configure Docker Hub image builds

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

#### 6.3 Push to Remote

```bash
# Push to main branch
git push origin main

# Or push to feature branch first
git push -u origin feature/add-<service-name>-service
```

#### 6.4 Verify GitHub Actions

1. Go to GitHub repository
2. Navigate to **Actions** tab
3. Check workflow runs for your new service:
   - ✅ CI workflow should run tests and quality checks
   - ✅ Docker workflow should build and push image

---

## Complete Example: Product Service

Here's a complete example for adding a `product` service.

### Directory Structure

```
product/
├── cmd/
│   └── main.go
├── domain/
│   └── entity/
│       └── product.go
├── go.mod
├── go.sum
├── Dockerfile
├── .dockerignore
└── README.md

.github/
└── workflows/
    ├── product-ci.yml        # Product CI workflow
    └── product-docker.yml    # Product Docker workflow
```

### CI Workflow: `.github/workflows/product-ci.yml`

```yaml
name: CI - Product Service

on:
  push:
    branches:
      - main
      - develop
    paths:
      - 'product/**'
      - '.github/workflows/product-ci.yml'
  pull_request:
    branches:
      - main
      - develop
    paths:
      - 'product/**'
      - '.github/workflows/product-ci.yml'

jobs:
  call-ci-workflow:
    uses: ./.github/workflows/ci-template.yml
    with:
      service_name: product
      service_path: product
      go_version: '1.23'
```

### Docker Workflow: `.github/workflows/product-docker.yml`

```yaml
name: Docker - Product Service

on:
  push:
    branches:
      - main
      - develop
    paths:
      - 'product/**'
      - '.github/workflows/product-docker.yml'
  pull_request:
    branches:
      - main
      - develop
    paths:
      - 'product/**'
      - '.github/workflows/product-docker.yml'
  workflow_dispatch:

jobs:
  call-docker-workflow:
    uses: ./.github/workflows/docker-template.yml
    secrets:
      docker_hub_username: ${{ secrets.DOCKER_HUB_USERNAME }}
      docker_hub_token: ${{ secrets.DOCKER_HUB_TOKEN }}
    with:
      service_name: product
      service_path: product
      docker_image: hanapi/product-service
      docker_platforms: linux/amd64,linux/arm64
```

### Expected Results

After pushing to GitHub:

✅ **CI Workflow** will:
- Run all tests
- Check code quality (go vet, go fmt, go mod tidy)
- Build the application binary
- Run security scans with Gosec

✅ **Docker Workflow** will:
- Build multi-architecture Docker images (AMD64 & ARM64)
- Push images to Docker Hub: `hanapi/product-service`
- Run container security scanning with Trivy

✅ **Docker Images Available**:
- `hanapi/product-service:latest` (main branch)
- `hanapi/product-service:develop` (develop branch)
- `hanapi/product-service:main-<sha>` (specific commits)
- `hanapi/product-service:pr-<number>` (pull requests)

---

## Troubleshooting

### Workflow Not Triggering

**Problem**: GitHub Actions workflow doesn't run after pushing.

**Solutions**:
1. Check file paths in workflow triggers match your directory structure
2. Ensure workflow files are placed **directly** in `.github/workflows/` (NOT in subdirectories — GitHub Actions does not scan subfolders)
3. Verify workflow files are named `<service>-ci.yml` and `<service>-docker.yml`
4. Check workflow syntax is correct (GitHub will show syntax errors)

### Docker Build Fails

**Problem**: Docker build fails in workflow.

**Solutions**:
1. Verify Dockerfile exists in service directory
2. Check Dockerfile syntax locally: `docker build -t test .`
3. Ensure `cmd/main.go` exists and is the entry point
4. Check Go module files exist: `go.mod`, `go.sum`

### Tests Fail in CI

**Problem**: Tests pass locally but fail in CI.

**Solutions**:
1. Check Go version matches (default: 1.23)
2. Verify all dependencies are in `go.mod`
3. Run tests with race detector locally: `go test -race ./...`
4. Check environment-specific code

### Docker Push Fails

**Problem**: Docker image doesn't push to Docker Hub.

**Solutions**:
1. Verify GitHub secrets are set:
   - `DOCKER_HUB_USERNAME`
   - `DOCKER_HUB_TOKEN`
2. Check Docker Hub token has "Read & Write" permissions
3. Ensure Docker Hub image name is correct
4. Check if you're authenticated locally: `docker login`

### Build Binary Not Found

**Problem**: Build job fails with "binary not found" error.

**Solutions**:
1. Ensure `cmd/main.go` exists
2. Check build command in `cmd/main.go`
3. Verify service can build locally: `go build cmd/main.go`
4. Check for any platform-specific code

### Path Filter Not Working

**Problem**: Workflow doesn't trigger when pushing to service directory.

**Solutions**:
1. Check path filter uses correct syntax: `'service/**'`
2. Ensure path matches exactly (case-sensitive)
3. Test with: `git push` after making changes in service directory
4. Check workflow logs for path matching information

---

## Best Practices

### 1. Service Naming

- Use lowercase: `product`, `user`, `order`
- Avoid special characters or spaces
- Keep names descriptive and concise

### 2. Docker Image Naming

Follow this pattern:
```
<docker-hub-username>/<service-name>-service
```

Examples:
- `hanapi/product-service`
- `hanapi/user-service`
- `hanapi/order-service`

### 3. Service Structure

Maintain consistency with existing services:
- Use Clean Architecture
- Follow same directory structure
- Include comprehensive tests
- Add proper documentation

### 4. Environment Variables

Use a `.env.example` file to document required variables:

```bash
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=product_db
DB_USER=root
DB_PASSWORD=password

# Service-Specific Configuration
```

### 5. Testing

- Aim for high test coverage (80%+)
- Include unit tests for business logic
- Add integration tests where needed
- Mock external dependencies

### 6. Documentation

Always include a `README.md` in your service directory with:
- Service purpose and features
- API endpoints documentation
- Setup instructions
- Environment variables reference

---

## Checklist

Use this checklist before committing a new service:

- [ ] Service directory created at root level
- [ ] Go module initialized (`go.mod` exists)
- [ ] Entry point exists (`cmd/main.go`)
- [ ] Service can build locally (`go build cmd/main.go`)
- [ ] Service runs without errors (`go run cmd/main.go`)
- [ ] Tests exist and pass (`go test ./...`)
- [ ] Dockerfile created and works locally
- [ ] `.dockerignore` created
- [ ] CI workflow created (`.github/workflows/<service>-ci.yml`)
- [ ] Docker workflow created (`.github/workflows/<service>-docker.yml`)
- [ ] All placeholders replaced in workflows
- [ ] Docker image name is correct
- [ ] Documentation updated (README.md)
- [ ] Files committed and pushed

---

## Additional Resources

- [GitHub Reusable Workflows Documentation](https://docs.github.com/en/actions/using-workflows/reusing-workflows)
- [Docker Multi-Platform Builds](https://docs.docker.com/build/building/multi-platform/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

---

**Need help?** Check the [main README](../README.md) or [workflow documentation](.github/workflows/README.md) for more information.
