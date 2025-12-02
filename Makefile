.PHONY: help test test-backend test-frontend test-coverage build run clean docker-up docker-down docker-rebuild lint fmt

# Default target
help:
	@echo "FleetPass - Available Commands:"
	@echo "  make test              - Run all tests (backend + frontend)"
	@echo "  make test-backend      - Run backend Go tests"
	@echo "  make test-frontend     - Run frontend React tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make build             - Build backend and frontend"
	@echo "  make run               - Run the application"
	@echo "  make lint              - Run linters"
	@echo "  make fmt               - Format code"
	@echo "  make docker-up         - Start Docker containers"
	@echo "  make docker-down       - Stop Docker containers"
	@echo "  make docker-rebuild    - Rebuild and restart Docker containers"
	@echo "  make clean             - Clean build artifacts"

# Testing
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	go test -v -race ./...

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm test -- --watchAll=false

test-coverage:
	@echo "Running backend tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo ""
	@echo "Running frontend tests with coverage..."
	cd frontend && npm test -- --coverage --watchAll=false

# Building
build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	go build -o bin/fleetpass-api .

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Running
run:
	@echo "Starting application..."
	docker-compose up

# Linting
lint: lint-backend lint-frontend

lint-backend:
	@echo "Running Go linters..."
	go vet ./...
	gofmt -s -l .

lint-frontend:
	@echo "Running frontend linters..."
	cd frontend && npm run lint --if-present

# Formatting
fmt:
	@echo "Formatting Go code..."
	gofmt -s -w .

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose up -d --build

docker-logs:
	docker-compose logs -f

# Database commands
db-migrate:
	@echo "Running database migrations..."
	@echo "Migrations are auto-applied on startup"

db-seed:
	@echo "Seeding database..."
	# Add seed command here when ready

# Cleaning
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf frontend/build
	rm -rf frontend/coverage
	go clean -cache

# Development setup
setup:
	@echo "Setting up development environment..."
	go mod download
	cd frontend && npm install

# Install tools
install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
