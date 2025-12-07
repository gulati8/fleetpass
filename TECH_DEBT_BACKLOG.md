# FleetPass Technical Debt Backlog

**Generated:** December 6, 2025
**Source:** ARCHITECTURE_REVIEW.md
**Backlog Size:** 15 active items (max)

---

## Active Backlog

Items are scoped to **1-2 small PRs**, modify **≤3 code units**, and are focused on a **single concern**.

| ID | Title | Area | Business Impact | Tech Impact | Effort | Type |
|----|-------|------|-----------------|-------------|--------|------|
| **TD-001** | Move JWT secret to environment variable | Security/Correctness | High | High | Low | Implementation |
| **TD-002** | Fix connection pool bug in db.go | Security/Correctness | High | High | Low | Implementation |
| **TD-003** | Add org isolation to read vehicle handlers | Security/Correctness | High | High | Low | Implementation |
| **TD-004** | Add org isolation to write vehicle handlers | Security/Correctness | High | High | Medium | Implementation |
| **TD-005** | Add rate limiting to auth endpoints | Security/Correctness | High | High | Medium | Implementation |
| **TD-006** | Create VehicleRepository interface + implementation | Architecture/Maintainability | Medium | High | Medium | Implementation |
| **TD-007** | Extract org_id lookup to helper function | Architecture/Maintainability | Low | High | Low | Implementation |
| **TD-008** | Add go-playground/validator to vehicle requests | Architecture/Maintainability | Medium | Medium | Low | Implementation |
| **TD-009** | Create typed error handling with AppError | Architecture/Maintainability | Medium | Medium | Low | Implementation |
| **TD-010** | Add structured logging with zerolog | Architecture/Maintainability | Low | Medium | Low | Implementation |
| **TD-011** | Install TypeScript and create base types | DX/Infra | Medium | High | Low | Spike |
| **TD-012** | Create useVehicles custom hook | DX/Infra | Low | High | Low | Implementation |
| **TD-013** | Extract LoadingSpinner and ErrorAlert components | DX/Infra | Low | Medium | Low | Implementation |
| **TD-014** | Create typed vehiclesAPI namespace | DX/Infra | Low | High | Low | Implementation |
| **TD-015** | Add ErrorBoundary to App | DX/Infra | Medium | Medium | Low | Implementation |

---

## Item Details

### TD-001: Move JWT secret to environment variable

**Business Impact:** High - Production deployment blocker. Hardcoded secrets are security vulnerability.

**Description:**
Remove hardcoded JWT secret from `main.go:22`. Add `JWT_SECRET` environment variable with startup validation (fail fast if missing). Update Docker Compose with example secret placeholder.

**Files Modified:**
- `main.go` (line 22: replace hardcoded secret with `os.Getenv("JWT_SECRET")`)
- `docker-compose.yml` (add `JWT_SECRET` env var)
- `.env.example` (create with `JWT_SECRET=your-secret-here-change-in-production`)

**Definition of Done:**
- [ ] `main.go:22` reads `JWT_SECRET` from environment using `os.Getenv()`
- [ ] Application exits with clear error message if `JWT_SECRET` is empty on startup
- [ ] `docker-compose.yml` includes `JWT_SECRET` environment variable with placeholder value
- [ ] `.env.example` created with `JWT_SECRET` entry and comment explaining usage
- [ ] Manual test: Start app without JWT_SECRET → fails with helpful error
- [ ] Manual test: Start app with JWT_SECRET → login works correctly

---

### TD-002: Fix connection pool bug in db.go

**Business Impact:** High - Database connections expiring immediately causes performance degradation and connection exhaustion.

**Description:**
Fix `db.go:64-67` where `SetConnMaxLifetime(5 * 60)` is treating integer as nanoseconds (300ns) instead of 5 minutes. Change to `SetConnMaxLifetime(5 * time.Minute)`.

**Files Modified:**
- `internal/database/db.go` (lines 64-67)

**Definition of Done:**
- [ ] Line 66 changed from `sqlDB.SetConnMaxLifetime(5 * 60)` to `sqlDB.SetConnMaxLifetime(5 * time.Minute)`
- [ ] Line 67 changed from `sqlDB.SetConnMaxIdleTime(5 * 60)` to `sqlDB.SetConnMaxIdleTime(5 * time.Minute)`
- [ ] Import `time` package added to file if not present
- [ ] Manual test: Query PostgreSQL `pg_stat_activity` → connections stay alive for ~5 minutes
- [ ] No existing tests broken by change

