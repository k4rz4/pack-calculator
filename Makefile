# Pack Calculator Makefile

# Variables
APP_NAME=pack-calculator
DOCKER_IMAGE=pack-calculator:latest
# DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE=docker compose
GO_VERSION=1.22

# Colors for output
YELLOW=\033[1;33m
GREEN=\033[1;32m
RED=\033[1;31m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-run docker-stop lint fmt vet deps migrate

# Default target
help: ## Show this help message
	@echo "$(YELLOW)Pack Calculator Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the application
	@echo "$(YELLOW)Building application...$(NC)"
	go build -o bin/$(APP_NAME) ./cmd/api

run: ## Run the application locally
	@echo "$(YELLOW)Running application...$(NC)"
	go run ./cmd/api

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	go test -v -race -tags=integration ./tests/integration/...

bench: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run

fmt: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	go vet ./...

deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy

# Docker
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application with Docker Compose
	@echo "$(YELLOW)Starting application with Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) up -d

docker-stop: ## Stop Docker containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	@echo "$(YELLOW)Showing Docker logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f api

docker-restart: docker-stop docker-run ## Restart Docker containers

# Development environment
dev-setup: deps ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)

dev: ## Run with hot reload (requires air)
	@echo "$(YELLOW)Starting development server with hot reload...$(NC)"
	air

# Cleanup
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f

clean-all: clean ## Clean everything including Docker images
	@echo "$(YELLOW)Cleaning everything...$(NC)"
	docker image rm $(DOCKER_IMAGE) 2>/dev/null || true
	$(DOCKER_COMPOSE) down -v --rmi all

# Health checks
health: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	curl -f http://pack-calculator.localhost/health || (echo "$(RED)Health check failed$(NC)" && exit 1)
	@echo "$(GREEN)Application is healthy$(NC)"

# Load testing
load-test: ## Run load tests (requires artillery or similar)
	@echo "$(YELLOW)Running load tests...$(NC)"
	@echo "Load testing not implemented yet"

# API testing
api-test: ## Test API endpoints
	@echo "$(YELLOW)Testing API endpoints...$(NC)"
	# Test health endpoint
	curl -f http://pack-calculator.localhost/health
	# Test pack creation
	curl -X POST http://pack-calculator.localhost/api/v1/packs \
		-H "Content-Type: application/json" \
		-d '{"size": 250, "name": "Small Pack"}'
	# Test calculation
	curl -X POST http://pack-calculator.localhost/api/v1/calculate \
		-H "Content-Type: application/json" \
		-d '{"pack_sizes": [250, 500, 1000], "order_quantity": 263}'

