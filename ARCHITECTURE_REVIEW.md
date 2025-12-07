# FleetPass Architecture Review

**Date:** December 6, 2025
**Reviewer:** Archie (Senior Software Architect)
**Tech Stack:** Go 1.22+ Backend, React 19 Frontend, PostgreSQL, Docker

---

## 1. Overview

### What FleetPass Does Today

FleetPass is a **multi-tenant fleet vehicle management platform** that enables organizations to:

- Manage vehicle inventory across multiple locations
- Track vehicle details (VIN, make, model, specifications, pricing)
- Implement role-based access control (RBAC) with admin, org_owner, manager, and user roles
- Perform bulk vehicle imports via CSV
- Manage organizations and locations in a hierarchical structure
- Authenticate users with JWT-based session management

### High-Level Architecture

```
┌─────────────────┐
│  React Frontend │  (JavaScript, Bootstrap 5, Axios)
│   Port: 3000    │
└────────┬────────┘
         │ HTTP/JSON
         ↓
┌─────────────────┐
│   Nginx Proxy   │  (Reverse proxy, static file serving)
│   Port: 80      │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  Go Backend API │  (Chi router, GORM, JWT auth)
│   Port: 8080    │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│   PostgreSQL    │  (Primary data store)
│   Port: 5432    │
└─────────────────┘
```

**Infrastructure:**
- Containerized with Docker Compose
- SQL migrations for schema management
- Environment-based configuration
- No CI/CD visible in repository

---

## 2. Backend (Go)

### Current Architecture

**Layers:**
- **HTTP Layer:** Chi router with middleware (CORS, logging, JWT auth)
- **Handlers:** HTTP request/response handling (`internal/handlers/`)
- **Models:** GORM-based domain models (`internal/models/`)
- **Database:** Global GORM instance (`database.DB`)

**Request Flow:**
```
HTTP Request → Chi Router → Handler → database.DB (GORM) → PostgreSQL
```

**Key Technologies:**
- Router: `go-chi/chi` ✓
- ORM: GORM
- Auth: `go-chi/jwtauth`
- Database: PostgreSQL with GORM AutoMigrate
- Config: Environment variables

### Strengths

1. ✓ **Chi router** - Matches preferred stack, clean routing structure
2. ✓ **JWT authentication** - Proper token-based auth with protected routes
3. ✓ **GORM models** - Type-safe database models with proper indexes
4. ✓ **Middleware stack** - Logger, Recoverer, CORS properly configured
5. ✓ **Database schema** - Well-designed with foreign keys, indexes, UUIDs
6. ✓ **SQL migrations** - Version-controlled schema in `migrations/001_initial_schema.sql`
7. ✓ **Docker Compose** - Reproducible development environment
8. ✓ **Separation of concerns** - Clear directory structure (`handlers/`, `models/`, `database/`)

### Weaknesses / Misalignments

#### Critical Issues

1. **No Service Layer** (Violation: Business logic in handlers)
   - Example: `vehicle.go:50-55` - Handler fetches location to extract organization_id
   - Business rules scattered across HTTP layer
   - Cannot reuse logic outside HTTP context

2. **No Repository Layer** (Violation: Direct database coupling)
   - Example: `vehicle.go:15` - `database.DB.Order("created_at DESC").Find(&vehicles)`
   - Handlers directly call GORM
   - Impossible to test without real database
   - Cannot swap implementations (e.g., caching layer)

3. **Global Database Variable** (`database/db.go:14`)
   - `var DB *gorm.DB` - Anti-pattern
   - Tight coupling across application
   - No dependency injection
   - Hard to test, hard to mock

4. **No Authorization Enforcement**
   - JWT validates identity, but no permission checks
   - Any authenticated user can access any organization's data
   - RBAC models exist (`roles`, `permissions`) but unused in handlers

5. **Manual Validation** (`vehicle.go:45-48`)
   - String-based validation: `if req.LocationID == "" || req.VIN == ""`
   - Should use `go-playground/validator` with struct tags
   - No VIN format validation, no business rule validation

#### High Priority