---

### TD-003: Add org isolation to read vehicle handlers

**Business Impact:** High - Critical security gap. Users can view vehicles from other organizations.

**Description:**
Extract user's `organization_id` from JWT claims in `GetVehicles` and `GetVehicle` handlers. Filter queries to only return vehicles where `organization_id` matches user's org. Return 403 if org mismatch detected.

**Files Modified:**
- `internal/handlers/vehicle.go` (modify `GetVehicles` and `GetVehicle` functions only)

**Definition of Done:**
- [ ] `GetVehicles` extracts `organization_id` from JWT claims via `jwtauth.FromContext(r.Context())`
- [ ] `GetVehicles` adds `Where("organization_id = ?", userOrgID)` to GORM query
- [ ] `GetVehicle` extracts `organization_id` from JWT claims
- [ ] `GetVehicle` validates fetched vehicle's `organization_id` matches user's org, returns 403 if mismatch
- [ ] Manual test: User A logs in → calls GET /api/vehicles → only sees their org's vehicles
- [ ] Manual test: User A tries GET /api/vehicles/{vehicleFromOrgB} → receives 403 Forbidden
- [ ] Existing vehicle handler tests updated to mock JWT claims with org_id

---

### TD-004: Add org isolation to write vehicle handlers

**Business Impact:** High - Critical security gap. Users can modify/delete vehicles from other organizations.

**Description:**
Add `organization_id` validation to `CreateVehicle`, `UpdateVehicle`, and `DeleteVehicle` handlers. Ensure user can only create vehicles in their org and cannot modify/delete vehicles from other orgs.

**Files Modified:**
- `internal/handlers/vehicle.go` (modify `CreateVehicle`, `UpdateVehicle`, `DeleteVehicle` functions)

**Definition of Done:**
- [ ] `CreateVehicle` validates that location's `organization_id` matches user's org before creating vehicle
- [ ] `UpdateVehicle` extracts `organization_id` from JWT, validates vehicle belongs to user's org before update
- [ ] `DeleteVehicle` extracts `organization_id` from JWT, validates vehicle belongs to user's org before delete
- [ ] All three handlers return 403 Forbidden with clear error message on org mismatch
- [ ] Manual test: User A tries to delete vehicle from Org B → receives 403
- [ ] Manual test: User A creates vehicle with location from Org B → receives 403
- [ ] Handler tests updated to verify org isolation for create/update/delete operations

---

### TD-005: Add rate limiting to auth endpoints

**Business Impact:** High - Prevents brute force attacks on login. Required for production security.

**Description:**
Install `tollbooth` rate limiting middleware. Apply rate limit of 5 requests per minute per IP to `/api/login`, `/api/register`, `/api/forgot-password`, `/api/reset-password` endpoints.

**Files Modified:**
- `main.go` (add rate limiting middleware to auth routes)
- `go.mod` (add `github.com/didip/tollbooth` dependency)

**Definition of Done:**
- [ ] `go.mod` includes `github.com/didip/tollbooth/v7` dependency
- [ ] `main.go` imports tollbooth
- [ ] Rate limiter created with `tollbooth.NewLimiter(5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Minute})`
- [ ] Rate limiter middleware applied to `/api/login`, `/api/register`, `/api/forgot-password`, `/api/reset-password` routes
- [ ] Manual test: Send 6 login requests within 1 minute from same IP → 6th request returns 429 Too Many Requests
- [ ] Manual test: Wait 1 minute → requests work again
- [ ] Rate limiting does not apply to other endpoints (vehicles, locations, etc.)

---

### TD-006: Create VehicleRepository interface + implementation

**Business Impact:** Medium - Improves testability and enables future refactoring. Not user-facing.

