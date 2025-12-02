# Deployment Guide

This guide covers deploying FleetPass to production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Deployment Options](#deployment-options)
- [Production Checklist](#production-checklist)
- [Monitoring](#monitoring)
- [Backup & Recovery](#backup--recovery)

## Prerequisites

- Docker and Docker Compose installed
- PostgreSQL 16+ database
- Domain name configured
- SSL/TLS certificates
- Secrets management solution

## Environment Configuration

### 1. Copy Production Environment Template

```bash
cp .env.production.example .env.production
```

### 2. Configure Required Variables

Edit `.env.production` with your production values:

```bash
# Database - Use managed PostgreSQL service in production
POSTGRES_DB=fleetpass_production
POSTGRES_USER=fleetpass_prod_user
POSTGRES_PASSWORD=<SECURE_PASSWORD>
DB_HOST=<your-db-host>
DB_PORT=5432
DB_SSLMODE=require

# Application
JWT_SECRET=<GENERATE_SECURE_32+_CHAR_STRING>
API_PORT=8080

# Frontend
REACT_APP_API_URL=https://api.yourfleetpass.com

# Security
CORS_ALLOWED_ORIGINS=https://yourfleetpass.com,https://www.yourfleetpass.com
```

### 3. Generate Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Generate database password
openssl rand -base64 24
```

## Deployment Options

### Option 1: Docker Compose (Small Scale)

**Best for**: Single server, low-moderate traffic

1. **Prepare production docker-compose file:**

```yaml
# docker-compose.production.yml
version: '3.8'

services:
  api:
    image: yourusername/fleetpass-api:latest
    restart: always
    environment:
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${POSTGRES_DB}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_SSLMODE=require
      - JWT_SECRET=${JWT_SECRET}
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: yourusername/fleetpass-frontend:latest
    restart: always
    environment:
      - REACT_APP_API_URL=${REACT_APP_API_URL}
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - api
```

2. **Deploy:**

```bash
docker-compose -f docker-compose.production.yml up -d
```

### Option 2: Kubernetes (Large Scale)

**Best for**: Multi-server, high traffic, auto-scaling

See `k8s/` directory for Kubernetes manifests.

### Option 3: Cloud Platforms

#### AWS ECS/Fargate

1. Build and push images to ECR
2. Create ECS task definitions
3. Configure load balancer
4. Set up RDS for PostgreSQL

#### Google Cloud Run

```bash
# Build and deploy API
gcloud builds submit --tag gcr.io/PROJECT_ID/fleetpass-api
gcloud run deploy fleetpass-api \
  --image gcr.io/PROJECT_ID/fleetpass-api \
  --platform managed \
  --region us-central1 \
  --set-env-vars DB_HOST=...,JWT_SECRET=...

# Build and deploy Frontend
gcloud builds submit --tag gcr.io/PROJECT_ID/fleetpass-frontend ./frontend
gcloud run deploy fleetpass-frontend \
  --image gcr.io/PROJECT_ID/fleetpass-frontend \
  --platform managed \
  --region us-central1
```

#### Azure Container Instances

```bash
# Create resource group
az group create --name fleetpass-rg --location eastus

# Deploy containers
az container create \
  --resource-group fleetpass-rg \
  --name fleetpass-api \
  --image yourusername/fleetpass-api:latest \
  --environment-variables DB_HOST=... JWT_SECRET=... \
  --ports 8080
```

## Production Checklist

### Security

- [ ] Change all default passwords
- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable SSL/TLS for all connections
- [ ] Configure CORS properly
- [ ] Use environment variables for secrets
- [ ] Enable database SSL mode
- [ ] Set up firewall rules
- [ ] Regular security updates
- [ ] Enable rate limiting
- [ ] Set up WAF (Web Application Firewall)

### Database

- [ ] Use managed database service (AWS RDS, Google Cloud SQL, etc.)
- [ ] Enable automated backups
- [ ] Set up point-in-time recovery
- [ ] Configure connection pooling
- [ ] Enable SSL connections
- [ ] Set up read replicas (if needed)
- [ ] Monitor database performance
- [ ] Regular vacuum/analyze operations

### Application

- [ ] Set production environment variables
- [ ] Configure proper logging
- [ ] Set up error tracking (Sentry, etc.)
- [ ] Configure monitoring (Prometheus, Datadog, etc.)
- [ ] Set up health checks
- [ ] Configure auto-scaling
- [ ] Set resource limits
- [ ] Enable compression
- [ ] Configure caching headers

### Networking

- [ ] Configure DNS
- [ ] Set up SSL certificates (Let's Encrypt, etc.)
- [ ] Configure load balancer
- [ ] Set up CDN for frontend assets
- [ ] Configure reverse proxy (nginx, etc.)
- [ ] Enable HTTP/2
- [ ] Set up GZIP compression

### CI/CD

- [ ] Set up GitHub Actions secrets
- [ ] Configure Docker Hub credentials
- [ ] Set up deployment pipeline
- [ ] Configure deployment approvals
- [ ] Set up rollback procedures
- [ ] Enable deployment notifications

## Monitoring

### Application Metrics

Monitor these key metrics:

- Request rate and latency
- Error rates
- CPU and memory usage
- Database connection pool status
- Active user sessions

### Recommended Tools

- **Application Performance**: New Relic, Datadog, or Application Insights
- **Error Tracking**: Sentry or Rollbar
- **Uptime Monitoring**: Pingdom or UptimeRobot
- **Log Aggregation**: ELK Stack, Splunk, or CloudWatch

### Health Check Endpoints

```bash
# API Health
curl https://api.yourfleetpass.com/health

# Database connectivity
curl https://api.yourfleetpass.com/health/db
```

## Backup & Recovery

### Database Backups

```bash
# Automated daily backups (managed service handles this)
# Or manual backup:
pg_dump -h $DB_HOST -U $POSTGRES_USER -d $POSTGRES_DB > backup_$(date +%Y%m%d).sql
```

### Disaster Recovery

1. **Regular Backups**: Daily automated backups with 30-day retention
2. **Point-in-Time Recovery**: Enabled for last 7 days
3. **Multi-Region**: Database replicas in different regions
4. **Documentation**: Recovery procedures documented and tested

### Recovery Procedure

```bash
# Restore from backup
psql -h $DB_HOST -U $POSTGRES_USER -d $POSTGRES_DB < backup_20231201.sql

# Verify data integrity
# Test critical features
# Monitor for errors
```

## SSL/TLS Setup

### Using Let's Encrypt

```bash
# Install certbot
sudo apt-get install certbot

# Generate certificate
sudo certbot certonly --standalone -d yourfleetpass.com -d www.yourfleetpass.com

# Auto-renewal
sudo certbot renew --dry-run
```

### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name yourfleetpass.com;

    ssl_certificate /etc/letsencrypt/live/yourfleetpass.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourfleetpass.com/privkey.pem;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://frontend:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api {
        proxy_pass http://api:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Scaling

### Horizontal Scaling

- Use load balancer (AWS ALB, nginx, etc.)
- Scale API and frontend containers independently
- Database read replicas for read-heavy operations

### Vertical Scaling

- Monitor resource usage
- Increase container resources as needed
- Optimize database queries

## Troubleshooting

### Common Issues

**Database Connection Errors**
```bash
# Check database connectivity
docker exec -it fleetpass-api ping $DB_HOST
docker logs fleetpass-api | grep -i database
```

**High Memory Usage**
```bash
# Check container stats
docker stats
# Increase memory limits if needed
```

**SSL Certificate Issues**
```bash
# Verify certificate
openssl s_client -connect yourfleetpass.com:443 -servername yourfleetpass.com
```

## Rollback Procedures

If deployment fails:

1. **Docker Compose:**
   ```bash
   docker-compose -f docker-compose.production.yml down
   # Deploy previous version
   docker-compose -f docker-compose.production.yml up -d
   ```

2. **Kubernetes:**
   ```bash
   kubectl rollout undo deployment/fleetpass-api
   kubectl rollout undo deployment/fleetpass-frontend
   ```

## Support

For deployment issues:
- Check logs: `docker-compose logs -f`
- Review monitoring dashboards
- Check GitHub Issues
- Contact support team
