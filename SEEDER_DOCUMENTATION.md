# Database Seeder Documentation

This document explains how to use the database seeder to populate your Chatbot SaaS application with demo data.

## Overview

The seeder creates:
- **8 Organizations** with different plan types (Free, Premium, Enterprise)
- **11 Users** with various roles (admin, user, manager, support)
- **3 Sample Chatbots** with configured settings
- **4 Knowledge Base Entries** for testing

## Available Seeder Methods

### Method 1: Go Seeder (Recommended)

The Go seeder provides interactive feedback and proper password hashing.

**Run the seeder:**
```bash
./seed.sh go
```

Or directly:
```bash
cd backend
export DB_USER=admin
export DB_PASSWORD=Admin@123
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=chatbot_saas
go run cmd/seeder/main.go
```

**Features:**
- ✅ Proper bcrypt password hashing
- ✅ Interactive prompts if data exists
- ✅ Detailed logging of each operation
- ✅ Automatic API key generation
- ✅ Comprehensive error handling
- ✅ Summary report with credentials

### Method 2: SQL Seeder (Quick Deploy)

The SQL seeder is faster for production deployments.

**Run the seeder:**
```bash
./seed.sh sql
```

Or directly:
```bash
mysql -u admin -pAdmin@123 chatbot_saas < backend/migrations/003_seed_data.sql
```

**Features:**
- ✅ Fast execution
- ✅ Idempotent (uses INSERT IGNORE)
- ✅ Can be integrated into deployment pipelines
- ✅ No additional dependencies

## Make Scripts Executable

```bash
chmod +x seed.sh
chmod +x start-backend.sh
chmod +x start-frontend.sh
chmod +x stop-backend.sh
```

## Seeded Data Details

### Organizations

| ID | Name | Plan Type | Status | API Key (Prefix) |
|----|------|-----------|--------|------------------|
| 1 | Test Company | free | active | 06efecc90c... |
| 2 | Demo Corp | premium | active | a1b2c3d4e5... |
| 3 | Enterprise Solutions Ltd | enterprise | active | e1f2g3h4i5... |
| 4 | Startup Inc | free | active | f1g2h3i4j5... |
| 5 | Tech Innovations | premium | active | g1h2i3j4k5... |
| 6 | Global Services | enterprise | active | h1i2j3k4l5... |
| 7 | Small Business Co | free | trial | i1j2k3l4m5... |
| 8 | Medium Enterprise | premium | active | j1k2l3m4n5... |

### Users (All passwords: `password123`)

| Organization | Email | Role | Access Level |
|--------------|-------|------|--------------|
| Test Company | admin@test.com | admin | Full access |
| Test Company | user@test.com | user | Limited access |
| Demo Corp | demo@democorp.com | admin | Full access |
| Demo Corp | manager@democorp.com | manager | Management |
| Enterprise Solutions | admin@enterprise.com | admin | Full access |
| Enterprise Solutions | support@enterprise.com | support | Support only |
| Startup Inc | contact@startup.com | admin | Full access |
| Tech Innovations | admin@techinnovations.com | admin | Full access |
| Global Services | admin@globalservices.com | admin | Full access |
| Small Business Co | owner@smallbiz.com | admin | Full access |
| Medium Enterprise | admin@mediumenterprise.com | admin | Full access |

### Sample Chatbots