**Description:**
Create `internal/repositories/vehicle.go` with `VehicleRepository` interface defining methods: `FindAll(ctx, orgID)`, `FindByID(ctx, id)`, `Create(ctx, vehicle)`, `Update(ctx, vehicle)`, `Delete(ctx, id)`. Implement with GORM. Do not integrate with handlers yet (that's a future PR).

**Files Modified:**
- `internal/repositories/vehicle.go` (new file)

**Definition of Done:**
- [ ] `internal/repositories/vehicle.go` created with package `repositories`
- [ ] `VehicleRepository` interface defined with 5 methods (FindAll, FindByID, Create, Update, Delete)
- [ ] `vehicleRepository` struct created with `db *gorm.DB` field
- [ ] `NewVehicleRepository(db *gorm.DB) VehicleRepository` constructor function
- [ ] All 5 methods implemented using GORM (e.g., `FindAll` uses `db.Where("organization_id = ?", orgID).Find(&vehicles)`)
- [ ] Unit tests created in `vehicle_test.go` for at least FindAll and Create methods
- [ ] Tests use in-memory SQLite or mock GORM for isolation
- [ ] Code compiles and tests pass (`go test ./internal/repositories`)

---

### TD-007: Extract org_id lookup to helper function

**Business Impact:** Low - Code quality improvement. Reduces duplication.

**Description:**
Extract `vehicle.go:50-55` (location lookup to get organization_id) into reusable function `getOrgIDFromLocation(locationID string) (string, error)` in new file `internal/handlers/helpers.go`. Use in `CreateVehicle` handler.

**Files Modified:**
- `internal/handlers/helpers.go` (new file)
- `internal/handlers/vehicle.go` (replace inline lookup with helper call in `CreateVehicle`)

**Definition of Done:**
- [ ] `internal/handlers/helpers.go` created with package `handlers`
- [ ] Function `getOrgIDFromLocation(locationID string) (string, error)` defined
- [ ] Function queries database for location, returns organization_id or error if not found
- [ ] `CreateVehicle` in `vehicle.go` calls `getOrgIDFromLocation(req.LocationID)` instead of inline query
- [ ] Original functionality preserved (error handling, HTTP status codes unchanged)
- [ ] Existing tests for `CreateVehicle` pass without modification
- [ ] Consider adding unit test for helper function (optional for this PR)

---

### TD-008: Add go-playground/validator to vehicle requests

**Business Impact:** Medium - Improves data quality. Prevents invalid vehicles from being created.

**Description:**
Install `go-playground/validator/v10`. Add `validate` struct tags to `CreateVehicleRequest` and `UpdateVehicleRequest` in `models/vehicle.go`. Add validation call in handlers before processing requests.

**Files Modified:**
- `internal/models/vehicle.go` (add `validate` tags to request structs)
- `internal/handlers/vehicle.go` (add validation calls in `CreateVehicle` and `UpdateVehicle`)
- `go.mod` (add validator dependency)

**Definition of Done:**
- [ ] `go.mod` includes `github.com/go-playground/validator/v10` dependency
- [ ] `CreateVehicleRequest` has validate tags: `VIN` tagged `validate:"required,len=17,alphanum"`, `Make` tagged `validate:"required,max=100"`, `Year` tagged `validate:"required,min=1900,max=2026"`
- [ ] Similar tags added to other required fields (Model, LocationID)
- [ ] `CreateVehicle` handler creates validator instance and calls `validate.Struct(req)`
- [ ] Validation errors returned as 400 Bad Request with field-specific error messages
- [ ] Manual test: POST vehicle with VIN="ABC" → returns 400 with error "VIN must be 17 characters"
- [ ] Manual test: POST vehicle with Year=1800 → returns 400 with error about year range
- [ ] Existing valid vehicle creation still works

---

### TD-009: Create typed error handling with AppError

**Business Impact:** Medium - Improves debugging and user experience. Better error messages.

**Description:**
Create `internal/errors/errors.go` with `AppError` struct containing `Code` (string), `Message` (string), and `HTTPStatus` (int). Define error codes: NOT_FOUND, VALIDATION_ERROR, FORBIDDEN, INTERNAL_ERROR. Add `WriteError(w, err)` helper function to serialize errors as JSON.

**Files Modified:**
- `internal/errors/errors.go` (new file)

**Definition of Done:**
- [ ] `internal/errors/errors.go` created with package `errors`
- [ ] `AppError` struct defined with fields: Code string, Message string, HTTPStatus int
- [ ] Constants defined: `ErrCodeNotFound = "NOT_FOUND"`, `ErrCodeValidation = "VALIDATION_ERROR"`, `ErrCodeForbidden = "FORBIDDEN"`, `ErrCodeInternal = "INTERNAL_ERROR"`
- [ ] Constructor functions: `NewNotFoundError(resource string)`, `NewValidationError(msg string)`, `NewForbiddenError(msg string)`, `NewInternalError(err error)`
- [ ] `WriteError(w http.ResponseWriter, err error)` function that type-switches on AppError and writes JSON response with appropriate HTTP status
- [ ] Unit test: `WriteError` with NotFoundError writes 404 status and JSON body with code/message
- [ ] Code compiles and tests pass

---

### TD-010: Add structured logging with zerolog

**Business Impact:** Low - Developer experience improvement. Better debugging in production.

**Description:**
Install `zerolog`. Replace `log.Println` statements in `vehicle.go` with structured logging. Add request ID middleware to `main.go` that generates UUID per request and adds to context.

**Files Modified:**
- `internal/handlers/vehicle.go` (replace log.Println with zerolog)
- `main.go` (add request ID middleware)
- `go.mod` (add zerolog and uuid dependencies)

**Definition of Done:**
- [ ] `go.mod` includes `github.com/rs/zerolog` and `github.com/google/uuid`
- [ ] `main.go` has middleware function that generates UUID, adds to request context, logs request start/end
- [ ] Request ID middleware applied to all routes in `main.go`
- [ ] `vehicle.go` imports zerolog and creates logger from request context
- [ ] At least 3 log statements replaced: request received, operation success, operation error
- [ ] Logs include fields: request_id, user_id (from JWT), action (e.g., "create_vehicle"), vehicle_id
- [ ] Manual test: Create vehicle → logs show JSON format with structured fields and request_id
- [ ] All log statements use appropriate levels (Info, Error, Warn)

---

### TD-011: Install TypeScript and create base types

**Business Impact:** Medium - Foundation for type safety. Reduces bugs in frontend development.

**Description:**
Install TypeScript dependencies in frontend. Create `tsconfig.json` with strict settings. Create type definition files in `src/types/` for Vehicle, User, Location, Organization matching backend models. No code migration yet - just setup and types.

**Files Modified:**
- `frontend/package.json` (add TS dependencies)
- `frontend/tsconfig.json` (new file)
- `frontend/src/types/vehicle.ts` (new file)
- `frontend/src/types/user.ts` (new file)
- `frontend/src/types/location.ts` (new file)
- `frontend/src/types/organization.ts` (new file)

**Definition of Done:**
- [ ] `package.json` includes: `typescript`, `@types/react`, `@types/react-dom`, `@types/node`
- [ ] `tsconfig.json` created with `strict: true`, `jsx: "react-jsx"`, `target: "ES2020"`
- [ ] `src/types/vehicle.ts` exports `Vehicle` interface with all fields matching backend model (id, vin, make, model, year, etc.)
- [ ] `src/types/user.ts` exports `User` interface (id, email, role, organization_id, etc.)
- [ ] `src/types/location.ts` exports `Location` interface
- [ ] `src/types/organization.ts` exports `Organization` interface
- [ ] All type files export request/response types (e.g., `CreateVehicleRequest`, `LoginResponse`)
- [ ] `npm run build` succeeds (TypeScript compiler runs without errors on type files)
- [ ] No existing .js files migrated yet

---

### TD-012: Create useVehicles custom hook

**Business Impact:** Low - Developer experience. Eliminates code duplication.

**Description:**
Extract duplicated fetch logic from `Vehicles.js` into `hooks/useVehicles.ts`. Hook manages vehicles state, loading state, error state, and provides refetch function. Return type is `{ vehicles: Vehicle[], loading: boolean, error: string | null, refetch: () => Promise<void> }`.

**Files Modified:**
- `frontend/src/hooks/useVehicles.ts` (new file)
- `frontend/src/pages/Vehicles.js` (replace manual fetch logic with hook)

**Definition of Done:**
- [ ] `frontend/src/hooks/useVehicles.ts` created (TypeScript file, no JSX)
- [ ] Hook imports `Vehicle` type from `../types/vehicle`
- [ ] Hook uses `useState` for vehicles, loading, error
- [ ] Hook uses `useEffect` to call `fetchVehicles` on mount
- [ ] `fetchVehicles` function calls `api.get('/api/vehicles')` and updates state
- [ ] Hook returns object with typed properties: `{ vehicles: Vehicle[], loading: boolean, error: string | null, refetch: () => Promise<void> }`
- [ ] `Vehicles.js` imports and uses hook: `const { vehicles, loading, error, refetch } = useVehicles()`
- [ ] Remove manual useState/useEffect from `Vehicles.js`
- [ ] Vehicle list still renders correctly, loading spinner works, errors display

---

### TD-013: Extract LoadingSpinner and ErrorAlert components

**Business Impact:** Low - Code quality. Reduces markup duplication.

**Description:**
Create reusable `LoadingSpinner.jsx` and `ErrorAlert.jsx` components to eliminate repeated markup across 5+ pages. Use Bootstrap classes for styling. Replace duplicated markup in `Vehicles.js` and `VehicleForm.js`.

**Files Modified:**
- `frontend/src/components/LoadingSpinner.jsx` (new file)
- `frontend/src/components/ErrorAlert.jsx` (new file)
- `frontend/src/pages/Vehicles.js` (use components)
- `frontend/src/pages/VehicleForm.js` (use components)

**Definition of Done:**
- [ ] `LoadingSpinner.jsx` created with Bootstrap spinner markup (`<div className="spinner-border">`)
- [ ] `LoadingSpinner` accepts no props, renders centered spinner with "Loading..." text
- [ ] `ErrorAlert.jsx` created accepting `message` prop (string) and optional `onDismiss` prop (function)
- [ ] `ErrorAlert` renders Bootstrap alert-danger with message and optional dismiss button
- [ ] `Vehicles.js` imports both components
- [ ] `Vehicles.js` renders `<LoadingSpinner />` when `loading === true`
- [ ] `Vehicles.js` renders `<ErrorAlert message={error} />` when `error` exists
- [ ] Same changes applied to `VehicleForm.js`
- [ ] Visual appearance unchanged from before
- [ ] Components are reusable (no hardcoded text or logic)

---

### TD-014: Create typed vehiclesAPI namespace

**Business Impact:** Low - Developer experience. Type safety for API calls.

**Description:**
Create `services/vehicles.api.ts` with typed functions `getAll()`, `getById(id)`, `create(data)`, `update(id, data)`, `delete(id)`. Use TypeScript generics with types from `types/vehicle.ts`. Return typed promises.

**Files Modified:**
- `frontend/src/services/vehicles.api.ts` (new file)

**Definition of Done:**
- [ ] `vehicles.api.ts` created (TypeScript file)
- [ ] Imports `api` from `./api`, `Vehicle`, `CreateVehicleRequest`, `UpdateVehicleRequest` from `../types/vehicle`
- [ ] Exports `vehiclesAPI` object with 5 methods
- [ ] `getAll(): Promise<Vehicle[]>` calls `api.get<Vehicle[]>('/api/vehicles')` and returns `response.data`
- [ ] `getById(id: string): Promise<Vehicle>` calls `api.get<Vehicle>(`/api/vehicles/${id}`)` and returns data
- [ ] `create(data: CreateVehicleRequest): Promise<Vehicle>` calls `api.post<Vehicle>('/api/vehicles', data)`
- [ ] `update(id: string, data: UpdateVehicleRequest): Promise<Vehicle>` calls `api.put<Vehicle>(...)`
- [ ] `delete(id: string): Promise<void>` calls `api.delete(...)`
- [ ] TypeScript compiler accepts file with no errors
- [ ] Can import and use: `import { vehiclesAPI } from '../services/vehicles.api'`

---

### TD-015: Add ErrorBoundary to App

**Business Impact:** Medium - Improves reliability. Prevents component crashes from breaking entire application.

**Description:**
Install `react-error-boundary` package. Create `ErrorPage.jsx` component to display when errors occur. Wrap `<Routes>` in `<ErrorBoundary fallback={<ErrorPage />}>` in `App.js`.

**Files Modified:**
- `frontend/package.json` (add react-error-boundary)
- `frontend/src/components/ErrorPage.jsx` (new file)
- `frontend/src/App.js` (wrap Routes in ErrorBoundary)

**Definition of Done:**
- [ ] `package.json` includes `react-error-boundary` dependency
- [ ] `ErrorPage.jsx` created with user-friendly error UI (Bootstrap alert-danger, "Something went wrong" message, button to reload page)
- [ ] `App.js` imports `ErrorBoundary` from `react-error-boundary` and `ErrorPage`
- [ ] `<Routes>` wrapped in `<ErrorBoundary fallback={<ErrorPage />}>`
- [ ] Manual test: Add `throw new Error("test")` to Vehicles.js → ErrorPage displays instead of blank screen
- [ ] Manual test: Click "Reload" button on ErrorPage → app recovers
- [ ] Remove test error after verification
- [ ] Console logs error details when boundary catches error

---

## Parking Lot

### Deferred Until After Core Refactor

These items depend on active backlog completion and core architectural patterns being in place.

**Blocked by TD-006 (VehicleRepository):**
- **Integrate VehicleRepository into handlers** - Refactor vehicle handlers to use repository instead of direct GORM calls. Requires TD-006 complete + tests passing.
- **Create VehicleService layer** - Business logic layer using repository. Requires TD-006 + TD-009 (AppError) complete. Would consolidate authorization checks, validation, and org_id logic.

**Blocked by Service Layer:**
- **Full Dependency Injection Pattern** - Inject repository and service into handlers via constructors instead of global database.DB. Requires service layer implemented first. Multi-week effort affecting all handlers.
- **Context Propagation** - Pass user_id, org_id, request_id through request context. Requires DI pattern so context can be passed to services/repos. Enables distributed tracing.

**Blocked by TD-011 (TypeScript Setup):**
- **Migrate Vehicles.js to Vehicles.tsx** - Convert page to TypeScript. Requires TD-011 (base types), TD-012 (useVehicles hook), TD-014 (typed API) complete.
- **Migrate all pages to TypeScript** - Full frontend migration. Multi-week effort. Do after proving pattern with Vehicles page.
- **Break down VehicleForm component** - Split 590-line component into VehicleBasicInfoCard, VehicleSpecsCard, VehiclePricingCard, etc. Requires TypeScript + React Hook Form. Each card is 2-3 PRs.

**Blocked by TypeScript Migration:**
- **Add React Hook Form + Zod** - Replace manual form state in VehicleForm. Requires TypeScript for type inference from Zod schemas. 1-2 week effort.
- **Add form validation schemas** - Create Zod schemas for all forms (vehicles, locations, orgs). Requires TypeScript + RHF integration complete.

**Infrastructure Dependencies:**
- **Audit Trail Middleware** - Populate audit_logs table for all CRUD operations. Requires context propagation (user_id in context) + structured logging. 3-5 days effort.
- **Implement soft delete pattern** - Add deleted_at column to vehicles, organizations, locations tables. Update all queries to filter `deleted_at IS NULL`. Create migration + GORM scopes. 3-5 days effort.

---

### Future Scale Considerations

These items address performance and scale but are not needed at current volume (<10K vehicles, <100 orgs).

**Performance:**
- **Add database read replicas** - Configure PostgreSQL streaming replication. Route read queries to replicas. Only needed when read QPS >1000.
- **Add Redis caching layer** - Cache vehicle lists, organization data. Adds complexity. Only needed if database becomes bottleneck.
- **Database partitioning** - Partition vehicles table by organization_id or created_at. Only needed at 1M+ vehicles.

**Observability:**
- **Prometheus metrics endpoint** - Expose `/metrics` endpoint with HTTP duration, request count, DB query time. Good for production monitoring but not blocking development.
- **Distributed tracing (Jaeger/Zipkin)** - Add trace IDs across services. Only needed if microservices architecture adopted.

**Advanced Features:**
- **Background job processing (asynq)** - For async bulk uploads, email sending, report generation. Current synchronous flow works fine for <100 uploads/day.
- **React Query** - Client-side caching, optimistic updates, request deduplication. Nice-to-have but current approach works. 1 week effort.

---

### Not Yet Prioritized

Items with unclear business value or technical necessity. Evaluate after completing core refactor.

**Frontend Nice-to-Haves:**
- **Add Storybook for component library** - Visual component documentation. Low priority, adds maintenance burden.
- **Migrate to Vite from CRA** - Faster builds. CRA works fine for now. 2-3 days effort for minimal gain.
- **Add React Testing Library tests** - Component tests. Some exist but coverage low. Ongoing effort.

**Backend Quality:**
- **Replace GORM with sqlc** - Architecture review mentioned sqlc preference for type safety. GORM works fine and team knows it. Large effort (2-3 weeks), unclear benefit.
- **Reorganize packages by domain** - Move from `handlers/`, `models/` to `vehicles/`, `organizations/`, `auth/` packages. Clean but not urgent. 1 week effort.

**Infrastructure:**
- **CI/CD Pipeline** - Automated testing, linting, deployment. Important long-term but team can deploy manually for now. 1-2 weeks setup.
- **Database backup automation** - Automated daily backups to S3. Important for production but not blocking development. 2-3 days setup.
- **Environment-based CORS config** - Move allowed origins from hardcoded array to environment variable. Low risk, can defer. 1 hour effort.

**Security Enhancements:**
- **PostgreSQL Row-Level Security** - Defense-in-depth for multi-tenancy. Application-level checks (TD-003, TD-004) are sufficient. RLS adds complexity. Defer until multi-tenant isolation is proven concern.
- **Add security headers middleware** - CSP, X-Frame-Options, HSTS. Good practice but low risk for API-only backend. 2-3 hours effort.
- **Input sanitization for XSS** - Sanitize HTML input using bluemonday. GORM prevents SQL injection. XSS risk low since no server-side rendering. Can defer.

**Code Cleanup:**
- **Remove dead code** - Delete `frontend/src/api.js` (unused), run `goimports` and `eslint --fix`. Low impact. 1 hour effort.
- **Standardize error messages** - Ensure consistent error format across all handlers. Cosmetic improvement. Can defer.

---

## Suggested Execution Order

### Phase 1: Critical Security (Week 1)
**Goal:** Close security gaps, unblock production deployment.

| Days | Items | Why |
|------|-------|-----|
| Day 1 | TD-001, TD-002 | Quick wins, high impact, low effort |
| Day 2-3 | TD-003 | Critical: prevent cross-org data access on reads |
| Day 4 | TD-004 | Critical: prevent cross-org data modification |
| Day 5 | TD-005 | Prevent brute force attacks |

**Deliverable:** Application can be deployed to production with confidence.

---

### Phase 2: Architecture Foundation (Week 2)
**Goal:** Establish patterns for cleaner code and better testing.

| Days | Items | Why |
|------|-------|-----|
| Day 1 | TD-006 | Repository pattern enables testability |
| Day 2 | TD-007, TD-008 | Code quality improvements, use new patterns |
| Day 3 | TD-009, TD-010 | Better error handling and logging |
| Day 4-5 | TD-011 | TypeScript foundation for frontend work |

**Deliverable:** Repository pattern proven, TypeScript setup complete.

---

### Phase 3: Frontend Modernization (Week 3)
**Goal:** Reduce frontend duplication, add type safety.

| Days | Items | Why |
|------|-------|-----|
| Day 1 | TD-012 | Custom hooks pattern proven |
| Day 2 | TD-013 | Extract reusable components |
| Day 3 | TD-014 | Typed API layer for frontend-backend contract |
| Day 4 | TD-015 | Error boundaries for reliability |
| Day 5 | Buffer/testing | Catch up, integration testing |

**Deliverable:** Frontend has TypeScript foundation, hooks pattern, reusable components.

---

## Notes

- **Scope discipline:** Each item modifies ≤3 files or ≤3 handlers. TD-003 split into read/write operations to maintain this constraint.
- **Parallel tracks:** Security (TD-001 to TD-005) can run parallel to architecture work. Backend (TD-006 to TD-010) and frontend (TD-011 to TD-015) can run concurrently after security work complete.
- **No epics in active backlog:** Large efforts like "Migrate to TypeScript" or "Add Service Layer" broken into small PRs or deferred to parking lot until dependencies complete.
- **Business impact:** High = affects end users or security, Medium = affects developers or operations, Low = code quality only.
- **Tech impact:** High = foundational change enabling future work, Medium = improves maintainability, Low = minor improvement.

---

**Document Version:** 3.0
**Last Updated:** December 6, 2025
**Next Review:** After completing Phase 1 (5 security items)
