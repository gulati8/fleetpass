# Testing Authentication System

## Setup

### 1. Rebuild Docker Containers

Since we've made significant backend changes, rebuild:

```bash
docker-compose down
docker-compose up -d --build api
```

### 2. Watch the Logs

```bash
docker-compose logs -f api
```

You should see:
- Database migrations running
- Permissions being seeded
- Roles being created
- Super admin user created with credentials:
  - Email: `admin@fleetpass.com`
  - Password: `Admin123!`

---

## Test Scenarios

### Scenario 1: Login with Super Admin

**Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fleetpass.com",
    "password": "Admin123!"
  }'
```

**Expected Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@fleetpass.com",
    "first_name": "System",
    "last_name": "Administrator",
    "email_verified": true,
    "is_active": true,
    "roles": ["super_admin"],
    "permissions": ["vehicles.create", "vehicles.read", ...],
    "organization_id": null
  }
}
```

---

### Scenario 2: User Registration

**Request:**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

**Expected Response:**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "user_id": "uuid"
}
```

**Check Logs:**
You should see an email printed to console with verification token.

---

### Scenario 3: Email Verification

Copy the verification token from the console logs, then:

**Request:**
```bash
curl -X POST http://localhost:8080/api/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "token": "PASTE_TOKEN_HERE"
  }'
```

**Expected Response:**
```json
{
  "message": "Email verified successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "email_verified": true,
    "is_active": true,
    "roles": ["customer"],
    "permissions": ["vehicles.read", "rentals.create", "rentals.read"],
    "organization_id": null
  }
}
```

---

### Scenario 4: Login with Verified User

**Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Expected Response:**
User profile with JWT token and customer role.

---

### Scenario 5: Get Profile (Protected Endpoint)

Use the token from login:

**Request:**
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Expected Response:**
User profile with roles and permissions.

---

### Scenario 6: Forgot Password

**Request:**
```bash
curl -X POST http://localhost:8080/api/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

**Expected Response:**
```json
{
  "message": "If an account exists with this email, you will receive a password reset link."
}
```

**Check Logs:**
You should see a password reset email with token.

---

### Scenario 7: Reset Password

Copy reset token from logs:

**Request:**
```bash
curl -X POST http://localhost:8080/api/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "PASTE_RESET_TOKEN_HERE",
    "new_password": "NewSecure123!"
  }'
```

**Expected Response:**
```json
{
  "message": "Password reset successful. You can now log in with your new password."
}
```

---

### Scenario 8: Login with New Password

**Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "NewSecure123!"
  }'
```

**Expected Response:**
Successful login with token.

---

## Error Cases to Test

### 1. Weak Password

**Request:**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "weak",
    "first_name": "Test",
    "last_name": "User"
  }'
```

**Expected:** Error about password requirements

---

### 2. Duplicate Email

Try registering with `admin@fleetpass.com`:

**Expected:** "User with this email already exists"

---

### 3. Login Without Email Verification

1. Register a new user
2. Try to login WITHOUT verifying email

**Expected:** "Please verify your email address before logging in"

---

### 4. Invalid Credentials

**Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fleetpass.com",
    "password": "wrongpassword"
  }'
```

**Expected:** 401 Unauthorized - "Invalid credentials"

---

### 5. Access Protected Route Without Token

**Request:**
```bash
curl http://localhost:8080/api/profile
```

**Expected:** 401 Unauthorized

---

## Database Verification

Connect to the database:

```bash
docker exec -it fleetpass-db psql -U fleetpass_user -d fleetpass
```

### Check Users

```sql
SELECT id, email, first_name, last_name, email_verified, is_active FROM users;
```

### Check Roles

```sql
SELECT * FROM roles;
```

### Check Permissions

```sql
SELECT * FROM permissions;
```

### Check User Roles

```sql
SELECT u.email, r.name as role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id;
```

### Check Role Permissions

```sql
SELECT r.name as role_name, p.name as permission_name
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.name, p.name;
```

---

## Success Criteria

✅ **Phase 1 Complete if:**
1. Super admin can login
2. New users can register
3. Email verification works (token in console logs)
4. Password reset works (token in console logs)
5. JWT tokens include roles and permissions
6. Database has all permissions, roles, and user_roles correctly set up
7. Protected endpoints require authentication

---

## Common Issues

### Issue: Database connection error
**Solution:** Ensure PostgreSQL container is running:
```bash
docker-compose ps
docker-compose logs db
```

### Issue: No seed data
**Solution:** Check API logs for seed script output:
```bash
docker-compose logs api | grep -i seed
```

### Issue: Can't login with super admin
**Solution:** Verify super admin was created:
```bash
docker-compose logs api | grep "Super admin"
```

Should show: `Email: admin@fleetpass.com` and `Password: Admin123!`

---

## Next Steps After Testing

Once all tests pass:
1. ✅ Backend authentication is working
2. Build frontend registration/login pages
3. Build role/permission guards
4. Add user management endpoints for admins

---

**Ready to test?** Run `docker-compose up -d --build api` and start with Scenario 1!