6. **Hardcoded JWT Secret** (`main.go:22`)
   - `jwtauth.New("HS256", []byte("your-secret-key-change-this-in-production"), nil)`
   - Security risk, should be environment variable

7. **Poor Error Handling**
   - Generic HTTP errors: `http.Error(w, "Failed to fetch vehicles", 500)`
   - No typed domain errors
   - Implementation details leaked to client

8. **No Context Propagation**
   - Handlers don't use `r.Context()`
   - No user ID, org ID, or request ID in context
   - Cannot trace requests or implement per-request transactions

9. **No Input Sanitization**
   - Direct assignment from request to model
   - Potential for mass assignment vulnerabilities

10. **No Structured Logging**
    - Using `log.Println` instead of zerolog/slog
    - No request IDs, user context, or structured fields

#### Medium Priority

11. **GORM vs sqlc Decision**
    - Using GORM but system prompt prefers sqlc for type safety
    - GORM chosen but not justified
    - Dynamic querying not heavily used

12. **No Background Job Infrastructure**
    - System prompt mentions asynq or go-work
    - Bulk operations (CSV upload) run synchronously
    - No async processing for emails, reports

13. **Flat Package Structure**
    - All handlers in one directory
    - Should group by domain (`vehicles/`, `organizations/`, `auth/`)

14. **No Integration Tests**
    - Unit tests exist for handlers (`vehicle_test.go`)
    - Require mocking global `database.DB`
    - No end-to-end API tests

### Top Recommended Backend Refactors (Ranked)

| Rank | Refactor | Impact | Effort | ROI |
|------|----------|--------|--------|-----|
| 1 | **Introduce Repository Layer** | HIGH | MEDIUM | ⭐⭐⭐⭐⭐ |
| 2 | **Introduce Service Layer** | HIGH | MEDIUM | ⭐⭐⭐⭐⭐ |
| 3 | **Implement Authorization Checks** | HIGH | MEDIUM | ⭐⭐⭐⭐ |
| 4 | **Add Dependency Injection** | HIGH | HIGH | ⭐⭐⭐⭐ |
| 5 | **Add Structured Error Handling** | MEDIUM | LOW | ⭐⭐⭐⭐ |
| 6 | **Implement go-playground/validator** | MEDIUM | LOW | ⭐⭐⭐ |
| 7 | **Add Context Propagation** | MEDIUM | MEDIUM | ⭐⭐⭐ |
| 8 | **Move JWT Secret to Environment** | HIGH | LOW | ⭐⭐⭐ |
| 9 | **Add Structured Logging (zerolog/slog)** | MEDIUM | LOW | ⭐⭐⭐ |
| 10 | **Add Background Job Processing** | LOW | HIGH | ⭐⭐ |

**Recommended Implementation Order:**
1. Repository layer (enables testing)
2. Service layer (centralizes business logic)
3. Dependency injection (removes global state)
4. Authorization checks (security critical)
5. Validation framework (data integrity)
6. Context propagation (observability)

---

## 3. Frontend (React + JavaScript → TypeScript)

### Current Architecture and Main Flows

**Tech Stack:**
- React 19.2.0 (latest)
- JavaScript (no TypeScript)
- React Router v7
- Axios for API calls
- Bootstrap 5.3.8 ✓
- Context API for auth state

**Component Structure:**
```
App.js (Routing + AuthProvider)
├─ Navbar.js
├─ ProtectedRoute.js
└─ Pages:
   ├─ Login.js
   ├─ Dashboard.js
   ├─ Vehicles.js (list)
   ├─ VehicleForm.js (create/edit - 590 lines!)
   ├─ VehicleProfile.js (detail)
   ├─ VehicleBulkUpload.js
   ├─ Organizations.js
   └─ Locations.js
```

**State Management:**
- **Auth:** Context API (`AuthContext.js`)
- **Data:** Local component state (`useState`)
- **Persistence:** localStorage (token, user JSON)
- **No global store** (Redux, Zustand, etc.)

**API Layer:**
- Axios instance with interceptors (`services/api.js`)
- Token injection via request interceptor
- Auto-redirect on 401/403 via response interceptor
- Partial API namespacing (`authAPI.login()` but raw `api.get('/api/vehicles')`)

