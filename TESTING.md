# Testing Guide

This document outlines the testing strategy and how to run tests for FleetPass.

## Table of Contents

- [Testing Stack](#testing-stack)
- [Running Tests](#running-tests)
- [Writing Tests](#writing-tests)
- [CI/CD Pipeline](#cicd-pipeline)
- [Coverage Reports](#coverage-reports)

## Testing Stack

### Backend (Go)
- **Framework**: Go's built-in `testing` package
- **Test Database**: PostgreSQL (isolated test database)
- **Coverage**: Built-in Go coverage tools
- **Race Detection**: Enabled by default

### Frontend (React)
- **Framework**: Jest
- **Component Testing**: React Testing Library
- **Mocking**: Jest mocks
- **Coverage**: Istanbul (via Jest)

## Running Tests

### Quick Start

Run all tests:
```bash
make test
```

### Backend Tests

Run backend tests:
```bash
make test-backend
```

Or directly with Go:
```bash
go test -v ./...
```

Run with race detection:
```bash
go test -v -race ./...
```

Run specific package:
```bash
go test -v ./internal/handlers
```

### Frontend Tests

Run frontend tests:
```bash
make test-frontend
```

Or directly with npm:
```bash
cd frontend
npm test
```

Run once (no watch mode):
```bash
cd frontend
npm test -- --watchAll=false
```

### Coverage Reports

Generate coverage for both backend and frontend:
```bash
make test-coverage
```

Backend coverage only:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Frontend coverage only:
```bash
cd frontend
npm test -- --coverage --watchAll=false
```

## Writing Tests

### Backend Test Structure

Create test files alongside source files with `_test.go` suffix:

```go
package handlers

import (
    "testing"
    "fleetpass/internal/testutil"
)

func TestGetVehicles(t *testing.T) {
    // Setup test database
    db := testutil.SetupTestDB(t)
    defer testutil.CleanupTestDB(t, db)

    // Create test data
    org := testutil.CreateTestOrganization(t, db, "Test Org", "test-org")

    // Test logic
    // ...

    // Assertions
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### Frontend Test Structure

Create test files alongside components with `.test.js` suffix:

```javascript
import { render, screen, fireEvent } from '@testing-library/react';
import MyComponent from './MyComponent';

describe('MyComponent', () => {
  test('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });

  test('handles click events', () => {
    render(<MyComponent />);
    fireEvent.click(screen.getByRole('button'));
    // Assert expected behavior
  });
});
```

### Test Utilities

Backend test utilities are in `internal/testutil/`:
- `SetupTestDB(t)` - Creates test database connection
- `CleanupTestDB(t, db)` - Cleans up test data
- `CreateTestOrganization()` - Creates test organization
- `CreateTestLocation()` - Creates test location
- `CreateTestVehicle()` - Creates test vehicle

## Test Database Setup

Tests require a PostgreSQL test database. Configure in your environment:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fleetpass_test
export DB_USER=fleetpass_user
export DB_PASSWORD=fleetpass_password
export DB_SSLMODE=disable
```

Or run with Docker:
```bash
docker run -d \
  --name fleetpass-test-db \
  -e POSTGRES_DB=fleetpass_test \
  -e POSTGRES_USER=fleetpass_user \
  -e POSTGRES_PASSWORD=fleetpass_password \
  -p 5432:5432 \
  postgres:16-alpine
```

## CI/CD Pipeline

Our GitHub Actions workflow runs on every push and pull request:

1. **Backend Tests**
   - Runs Go tests with race detection
   - Generates coverage reports
   - Uploads to Codecov

2. **Frontend Tests**
   - Runs Jest tests
   - Generates coverage reports
   - Uploads to Codecov

3. **Build Validation**
   - Builds backend binary
   - Builds frontend production bundle

4. **Docker Images** (main branch only)
   - Builds and pushes Docker images
   - Tags with `latest` and commit SHA

5. **Security Scanning**
   - Runs Trivy vulnerability scanner
   - Uploads results to GitHub Security

## Coverage Targets

We aim for:
- **Backend**: 80%+ code coverage
- **Frontend**: 70%+ code coverage

Current coverage badges:
[![codecov](https://codecov.io/gh/yourusername/fleetpass/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/fleetpass)

## Best Practices

1. **Write tests first** (TDD approach when possible)
2. **Test behavior, not implementation**
3. **Keep tests isolated** - no dependencies between tests
4. **Use descriptive test names**
5. **Clean up after tests** - use defer for cleanup
6. **Mock external dependencies**
7. **Test error cases** - don't just test happy paths
8. **Keep tests fast** - use test databases, not production

## Debugging Tests

### Backend
```bash
# Verbose output
go test -v ./internal/handlers

# Run specific test
go test -v -run TestGetVehicles ./internal/handlers

# With debugging
dlv test ./internal/handlers -- -test.run TestGetVehicles
```

### Frontend
```bash
# Debug in watch mode
cd frontend
npm test

# Run specific test file
npm test -- VehicleBulkUpload.test.js

# With debugging
node --inspect-brk node_modules/.bin/jest --runInBand
```

## Continuous Improvement

- Review coverage reports regularly
- Add tests for new features
- Update tests when refactoring
- Remove obsolete tests
- Keep test utilities up to date
