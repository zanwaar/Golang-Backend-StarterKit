.PHONY: help build run test clean deps fmt lint docker-build docker-run swagger

# Variables
BINARY_NAME=api-gin-production
GO_FILES=$(shell find . -name '*.go' | grep -v vendor)

help:
	@echo "Available commands:"
	@echo "  make deps          - Download and install dependencies"
	@echo "  make build         - Build binary"
	@echo "  make run           - Run application"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make test          - Run tests"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make db-migrate    - Run database migrations"
	@echo "  make swagger       - Generate Swagger documentation"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

# Build binary
build:
	@echo "🔨 Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(BINARY_NAME) main.go
	@echo "✓ Build complete: bin/$(BINARY_NAME)"

# Build for Windows
build-windows:
	@echo "🔨 Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe main.go
	@echo "✓ Build complete: bin/$(BINARY_NAME).exe"

# Build for macOS
build-macos:
	@echo "🔨 Building for macOS..."
	CGO_ENABLED=0 GOOS=darwin go build -o bin/$(BINARY_NAME)-macos main.go
	@echo "✓ Build complete: bin/$(BINARY_NAME)-macos"

# Run application
run:
	@echo "🚀 Running application..."
	go run main.go

# Run application in release mode
run-release:
	@echo "🚀 Running application in RELEASE mode..."
	GIN_MODE=release go run main.go

# Run with hot reload (requires: go install github.com/cosmtrek/air@latest)
dev:
	@echo "🔄 Running with hot reload..."
	$(shell go env GOPATH)/bin/air

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v -race -timeout 30s ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -race -timeout 30s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Lint code (requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@echo "🔍 Linting code..."
	golangci-lint run ./...

# Vet code
vet:
	@echo "🔍 Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "✓ Clean complete"

# Generate mocks (requires: go install github.com/golang/mock/cmd/mockgen@latest)
mocks:
	@echo "🎭 Generating mocks..."
	mockgen -source=repository/user_repository.go -destination=mocks/mock_user_repository.go
	@echo "✓ Mocks generated"

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "✓ Docker image built"

# Docker run
docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# Database migration (requires migration tool)
db-migrate:
	@echo "🗄️  Running database migrations..."
	@echo "🗄️  Running database migrations..."
	go run main.go -migrate
	@echo "✓ Migrations complete"

# Database seeder
db-seed:
	@echo "🌱 Running database seeder..."
	go run main.go -seed
	@echo "✓ Seeding complete"

# Generate Swagger docs
swagger:
	@echo "📝 Generating Swagger documentation..."
	$(shell go env GOPATH)/bin/swag init
	@echo "✓ Swagger docs generated"

# Install development tools
install-tools:
	@echo "📦 Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/cmd/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✓ Tools installed"

# Security Audit
audit:
	@echo "🛡️ Running security audit..."
	$(shell go env GOPATH)/bin/govulncheck ./...
	@echo "✓ Audit complete"

# Format, vet and lint
check: fmt vet lint audit
	@echo "✓ All checks passed"

# Full build and test
all: clean deps fmt vet build test
	@echo "✓ Build complete and all tests passed"