**Main User Flows:**

1. **Authentication:**
   ```
   /login → authAPI.login() → Store token + user in localStorage
   → Navigate to /dashboard
   ```

2. **Vehicle CRUD:**
   ```
   /vehicles → useEffect → api.get('/api/vehicles') → Render table
   /vehicles/new → VehicleForm (590 lines, 42 form fields) → api.post()
   /vehicles/:id/edit → Fetch vehicle → Populate form → api.put()
   /vehicles/:id → Fetch vehicle → Render profile
   ```

3. **Session Expiry:**
   ```
   API returns 401 → Axios interceptor → Clear localStorage
   → window.location.href = '/login'
   ```

### Strengths

1. ✓ **Bootstrap 5** - Matches preferred stack, responsive design
2. ✓ **React 19** - Latest version, modern hooks
3. ✓ **Protected Routes** - ProtectedRoute wrapper prevents unauthorized access
4. ✓ **Axios Interceptors** - Centralized token injection and session handling
5. ✓ **Auth Context** - Centralized authentication state
6. ✓ **React Router v7** - Modern routing with proper nesting
7. ✓ **Loading States** - Spinners shown during data fetches
8. ✓ **Error Alerts** - User-friendly error messages displayed
9. ✓ **Responsive Grid** - Bootstrap grid system used consistently
10. ✓ **Clean Component Separation** - Pages vs Components vs Context

### Weaknesses / Misalignments

#### Critical Issues

1. **No TypeScript** (PRIMARY ISSUE)
   - All files are `.js` instead of `.tsx`
   - No type safety on props, API responses, or state
   - Example: `VehicleForm.js:93` - `vehicle.features ? vehicle.features.join('\n') : ''`
     - Is `features` always an array? Could be null, undefined, or something else - no way to know
   - Refactoring is risky without types

2. **Massive Form Component** (`VehicleForm.js` - 590 lines)
   - Single component handles:
     - Fetching locations
     - Fetching existing vehicle (edit mode)
     - Managing 42 form fields
     - Validation
     - Data transformation (arrays ↔ newline-separated strings)
     - Five Bootstrap cards
   - Should be split into 5+ smaller components

3. **No Custom Hooks** (Violation: DRY principle)
   - Repeated fetch pattern in every component:
     ```javascript
     const [data, setData] = useState([]);
     const [loading, setLoading] = useState(true);
     const [error, setError] = useState('');
     useEffect(() => { fetchData(); }, []);
     ```
   - Should have `useVehicles()`, `useVehicle(id)`, `useLocations()`, etc.

4. **Implicit API Contracts**
   - No TypeScript interfaces for API requests/responses
   - Example: What fields does `Vehicle` have? Must read backend code or inspect runtime
   - Frontend/backend contract violations only discovered at runtime

5. **No Validation Library**
   - Only HTML5 validation (`required`, `min`, `max`)
   - No VIN format validation (must be 17 chars)
   - No cross-field validation
   - No custom error messages
   - Should use React Hook Form + Zod

#### High Priority

6. **Duplicated Logic Across Components**
   - Loading spinner markup repeated 5+ times
   - Error alert markup repeated 5+ times
   - Fetch-data-on-mount pattern repeated 10+ times
   - Table rendering logic not extracted

7. **Inefficient Re-renders**
   - `Vehicles.js:33` - After delete, refetches entire vehicle list
   - No optimistic updates
   - No caching (React Query would help)

8. **Inconsistent API Patterns**
   - Auth endpoints: `authAPI.login()` (namespaced)
   - Other endpoints: `api.get('/api/vehicles')` (raw)
   - Should have `vehiclesAPI.getAll()`, `locationsAPI.getAll()`, etc.

9. **localStorage Anti-pattern** (`AuthContext.js:20-24`)
   ```javascript
   const savedUser = localStorage.getItem('user');
   setUser(JSON.parse(savedUser));
   ```
   - Storing entire user object as JSON
   - Should derive user from JWT claims
   - Type safety lost across serialization boundary

