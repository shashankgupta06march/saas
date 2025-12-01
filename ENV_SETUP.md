# Environment Configuration Setup

## Quick Setup

### 1. Create .env File

```bash
cd /var/www/html/chatbot/backend
cp .env.example .env
```

### 2. Edit Configuration

```bash
nano .env
```

Update with your actual values:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=Admin@123
DB_NAME=chatbot_saas

# Server Configuration
SERVER_PORT=8081
ALLOWED_ORIGINS=http://localhost:3000

# OpenAI API Configuration (Get from https://platform.openai.com/api-keys)
OPENAI_API_KEY=sk-your-actual-key-here

# JWT Secret (Change in production)
JWT_SECRET=your-super-secret-jwt-key-change-this

# Environment
ENVIRONMENT=development
```

### 3. Run Seeder (Auto-loads .env)

```bash
cd /var/www/html/chatbot
./seed.sh go
```

The seeder now automatically reads from `backend/.env` file!

### 4. Start Backend (Auto-loads .env)

Update your start script or manually export:

```bash
cd /var/www/html/chatbot/backend
export $(grep -v '^#' .env | xargs)
go run cmd/api/main.go
```

## Configuration Details

### Database Settings

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | MySQL host | localhost |
| DB_PORT | MySQL port | 3306 |
| DB_USER | Database user | admin |
| DB_PASSWORD | Database password | Admin@123 |
| DB_NAME | Database name | chatbot_saas |

### Server Settings

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | API server port | 8081 |
| ALLOWED_ORIGINS | CORS allowed origins | http://localhost:3000 |

### OpenAI Settings

| Variable | Description | Required |
|----------|-------------|----------|
| OPENAI_API_KEY | OpenAI API key | Yes (for AI features) |

Get your key from: https://platform.openai.com/api-keys

### Security Settings

| Variable | Description | Required |
|----------|-------------|----------|
| JWT_SECRET | Secret for JWT tokens | Yes |

Generate a secure secret:
```bash
openssl rand -base64 32
```

## Production Deployment

### Security Checklist

- [ ] Change default database password
- [ ] Generate new JWT_SECRET
- [ ] Use real OpenAI API key
- [ ] Set ENVIRONMENT=production
- [ ] Use HTTPS for ALLOWED_ORIGINS
- [ ] Never commit .env file to git
- [ ] Use environment variables or secrets management

### Docker Deployment

Create `.env` file and use with docker-compose:

```yaml
services:
  backend:
    env_file:
      - .env
```

### Kubernetes Deployment

Create ConfigMap and Secret:

```bash
# Create secret for sensitive data
kubectl create secret generic chatbot-secrets \
  --from-literal=DB_PASSWORD=YourSecurePassword \
  --from-literal=OPENAI_API_KEY=sk-your-key \
  --from-literal=JWT_SECRET=your-jwt-secret

# Create configmap for non-sensitive data
kubectl create configmap chatbot-config \
  --from-literal=DB_HOST=mysql-service \
  --from-literal=DB_PORT=3306 \
  --from-literal=DB_NAME=chatbot_saas \
  --from-literal=SERVER_PORT=8081
```

## Troubleshooting

### .env Not Loading

**Problem**: Seeder or backend not reading .env file

**Solutions**:
```bash
# Check file exists
ls -la backend/.env

# Check file permissions
chmod 644 backend/.env

# Manually export
export $(grep -v '^#' backend/.env | xargs)
```

### Invalid .env Format

**Problem**: Errors when loading .env

**Solution**: Ensure no spaces around = and proper formatting:
```bash
# ✅ Correct
DB_USER=admin

# ❌ Wrong
DB_USER = admin
```

### Database Connection Failed

**Problem**: Cannot connect to database

**Solutions**:
```bash
# Verify credentials
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD -e "SELECT 1"

# Check MySQL is running
systemctl status mysql

# Check .env values
cat backend/.env | grep DB_
```

## Summary

✅ **Seeder now auto-loads from .env**  
✅ **No more hardcoded credentials**  
✅ **Production-ready configuration**  
✅ **Secure secrets management**

**Quick Start:**
```bash
# 1. Setup .env
cd /var/www/html/chatbot/backend
cp .env.example .env
nano .env  # Edit your values

# 2. Run seeder (auto-loads .env)
cd ..
./seed.sh go

# 3. Start backend
cd backend
export $(grep -v '^#' .env | xargs)
go run cmd/api/main.go
```

🚀 **That's it! Everything now uses .env configuration!**
