# Database Seeder - Quick Usage Guide

## Quick Start

### Run the Seeder

```bash
# Make executable (first time only)
chmod +x seed.sh

# Run Go seeder (recommended)
./seed.sh go

# Or run SQL seeder (faster)
./seed.sh sql
```

## What Gets Created

### 8 Organizations
- **Free Plan**: Test Company, Startup Inc, Small Business Co
- **Premium Plan**: Demo Corp, Tech Innovations, Medium Enterprise  
- **Enterprise Plan**: Enterprise Solutions Ltd, Global Services

### 11 Users (All passwords: `password123`)

| Email | Organization | Role | Plan |
|-------|--------------|------|------|
| admin@test.com | Test Company | admin | Free |
| user@test.com | Test Company | user | Free |
| demo@democorp.com | Demo Corp | admin | Premium |
| manager@democorp.com | Demo Corp | manager | Premium |
| admin@enterprise.com | Enterprise Solutions | admin | Enterprise |
| support@enterprise.com | Enterprise Solutions | support | Enterprise |
| contact@startup.com | Startup Inc | admin | Free |
| admin@techinnovations.com | Tech Innovations | admin | Premium |
| admin@globalservices.com | Global Services | admin | Enterprise |
| owner@smallbiz.com | Small Business Co | admin | Free |
| admin@mediumenterprise.com | Medium Enterprise | admin | Premium |

### 3 Chatbots
- Test Company Support Bot
- Demo Corp Support Bot
- Enterprise Solutions Ltd Support Bot

### 4 Knowledge Base Entries
- Sample content for testing chatbots

## Test Login

After seeding, login with:
```
Email: admin@test.com
Password: password123
```

## Environment Variables

```bash
export DB_USER=admin
export DB_PASSWORD=Admin@123
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=chatbot_saas
```

## Troubleshooting

### Connection Error
```bash
# Check MySQL is running
systemctl status mysql

# Test connection
mysql -u admin -pAdmin@123 -e "SELECT 1"
```

### Duplicate Data
The seeder will detect existing data and ask if you want to continue (Go seeder) or skip duplicates (SQL seeder).

### Go Not Found
```bash
# Use SQL seeder instead
./seed.sh sql
```

## Advanced Usage

### Customize Data
Edit `backend/cmd/seeder/main.go` to add/modify:
- Organizations
- Users
- Passwords
- Roles

### Run from Code
```bash
cd backend
export DB_USER=admin DB_PASSWORD=Admin@123
go run cmd/seeder/main.go
```

### Run SQL Directly
```bash
mysql -u admin -pAdmin@123 chatbot_saas < backend/migrations/003_seed_data.sql
```

## Complete Documentation

See [SEEDER_DOCUMENTATION.md](SEEDER_DOCUMENTATION.md) for:
- Detailed implementation
- Production deployment
- Security considerations
- CI/CD integration
- Customization guide

## Summary

**Fastest Method:**
```bash
./seed.sh sql
```

**Interactive Method:**
```bash
./seed.sh go
```

**Both create complete demo data for testing all features!**