10. **No Error Boundaries**
    - Component errors crash entire app
    - Should wrap routes in ErrorBoundary

11. **Side Effects in Interceptor** (`api.js:31`)
    ```javascript
    window.location.href = '/login';  // Hard redirect
    ```
    - Interceptor directly manipulates browser location
    - Hard to test
    - Should emit event or throw typed error

#### Medium Priority

12. **Large Page Components**
    - `Vehicles.js` is 165 lines but could be 40 lines with extracted components
    - Table rendering should be `<VehicleTable />` component
    - Actions should be extracted

13. **No Form State Management**
    - VehicleForm uses raw `useState` for 42 fields
    - Should use React Hook Form for performance and validation

14. **No Request Deduplication**
    - Multiple components can fetch same data simultaneously
    - React Query would solve this

15. **No Optimistic Updates**
    - User waits for server response on every action
    - Should optimistically update UI, rollback on error

16. **Old api.js File Unused** (`src/api.js`)
    - Legacy file with only `getHello()`
    - Dead code should be removed

### Top Recommended Frontend Refactors (Ranked)

| Rank | Refactor | Impact | Effort | ROI |
|------|----------|--------|--------|-----|
| 1 | **Introduce TypeScript** | HIGH | HIGH | ⭐⭐⭐⭐⭐ |
| 2 | **Create Custom Hooks** (useVehicles, etc.) | HIGH | LOW | ⭐⭐⭐⭐⭐ |
| 3 | **Extract Reusable Components** | MEDIUM | LOW | ⭐⭐⭐⭐ |
| 4 | **Add Typed API Layer** | HIGH | MEDIUM | ⭐⭐⭐⭐ |
| 5 | **Break Down VehicleForm** | MEDIUM | MEDIUM | ⭐⭐⭐⭐ |
| 6 | **Add React Hook Form + Zod** | MEDIUM | MEDIUM | ⭐⭐⭐⭐ |
| 7 | **Implement Error Boundaries** | MEDIUM | LOW | ⭐⭐⭐ |
| 8 | **Add React Query** (optional) | MEDIUM | MEDIUM | ⭐⭐⭐ |
| 9 | **Optimize Renders** (optimistic updates) | LOW | MEDIUM | ⭐⭐ |
| 10 | **Clean Up Dead Code** | LOW | LOW | ⭐⭐ |

**Recommended Implementation Order:**
1. Install TypeScript, create type definitions
2. Migrate API layer to TypeScript (`services/*.ts`)
3. Create custom hooks (`hooks/useVehicles.ts`)
4. Extract reusable components (LoadingSpinner, ErrorAlert, VehicleTable)
5. Migrate pages to TypeScript (start with simplest)
6. Break down VehicleForm into smaller components
7. Add React Hook Form + Zod for validation

---

## 4. Data Model & Persistence

### Key Entities and Relationships

```
Organizations (Multi-tenant root)
    │
    ├─── Locations (1:N)
    │       │
    │       └─── Vehicles (1:N)
    │
    └─── Users (1:N)
             │
             └─── Roles → Permissions (RBAC)

Audit Logs (track all changes)
```

**Schema Highlights:**

1. **organizations** (UUID PK, slug unique index)
   - Root tenant entity
   - Has `is_active` flag

2. **locations** (UUID PK, FK to organizations)
   - Belongs to one organization
   - Contains address, contact info
   - Vehicles belong to locations

3. **vehicles** (UUID PK, FK to organizations + locations)
   - Comprehensive spec fields (VIN, make, model, year, mileage, etc.)
   - JSONB for features and images
   - Status enum: available, rented, maintenance, inactive
   - Pricing fields: daily_rate, weekly_rate, monthly_rate
   - **Unique constraint:** VIN (globally unique)
   - **Composite index:** (make, model)

4. **users** (UUID PK, FK to organizations)
   - Email unique constraint
   - Role enum: admin, org_owner, manager, user
   - Password hash stored

5. **roles & permissions** (RBAC models exist but unused)
   - Defined in `internal/models/role.go` and `permission.go`
   - Not enforced in handlers

6. **audit_logs** (UUID PK, JSONB for changes)
   - Tracks user actions
   - Not currently populated (no audit middleware)

