# GitHub Actions Workflows Structure

This directory contains GitHub Actions workflows for CI/CD pipelines across all services in the monorepo.

> ⚠️ **IMPORTANT**: GitHub Actions **only detects workflow files placed DIRECTLY in `.github/workflows/`**.
> Workflow files in subdirectories (e.g., `.github/workflows/services/auth/ci.yml`) will **NOT** be picked up.
> Always place workflow `.yml` files at the root of this directory.
> Ref: [github/community#18055](https://github.com/orgs/community/discussions/18055)

## Directory Structure

```
.github/
└── workflows/
    ├── auth-ci.yml           # Auth service CI workflow (caller)
    ├── auth-docker.yml       # Auth service Docker workflow (caller)
    ├── ci-template.yml       # Reusable CI template (test, lint, build, security)
    ├── docker-template.yml   # Reusable Docker build & push template
    └── _examples/            # Example files for new services (NOT workflows)
        ├── ci.yml.example
        └── docker.yml.example
```

### Naming Convention

Since all workflows must live at the root of `.github/workflows/`, we use a **flat naming convention** to keep things organized:

| Pattern | Purpose | Example |
|---------|---------|---------|
| `<service>-ci.yml` | CI caller for a service | `auth-ci.yml`, `product-ci.yml` |
| `<service>-docker.yml` | Docker caller for a service | `auth-docker.yml`, `product-docker.yml` |
| `ci-template.yml` | Reusable CI template | (shared) |
| `docker-template.yml` | Reusable Docker template | (shared) |

## How Templates Work

Reusable workflows allow you to define CI/CD logic once and reuse it across multiple services:

- **CI Template** (`ci-template.yml`): Runs tests, linting, build verification, and security scans
- **Docker Template** (`docker-template.yml`): Builds and pushes multi-arch Docker images with security scanning

Each service's caller workflow references the template via `uses:` and passes service-specific parameters.

## Adding a New Service

When adding a new service (e.g., `product`, `user`, `order`), follow these steps:

### 1. Create the CI Workflow

Copy the example to the **root** of `.github/workflows/` with a service-specific name:

```bash
cp .github/workflows/_examples/ci.yml.example \
   .github/workflows/<service-name>-ci.yml
```

Edit the file and replace placeholders:
- `<SERVICE_NAME>`: Service name (e.g., `product`)
- `<SERVICE_PATH>`: Path to service directory (e.g., `product`)
- Optionally adjust `go_version` if needed

### 2. Create the Docker Workflow

```bash
cp .github/workflows/_examples/docker.yml.example \
   .github/workflows/<service-name>-docker.yml
```

Edit the file and replace placeholders:
- `<SERVICE_NAME>`: Service name (e.g., `product`)
- `<SERVICE_PATH>`: Path to service directory (e.g., `product`)
- `<DOCKER_IMAGE_NAME>`: Docker Hub image name (e.g., `hanapi23/product-service`)

### 3. Commit and Push

```bash
git add .github/workflows/<service-name>-ci.yml .github/workflows/<service-name>-docker.yml
git commit -m "feat: add CI/CD workflows for <service-name> service"
git push origin main
```

## Example: Adding a Product Service

### CI Workflow: `.github/workflows/product-ci.yml`

```yaml
name: CI - Product Service

on:
  push:
    branches: [main, develop]
    paths:
      - 'product/**'
      - '.github/workflows/product-ci.yml'
      - '.github/workflows/ci-template.yml'
  pull_request:
    branches: [main, develop]
    paths:
      - 'product/**'
      - '.github/workflows/product-ci.yml'
      - '.github/workflows/ci-template.yml'

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
    branches: [main, develop]
    paths:
      - 'product/**'
      - '.github/workflows/product-docker.yml'
      - '.github/workflows/docker-template.yml'
  pull_request:
    branches: [main, develop]
    paths:
      - 'product/**'
      - '.github/workflows/product-docker.yml'
      - '.github/workflows/docker-template.yml'
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
      docker_image: hanapi23/product-service
```

## Workflow Parameters

### CI Template Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `service_name` | Yes | - | Service name (used for artifacts) |
| `service_path` | Yes | - | Path to service directory |
| `go_version` | No | `1.23` | Go version to use |

### Docker Template Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `service_name` | Yes | - | Service name |
| `service_path` | Yes | - | Path to service directory |
| `docker_image` | Yes | - | Docker Hub image name |
| `docker_platforms` | No | `linux/amd64,linux/arm64` | Target platforms |

## Required GitHub Secrets

Both templates require these repository secrets:

- `DOCKER_HUB_USERNAME`: Your Docker Hub username
- `DOCKER_HUB_TOKEN`: Docker Hub access token with Read & Write permissions

Set these in: `Settings` → `Secrets and variables` → `Actions`

## Workflow Features

### CI Template Includes:
- Automated testing with race detection
- Code quality checks (go vet, go fmt, go mod tidy)
- Build verification
- Security scanning with Gosec
- Coverage report generation

### Docker Template Includes:
- Multi-architecture builds (AMD64 & ARM64)
- Automated pushes to Docker Hub
- Container security scanning with Trivy
- Layer caching for faster builds
- Smart image tagging

## Troubleshooting

### Workflow Not Triggering / Not Showing in Actions Tab

Check that:
1. Workflow `.yml` files are placed **directly** in `.github/workflows/` (NOT in subdirectories)
2. Workflows are on the **default branch** (`main`)
3. GitHub Actions is enabled: `Settings → Actions → General → Allow all actions`
4. File path filters match your service directory
5. At least one file matching the `paths:` filter was changed in the commit

### Docker Push Fails

Check that:
1. GitHub secrets are correctly set (`DOCKER_HUB_USERNAME`, `DOCKER_HUB_TOKEN`)
2. Docker Hub token has Read & Write permissions
3. Docker image name is correct

### Build Fails

Check that:
1. Service has `cmd/main.go` entry point
2. `go.mod` and `go.sum` files exist
3. All dependencies are available

## Current Services

| Service | Directory | CI Workflow | Docker Workflow | Docker Image | Status |
|---------|-----------|-------------|-----------------|--------------|--------|
| Auth | `auth/` | `auth-ci.yml` | `auth-docker.yml` | `hanapi23/auth-service` | Active |
| (Add more here as you create them) | | | | | |
