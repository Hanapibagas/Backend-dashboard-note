# GitHub Actions Workflows Structure

This directory contains GitHub Actions workflows for CI/CD pipelines across all services in the monorepo.

## Directory Structure

```
.github/
└── workflows/
    ├── templates/              # Reusable workflow templates
    │   ├── ci-template.yml     # CI template (test, lint, build, security)
    │   └── docker-template.yml # Docker build & push template
    ├── services/               # Service-specific workflows
    │   └── auth/               # Auth service workflows
    │       ├── ci.yml
    │       └── docker.yml
    └── _examples/              # Example workflows for new services
        ├── ci.yml.example
        └── docker.yml.example
```

## How Templates Work

Reusable workflows allow you to define CI/CD logic once and reuse it across multiple services:

- **CI Template** (`ci-template.yml`): Runs tests, linting, build verification, and security scans
- **Docker Template** (`docker-template.yml`): Builds and pushes multi-arch Docker images with security scanning

Each service calls these templates with their specific parameters.

## Adding a New Service

When adding a new service (e.g., `product`, `user`, `order`), follow these steps:

### 1. Create Service Workflow Directory

```bash
mkdir -p .github/workflows/services/<service-name>
```

### 2. Create CI Workflow

Copy and customize the example:

```bash
cp .github/workflows/_examples/ci.yml.example \
   .github/workflows/services/<service-name>/ci.yml
```

Edit the file and replace placeholders:
- `<SERVICE_NAME>`: Service name (e.g., `product`)
- `<SERVICE_PATH>`: Path to service directory (e.g., `product`)
- Optionally adjust `go_version` if needed

### 3. Create Docker Workflow

```bash
cp .github/workflows/_examples/docker.yml.example \
   .github/workflows/services/<service-name>/docker.yml
```

Edit the file and replace placeholders:
- `<SERVICE_NAME>`: Service name (e.g., `product`)
- `<SERVICE_PATH>`: Path to service directory (e.g., `product`)
- `<DOCKER_IMAGE_NAME>`: Docker Hub image name (e.g., `hanapi/product-service`)

### 4. Commit and Push

```bash
git add .github/workflows/services/<service-name>/
git commit -m "feat: add CI/CD workflows for <service-name> service"
git push origin main
```

## Example: Adding a Product Service

Here's a complete example for adding a `product` service:

### Directory Structure

```
product/
├── cmd/
│   └── main.go
├── domain/
├── application/
├── infrastructure/
└── go.mod
```

### CI Workflow: `.github/workflows/services/product/ci.yml`

```yaml
name: CI - Product Service

on:
  push:
    branches: [main, develop]
    paths: ['product/**']
  pull_request:
    branches: [main, develop]
    paths: ['product/**']

jobs:
  call-ci-workflow:
    uses: ./.github/workflows/templates/ci-template.yml
    with:
      service_name: product
      service_path: product
      go_version: '1.23'
```

### Docker Workflow: `.github/workflows/services/product/docker.yml`

```yaml
name: Docker - Product Service

on:
  push:
    branches: [main, develop]
    paths: ['product/**']
  pull_request:
    branches: [main, develop]
    paths: ['product/**']
  workflow_dispatch:

jobs:
  call-docker-workflow:
    uses: ./.github/workflows/templates/docker-template.yml
    secrets:
      docker_hub_username: ${{ secrets.DOCKER_HUB_USERNAME }}
      docker_hub_token: ${{ secrets.DOCKER_HUB_TOKEN }}
    with:
      service_name: product
      service_path: product
      docker_image: hanapi/product-service
```

## Workflow Parameters

### CI Template Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `service_name` | ✅ | - | Service name (used for artifacts) |
| `service_path` | ✅ | - | Path to service directory |
| `go_version` | ❌ | `1.23` | Go version to use |

### Docker Template Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `service_name` | ✅ | - | Service name |
| `service_path` | ✅ | - | Path to service directory |
| `docker_image` | ✅ | - | Docker Hub image name |
| `docker_platforms` | ❌ | `linux/amd64,linux/arm64` | Target platforms |

## Required GitHub Secrets

Both templates require these repository secrets:

- `DOCKER_HUB_USERNAME`: Your Docker Hub username
- `DOCKER_HUB_TOKEN`: Docker Hub access token with Read & Write permissions

Set these in: `Settings` → `Secrets and variables` → `Actions`

## Workflow Features

### CI Template Includes:
- ✅ Automated testing with race detection
- ✅ Code quality checks (go vet, go fmt, go mod tidy)
- ✅ Build verification
- ✅ Security scanning with Gosec
- ✅ Coverage report generation

### Docker Template Includes:
- ✅ Multi-architecture builds (AMD64 & ARM64)
- ✅ Automated pushes to Docker Hub
- ✅ Container security scanning with Trivy
- ✅ Layer caching for faster builds
- ✅ Smart image tagging

## Troubleshooting

### Workflow Not Triggering

Check that:
1. File path filters match your service directory
2. Workflow files are in the correct location
3. File names are exactly `ci.yml` and `docker.yml`

### Docker Push Fails

Check that:
1. GitHub secrets are correctly set
2. Docker Hub token has Read & Write permissions
3. Docker image name is correct

### Build Fails

Check that:
1. Service has `cmd/main.go` entry point
2. `go.mod` and `go.sum` files exist
3. All dependencies are available

## Current Services

| Service | Directory | Docker Image | Status |
|---------|-----------|--------------|--------|
| Auth | `auth/` | `hanapi/auth-service` | ✅ Active |
| (Add more here as you create them) | | | |