### Strengths

1. ✓ **UUID Primary Keys** - Prevents enumeration attacks, distributed-friendly
2. ✓ **Proper Indexes** - ON (organization_id, location_id, vin, status, make/model, year)
3. ✓ **Foreign Keys with Cascades** - Data integrity enforced at DB level
4. ✓ **JSONB for Arrays** - Flexible storage for features/images
5. ✓ **Enum Constraints** - `CHECK (condition IN (...))` validates at DB level
6. ✓ **Timestamps** - `created_at`, `updated_at` with triggers
7. ✓ **Audit Log Table** - Foundation for compliance/tracking
8. ✓ **Multi-tenancy Ready** - Organization-based data isolation

### Risks and Issues

#### High Priority

1. **No Row-Level Security (RLS)**
   - PostgreSQL supports RLS but not configured
   - Application must enforce org-level isolation
   - Risk: Query bug could leak data across organizations

2. **No Soft Deletes**
   - `DELETE` operations are hard deletes
   - Cannot recover accidentally deleted vehicles
   - Audit trail incomplete (deletion not logged)

3. **Audit Logs Not Populated**
   - Table exists but no middleware/triggers to populate
   - Cannot track who created/modified/deleted records
   - Compliance risk for regulated industries

4. **No Data Validation at DB Level (Some Fields)**
   - VIN should be CHAR(17) not VARCHAR(17)
   - License plate formats vary by jurisdiction
   - Year should have CHECK constraint (currently only in Go validation)

#### Medium Priority

5. **Vehicle-Location Relationship** is `ON DELETE RESTRICT`
   - Cannot delete location if vehicles exist
   - May cause operational issues
   - Consider soft delete or cascade with warning

6. **No Composite Unique Constraints**
   - Example: (organization_id, license_plate) should be unique per org
   - Currently license_plate has no uniqueness constraint

7. **GORM AutoMigrate vs SQL Migrations**
   - Using both `AutoMigrate` (db.go:88) and SQL migrations (001_initial_schema.sql)
   - Risk: Schema drift between AutoMigrate and SQL
   - Recommendation: Pick one approach (prefer SQL migrations)

8. **No Partitioning Strategy**
   - As vehicles table grows, queries may slow
   - Consider partitioning by organization_id or created_at

9. **JSONB Arrays for Features/Images**
   - No schema validation on JSONB content
   - Frontend sends arrays, DB stores as JSONB
   - Risk: Corrupt data if frontend sends wrong format

#### Low Priority

10. **Connection Pool Configuration** (db.go:64-67)
    - `SetConnMaxLifetime(5 * 60)` - Sets 5 minutes but expects time.Duration
    - Should be `5 * time.Minute` (currently treated as nanoseconds = 300ns!)
    - Bug: Connections likely expiring immediately

11. **No Database Read Replicas**
    - All queries hit primary
    - Read-heavy workload (vehicle lists) could use replicas

12. **No Database Backups Visible**
    - No backup/restore procedures documented
    - Should have automated backups in production

---

## 5. Cross-Cutting Concerns

### Authentication & Authorization

**Current State:**

✓ **Authentication (Identity):**
- JWT-based authentication using `go-chi/jwtauth`
- Tokens include user_id, email, role, organization_id claims
- Frontend stores token in localStorage
- Axios interceptor injects `Authorization: Bearer <token>` header
- Backend verifies token via middleware (`main.go:68-69`)
- 401/403 triggers auto-redirect to login

✗ **Authorization (Permissions):**
- RBAC models exist (`roles`, `permissions` tables) but **not enforced**
- No permission checks in handlers
- Example: Any authenticated user can:
  - Create vehicles in any organization
  - Delete vehicles in any organization
  - Access any organization's data
- **CRITICAL SECURITY GAP**

**Issues:**

1. **No Organization Isolation**
   - User in Org A can access Org B's vehicles via API
   - Backend trusts `organization_id` from request, doesn't validate against JWT

2. **Role Field Unused**
   - JWT contains `role` claim but handlers don't check it
   - Admin vs user distinction not enforced

