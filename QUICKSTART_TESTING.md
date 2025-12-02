# Testing & CI/CD Quick Start

## Run Tests Locally

### All Tests
```bash
make test
```

### Backend Only
```bash
make test-backend
# or
go test -v ./...
```

### Frontend Only
```bash
make test-frontend
# or
cd frontend && npm test
```

### With Coverage
```bash
make test-coverage
```

## What's Included

### Backend Tests ✅
- **Unit Tests**: `internal/handlers/vehicle_test.go`
- **Bulk Upload Tests**: `internal/handlers/vehicle_bulk_test.go`
- **Test Utilities**: `internal/testutil/testutil.go`
- Coverage: ~80% target

### Frontend Tests ✅
- **Auth Context**: `frontend/src/context/AuthContext.test.js`
- **API Service**: `frontend/src/services/api.test.js`
- **Bulk Upload**: `frontend/src/pages/VehicleBulkUpload.test.js`
- Coverage: ~70% target

### CI/CD Pipeline ✅
- **File**: `.github/workflows/ci.yml`
- **Triggers**: Push to main/develop, Pull Requests
- **Steps**:
  1. Run backend tests with race detection
  2. Run frontend tests
  3. Build both applications
  4. Generate coverage reports
  5. Build Docker images (main branch only)
  6. Run security scans

## Quick Commands

```bash
# Development
make setup          # Install dependencies
make run            # Start with docker-compose
make docker-rebuild # Rebuild containers

# Testing
make test           # Run all tests
make lint           # Run linters
make fmt            # Format code

# Coverage
make test-coverage  # Generate coverage reports
```

## GitHub Actions Setup

Add these secrets to your GitHub repository:

1. `DOCKER_USERNAME` - Docker Hub username
2. `DOCKER_PASSWORD` - Docker Hub password/token

Go to: Repository → Settings → Secrets and variables → Actions

## Next Steps

1. **Add More Tests**: Expand test coverage for all handlers
2. **E2E Tests**: Add Cypress or Playwright for end-to-end testing
3. **Performance Tests**: Add load testing with k6 or Artillery
4. **Deploy**: Follow DEPLOYMENT.md for production deployment

## Documentation

- **Full Testing Guide**: [TESTING.md](TESTING.md)
- **Deployment Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)

## Test Database Setup (Local)

For local testing, start a test database:

```bash
docker run -d \
  --name fleetpass-test-db \
  -e POSTGRES_DB=fleetpass_test \
  -e POSTGRES_USER=fleetpass_user \
  -e POSTGRES_PASSWORD=fleetpass_password \
  -p 5432:5432 \
  postgres:16-alpine
```

Or use the existing development database (tests will use a separate test database).

## Troubleshooting

**Tests fail with database connection error:**
- Ensure PostgreSQL is running
- Check environment variables match test database config
- Verify firewall allows connection to port 5432

**Frontend tests fail:**
- Clear node_modules: `cd frontend && rm -rf node_modules && npm install`
- Clear Jest cache: `cd frontend && npm test -- --clearCache`

**Coverage too low:**
- Add more test cases for uncovered code paths
- Run `go tool cover -html=coverage.out` to see coverage visualization