1. **Test Company Support Bot**
   - Theme: Blue (#007bff)
   - Position: Bottom right
   - Size: Medium

2. **Demo Corp Support Bot**
   - Theme: Green (#28a745)
   - Position: Bottom right
   - Size: Medium

3. **Enterprise Solutions Ltd Support Bot**
   - Theme: Purple (#6f42c1)
   - Position: Bottom right
   - Size: Large

### Knowledge Base Entries

- Company Information (Test Company)
- Product Features (Test Company)
- Pricing Plans (Demo Corp)
- Support Hours (Enterprise Solutions)

## Environment Variables

Configure these before running the seeder:

```bash
export DB_USER=admin              # Database username
export DB_PASSWORD=Admin@123      # Database password
export DB_HOST=localhost          # Database host
export DB_PORT=3306               # Database port
export DB_NAME=chatbot_saas       # Database name
```

## Testing the Seeded Data

### 1. Login to Frontend

```bash
# Start the frontend
./start-frontend.sh

# Navigate to http://localhost:3000/login
# Use any of the credentials above
```

### 2. Test API Endpoints

```bash
# Login
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }'

# Get Chatbots (use token from login)
curl -X GET http://localhost:8081/api/chatbots \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 3. Check Database

```bash
mysql -u admin -pAdmin@123 chatbot_saas

# View organizations
SELECT id, name, plan_type, status FROM organizations;

# View users
SELECT id, email, organization_id, role FROM users;

# View chatbots
SELECT c.id, c.name, o.name as org_name 
FROM chatbots c 
JOIN organizations o ON c.organization_id = o.id;
```

## Integration with Deployment

### Docker Compose

Add to your `docker-compose.yml`:

```yaml
services:
  seeder:
    build: ./backend
    command: go run cmd/seeder/main.go
    environment:
      - DB_USER=admin
      - DB_PASSWORD=Admin@123
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=chatbot_saas
    depends_on:
      - mysql
```

### Kubernetes

Create a Job:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: chatbot-seeder
spec:
  template:
    spec:
      containers:
      - name: seeder
        image: chatbot-saas-backend:latest
        command: ["go", "run", "cmd/seeder/main.go"]
        env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
      restartPolicy: Never
```

### CI/CD Pipeline

Add to your deployment script:

```bash
#!/bin/bash

# Deploy database migrations
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < migrations/001_initial_schema.sql
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < migrations/002_enhanced_knowledge_base.sql

# Seed initial data (only on first deployment)
if [ "$SEED_DATA" = "true" ]; then
  mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < migrations/003_seed_data.sql
fi
```

## Resetting the Database

To start fresh:

```bash
# Drop and recreate database
mysql -u admin -pAdmin@123 -e "DROP DATABASE IF EXISTS chatbot_saas; CREATE DATABASE chatbot_saas;"

# Run migrations
mysql -u admin -pAdmin@123 chatbot_saas < backend/migrations/001_initial_schema.sql
mysql -u admin -pAdmin@123 chatbot_saas < backend/migrations/002_enhanced_knowledge_base.sql

# Run seeder
./seed.sh go
```

## Production Considerations

### Security

⚠️ **Important**: The seeded data is for development/demo purposes only.

For production:

1. **Change all passwords**
   ```bash
   # Use strong, unique passwords
   UPDATE users SET password_hash = ? WHERE email = ?;
   ```

2. **Regenerate API keys**
   ```bash
   # Generate secure API keys
   UPDATE organizations SET api_key = ? WHERE id = ?;
   ```

3. **Remove test accounts**
   ```bash
   # Delete demo accounts
   DELETE FROM users WHERE email LIKE '%@test.com';
   ```

4. **Use environment variables**
   ```bash
   # Never hardcode credentials
   export DB_PASSWORD=$(cat /run/secrets/db_password)
   ```

### Data Privacy

- Don't seed production databases with fake user data
- Use anonymized data if testing in production-like environments
- Follow GDPR/privacy regulations for test data

### Performance

- For large datasets, use batch inserts
- Consider using database transactions
- Run seeder during off-peak hours if seeding production

## Troubleshooting

### Error: "Failed to connect to database"

**Solution**: Check database credentials and ensure MySQL is running
```bash
# Test connection
mysql -h localhost -u admin -pAdmin@123 -e "SELECT 1"

# Check if MySQL is running
systemctl status mysql
```

### Error: "Duplicate entry for key 'email'"

**Solution**: Data already exists. Either:
- Choose 'n' when prompted by Go seeder
- Manually delete existing data
- Use different email addresses

### Error: "go: command not found"

**Solution**: Install Go or use SQL seeder
```bash
# Use SQL seeder instead
./seed.sh sql

# Or install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Error: "Access denied for user"

**Solution**: Update database credentials
```bash
export DB_USER=your_user
export DB_PASSWORD=your_password
./seed.sh go
```

## Customization

### Adding More Organizations

Edit `backend/cmd/seeder/main.go`:

```go
organizations := []Organization{
    // ... existing organizations
    {Name: "Your Company", PlanType: "premium", Status: "active"},
}

users := []User{
    // ... existing users
    {Email: "admin@yourcompany.com", Password: "secure123", Role: "admin"},
}
```

### Changing Default Password

Update the password in the seeder file:

```go
// Go seeder
{Email: "admin@test.com", Password: "YourSecurePassword123", Role: "admin"}

// SQL seeder - generate hash first
// In Go:
hash, _ := auth.HashPassword("YourSecurePassword123")
// Then update SQL file with the hash
```

### Adding Custom Roles

1. Update the seeder to include new roles
2. Modify your authorization middleware to handle new roles
3. Update frontend role checks

## Support

For issues or questions:
- Check logs: `tail -f backend/server.log`
- Review database: `mysql -u admin -pAdmin@123 chatbot_saas`
- Check seeder output for specific errors

## Summary

The seeder provides two flexible methods to populate your database:

1. **Go Seeder**: Best for development with interactive feedback
2. **SQL Seeder**: Best for automated deployments

Both methods create comprehensive demo data covering all user types and organization tiers, making it easy to test and demonstrate your Chatbot SaaS application.

**Quick Start:**
```bash
# Make executable
chmod +x seed.sh

# Run Go seeder
./seed.sh go

# Or run SQL seeder
./seed.sh sql

# Login with any seeded credentials
Email: admin@test.com
Password: password123
```

Happy coding! 🚀