3. **Hardcoded JWT Secret** (`main.go:22`)
   - Should be in environment variable
   - Rotation not possible

**Recommendations:**

1. Add authorization middleware:
   ```go
   func RequirePermission(permission string) func(next http.Handler) http.Handler
   func RequireOrgMatch(orgID string) func(next http.Handler) http.Handler
   ```

2. Enforce at service layer:
   ```go
   func (s *VehicleService) CreateVehicle(ctx context.Context, userID, orgID string, req CreateVehicleRequest) {
       if err := s.authz.CanCreateVehicle(userID, orgID); err != nil {
           return ErrForbidden
       }
       // ...
   }
   ```

3. Use PostgreSQL Row-Level Security as defense-in-depth

---

### Validation

**Current State:**

**Backend:**
- Manual string validation in handlers (`if req.LocationID == ""`)
- No validation framework
- No VIN format validation (should be exactly 17 alphanumeric)
- No year range validation (should be 1900-current+1)
- Database CHECK constraints exist but limited

**Frontend:**
- HTML5 validation only (`required`, `min`, `max`, `pattern`)
- No custom error messages
- No cross-field validation
- No async validation (e.g., "is VIN already registered?")

**Issues:**

1. **Inconsistent Validation**
   - Frontend allows submission of invalid data
   - Backend may accept bad data
   - Database constraints are last line of defense

2. **Poor User Experience**
   - Generic browser validation messages
   - Errors shown one-at-a-time (HTML5 limitation)
   - No real-time validation feedback

3. **Security Risk**
   - No input sanitization
   - SQL injection prevented by GORM parameterization, but no XSS prevention
   - Mass assignment vulnerabilities (accepting all fields from request)

**Recommendations:**

**Backend:**
```go
// Use go-playground/validator
type CreateVehicleRequest struct {
    VIN  string `json:"vin" validate:"required,len=17,alphanum"`
    Make string `json:"make" validate:"required,max=100"`
    Year int    `json:"year" validate:"required,min=1900,max=2026"`
}
```

**Frontend:**
```typescript
// Use Zod + React Hook Form
const vehicleSchema = z.object({
  vin: z.string().length(17).regex(/^[A-Z0-9]{17}$/),
  make: z.string().min(1).max(100),
  year: z.number().int().min(1900).max(2026),
});
```

---

### Error Handling

**Current State:**

**Backend:**
- Generic HTTP errors: `http.Error(w, "Failed to fetch vehicles", 500)`
- No typed errors
- No error codes
- Stack traces not logged
- Errors don't distinguish between user error vs system error

**Frontend:**
- Catch blocks set generic error messages: `setError('Failed to fetch vehicles')`
- Some components show `err.response?.data` but inconsistent
- No error boundaries (runtime errors crash app)
- No retry logic

**Issues:**

1. **Poor Debugging**
   - Generic "Failed to create vehicle" - why did it fail?
   - No request IDs to trace errors
   - No structured logging

2. **Poor UX**
   - User sees "Internal server error" with no context
   - Cannot distinguish between:
     - Network error (retry possible)
     - Validation error (fix input)
     - Authorization error (permission denied)
     - System error (contact support)

3. **No Error Tracking**
   - No Sentry, Rollbar, or error monitoring
   - Errors only visible in logs

**Recommendations:**

**Backend:**
```go
// Typed errors
type AppError struct {
    Code    ErrorCode
    Message string
    Err     error
}

const (
    ErrNotFound     ErrorCode = "NOT_FOUND"
    ErrValidation   ErrorCode = "VALIDATION_ERROR"
    ErrForbidden    ErrorCode = "FORBIDDEN"
    ErrInternal     ErrorCode = "INTERNAL_ERROR"
)
```

**Frontend:**
```typescript
// Error boundaries
<ErrorBoundary fallback={<ErrorPage />}>
  <Routes />
</ErrorBoundary>

// Typed error handling
try {
  await vehiclesAPI.create(data);
} catch (err) {
  if (err.code === 'VALIDATION_ERROR') {
    // Show field errors
  } else if (err.code === 'FORBIDDEN') {
    // Show permission denied
  }
}
```

---

### Logging, Metrics, Config, Security

#### Logging

**Current:**
- Using `log.Println` (stdlib)
- No structured logging
- No log levels (debug, info, warn, error)
- No request IDs
- No user context

**Should Be:**
- Use zerolog or slog
- Structured JSON logs
- Request ID middleware
- Log levels configurable via environment

#### Metrics

**Current:**
- None visible
- No Prometheus, StatsD, or metrics endpoint

**Should Be:**
- Expose `/metrics` endpoint for Prometheus
- Track:
  - HTTP request duration (histogram)
  - Request count by endpoint/status (counter)
  - Database query duration (histogram)
  - Active connections (gauge)

#### Configuration

**Current:**
- Environment variables via `os.Getenv` with defaults
- No validation of required vars
- JWT secret hardcoded (!)

**Should Be:**
- Use Viper or env validation library
- Fail fast on missing required config
- Support `.env` files for local dev
- Never hardcode secrets

#### Security

**Current Issues:**

1. ✗ Hardcoded JWT secret
2. ✗ No rate limiting
3. ✗ No CSRF protection (SPA so less critical)
4. ✗ No input sanitization (XSS risk)
5. ✗ No SQL injection prevention beyond GORM parameterization
6. ✗ No security headers (CSP, X-Frame-Options, etc.)
7. ✗ No secrets management (AWS Secrets Manager, Vault)
8. ✗ CORS allows all origins in allowed list (should be env-based)

**Should Add:**

- Rate limiting middleware (e.g., `tollbooth`)
- Helmet.js equivalent for Go (security headers)
- Input sanitization (bluemonday for HTML)
- Move JWT secret to environment variable
- Add HTTPS redirect in production
- Implement session management (currently stateless JWT)

---

## 6. Prioritized Improvement Backlog

### Legend
- **Impact:** How much this improves quality, security, or maintainability
- **Effort:** Development time required (Low = 1-2 days, Medium = 3-5 days, High = 1-2 weeks)

| # | Area | Impact | Effort | Description |
|---|------|--------|--------|-------------|
| 1 | Backend | HIGH | MEDIUM | **Introduce Repository Layer** - Abstract database access behind interfaces for testability |
| 2 | Backend | HIGH | MEDIUM | **Introduce Service Layer** - Move business logic out of handlers into services |
| 3 | Backend | HIGH | LOW | **Implement Authorization Checks** - Enforce RBAC and org-level data isolation |
| 4 | Security | HIGH | LOW | **Move JWT Secret to Environment Variable** - Remove hardcoded secret |
| 5 | Frontend | HIGH | HIGH | **Introduce TypeScript** - Migrate all `.js` files to `.tsx` for type safety |
| 6 | Frontend | HIGH | LOW | **Create Custom Data Hooks** - Extract useVehicles, useVehicle, useLocations |
| 7 | Backend | HIGH | HIGH | **Add Dependency Injection** - Remove global database.DB, inject dependencies |
| 8 | Frontend | HIGH | MEDIUM | **Create Typed API Layer** - vehiclesAPI.getAll(), locationsAPI, etc. with TypeScript |
| 9 | Backend | MEDIUM | LOW | **Add go-playground/validator** - Declarative validation with struct tags |
| 10 | Frontend | MEDIUM | LOW | **Extract Reusable Components** - LoadingSpinner, ErrorAlert, VehicleTable |
| 11 | Backend | MEDIUM | MEDIUM | **Implement Structured Error Handling** - Typed domain errors with codes |
| 12 | Frontend | MEDIUM | MEDIUM | **Break Down VehicleForm** - Split 590-line component into smaller cards |
| 13 | Data | HIGH | MEDIUM | **Implement Audit Trail** - Populate audit_logs table via middleware |
| 14 | Backend | MEDIUM | MEDIUM | **Add Context Propagation** - Pass user_id, org_id, request_id via context |
| 15 | Frontend | MEDIUM | MEDIUM | **Add React Hook Form + Zod** - Better form handling and validation |
| 16 | Security | HIGH | MEDIUM | **Add Rate Limiting** - Prevent brute force and DoS attacks |
| 17 | Backend | MEDIUM | LOW | **Add Structured Logging** - Replace log.Println with zerolog/slog |
| 18 | Data | HIGH | LOW | **Fix Connection Pool Bug** - Use time.Duration for SetConnMaxLifetime |
| 19 | Frontend | MEDIUM | LOW | **Add Error Boundaries** - Prevent component errors from crashing app |
| 20 | Data | MEDIUM | MEDIUM | **Implement Soft Deletes** - Add deleted_at column, prevent data loss |
| 21 | Security | MEDIUM | LOW | **Add Security Headers Middleware** - CSP, X-Frame-Options, HSTS |
| 22 | Backend | MEDIUM | MEDIUM | **Add Background Job Processing** - Use asynq for bulk operations, emails |
| 23 | Frontend | MEDIUM | MEDIUM | **Add React Query** - Caching, deduplication, optimistic updates |
| 24 | Infra | MEDIUM | MEDIUM | **Add Metrics Endpoint** - Prometheus /metrics for observability |
| 25 | Data | MEDIUM | MEDIUM | **Implement Row-Level Security** - PostgreSQL RLS for defense-in-depth |
| 26 | Backend | LOW | LOW | **Remove Dead Code** - Clean up unused API functions |
| 27 | Frontend | LOW | LOW | **Remove Dead Code** - Delete src/api.js (unused) |
| 28 | Infra | MEDIUM | HIGH | **Add CI/CD Pipeline** - Automated testing and deployment |
| 29 | Data | LOW | HIGH | **Add Database Partitioning** - Partition vehicles table by org or date |
| 30 | Infra | MEDIUM | MEDIUM | **Add Database Backups** - Automated backup and restore procedures |

---

## Recommended First Steps (Next 2 Weeks)

### Week 1: Backend Foundation

**Day 1-2:** Repository Layer
- Create `internal/repositories/interfaces.go`
- Implement `VehicleRepository` with GORM
- Migrate one handler to use repository

**Day 3-4:** Service Layer
- Create `internal/services/vehicle_service.go`
- Move business logic from handler
- Add authorization checks

**Day 5:** Error Handling & Validation
- Add typed error structs
- Install go-playground/validator
- Update one handler with validation

### Week 2: Frontend Foundation

**Day 1-2:** TypeScript Setup
- Install TypeScript dependencies
- Create type definitions (`types/vehicle.ts`, `types/user.ts`)
- Migrate API layer to TypeScript

**Day 3:** Custom Hooks
- Create `hooks/useVehicles.ts`
- Create `hooks/useVehicle.ts`
- Migrate Vehicles.tsx to use hooks

**Day 4:** Reusable Components
- Extract LoadingSpinner
- Extract ErrorAlert
- Extract VehicleTable

**Day 5:** Testing & Documentation
- Write tests for new hooks
- Update README with new architecture
- Document migration progress

---

## Conclusion

FleetPass has a **solid foundation** with modern tech choices (Go, React, PostgreSQL, Docker) and clean separation between frontend and backend. However, it suffers from **classic v1 technical debt**:

- Backend violates clean architecture (no service/repository layers)
- Frontend lacks type safety (JavaScript instead of TypeScript)
- Authorization exists in the data model but is not enforced
- Validation is minimal and inconsistent
- No observability (logging, metrics, tracing)

The **highest ROI improvements** are:

1. Backend: Repository + Service layers (enables testing, maintainability)
2. Backend: Authorization enforcement (closes security gap)
3. Frontend: TypeScript migration (prevents bugs, improves DX)
4. Frontend: Custom hooks (eliminates duplication)
5. Security: Move secrets to environment, add rate limiting

**Estimated Effort for Full Refactor:** 6-8 weeks with 2 developers

**Recommended Approach:** Strangler Fig Pattern
- Refactor one feature end-to-end (e.g., Vehicle Management)
- Use as template for other features
- Old and new code coexist during migration
- Continuous delivery without big-bang rewrite

---

**Document Version:** 1.0
**Last Updated:** December 6, 2025
**Next Review:** After completing first 5 backlog items
