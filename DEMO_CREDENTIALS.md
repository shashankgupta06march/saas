# Demo Credentials

## System Access

### Backend API
- **Base URL**: `http://localhost:8081/api`
- **Widget Endpoint**: `http://localhost:8081/widget.js`
- **Status**: ✅ Running

### Database
- **Host**: localhost
- **Port**: 3306
- **Database**: chatbot_saas
- **Username**: admin
- **Password**: Admin@123

---

## Demo User Accounts

> **Note**: To seed all demo accounts automatically, run: `./seed.sh go` or `./seed.sh sql`  
> See [SEEDER_DOCUMENTATION.md](SEEDER_DOCUMENTATION.md) for details.

### 1. Organization Administrator
**Organization**: Test Company  
**Email**: `admin@test.com`  
**Password**: `password123`  
**Role**: admin  
**Plan**: Free  
**Organization ID**: 1  
**API Key**: `06efecc90c6069f61c803ca415cdb01d5a7fe3c192f10a6bd70031e3bda484e6`

**Capabilities**:
- Create/manage chatbots
- Upload knowledge base
- Customize widget appearance
- View analytics
- Manage organization settings

### 2. Premium Organization Admin
**Organization**: Demo Corp  
**Email**: `demo@democorp.com`  
**Password**: `password123`  
**Role**: admin  
**Plan**: Premium  

### 3. Enterprise Organization Admin
**Organization**: Enterprise Solutions Ltd  
**Email**: `admin@enterprise.com`  
**Password**: `password123`  
**Role**: admin  
**Plan**: Enterprise  

### 4. Additional Test Users
**Organization**: Test Company  
**Email**: `user@test.com`  
**Password**: `password123`  
**Role**: user  

**Organization**: Demo Corp  
**Email**: `manager@democorp.com`  
**Password**: `password123`  
**Role**: manager  

**Organization**: Enterprise Solutions Ltd  
**Email**: `support@enterprise.com`  
**Password**: `password123`  
**Role**: support  

> **All Seeded Credentials**: Run the seeder to create 8 organizations with 11 users. All passwords are `password123` for demo purposes.

---

## API Testing Examples

### Login
```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }'
```

### Create a Chatbot (Requires Token)
```bash
curl -X POST http://localhost:8081/api/chatbots \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Support Bot",
    "description": "Customer support chatbot"
  }'
```

### List Chatbots
```bash
curl -X GET http://localhost:8081/api/chatbots \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Add Knowledge Base Entry
```bash
curl -X POST http://localhost:8081/api/knowledge \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "chatbot_id": 1,
    "title": "Company Info",
    "content": "We are a technology company providing AI solutions.",
    "content_type": "text"
  }'
```

---

## Creating Additional Demo Users

### Option 1: Use Automated Seeder (Recommended)

```bash
# Seed all demo data at once
./seed.sh go

# Or use SQL seeder for faster execution
./seed.sh sql
```

This creates:
- 8 Organizations (Free, Premium, Enterprise plans)
- 11 Users with various roles
- 3 Sample Chatbots
- 4 Knowledge Base Entries

See [SEEDER_DOCUMENTATION.md](SEEDER_DOCUMENTATION.md) for complete details.

### Option 2: Manual Registration

#### Create Second Organization
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Demo Corp",
    "email": "demo@democorp.com",
    "password": "demo123456"
  }'
```

#### Create Third Organization
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Sample Industries",
    "email": "contact@sampleind.com",
    "password": "sample123"
  }'
```

---

## Frontend Demo Access (When Available)

Once the frontend is running on `http://localhost:3000`:

### Login Page
Navigate to: `http://localhost:3000/login`

Use any of the credentials above to login.

---

## Widget Integration Demo

To test the embeddable widget on your website:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Widget Demo</title>
</head>
<body>
    <h1>My Website</h1>
    <p>This is a demo page with the chatbot widget.</p>
    
    <!-- Chatbot Widget -->
    <script 
        src="http://localhost:8081/widget.js" 
        data-chatbot-id="1"
        data-api-url="http://localhost:8081/api">
    </script>
</body>
</html>
```

Or open the provided demo file:
```bash
# Update the chatbot ID in the file first
open /var/www/html/chatbot/widget/demo.html
```

---

## User Roles Explained

### Admin Role
- **Full Access**: Yes
- **Can Create Chatbots**: Yes
- **Can Manage Knowledge Base**: Yes
- **Can View Analytics**: Yes
- **Can Customize Widgets**: Yes
- **Can Manage Organization**: Yes

**Note**: Currently, the system implements a single "admin" role per organization. All users in an organization have admin privileges for that organization's resources.

---

## Multi-Tenancy Architecture

Each organization is completely isolated:
- Organization 1 cannot access Organization 2's data
- Each organization has unique API keys
- Chatbots, knowledge base, and conversations are organization-scoped

---

## Database Direct Access

For development/testing, you can access the database directly:

```bash
mysql -u admin -pAdmin@123 chatbot_saas
```

### Useful Queries

**View all organizations:**
```sql
SELECT * FROM organizations;
```

**View all users:**
```sql
SELECT id, email, organization_id, role, created_at FROM users;
```

**View all chatbots:**
```sql
SELECT c.*, o.name as org_name 
FROM chatbots c 
JOIN organizations o ON c.organization_id = o.id;
```

**View knowledge base entries:**
```sql
SELECT id, chatbot_id, title, LEFT(content, 50) as content_preview, created_at 
FROM knowledge_base 
ORDER BY created_at DESC;
```

**View recent conversations:**
```sql
SELECT c.*, ch.name as chatbot_name 
FROM conversations c 
JOIN chatbots ch ON c.chatbot_id = ch.id 
ORDER BY started_at DESC 
LIMIT 10;
```

---

## OpenAI API Configuration

To enable AI features, you need a valid OpenAI API key.

### Get Your API Key
1. Go to: https://platform.openai.com/api-keys
2. Create a new API key
3. Replace `sk-placeholder` in the backend startup command

### Restart Backend with Real API Key
```bash
cd /var/www/html/chatbot/backend
export PATH=$HOME/go/bin:$PATH
SERVER_PORT=8081 \
DB_USER=admin \
DB_PASSWORD=Admin@123 \
DB_NAME=chatbot_saas \
OPENAI_API_KEY=sk-your-real-key-here \
ALLOWED_ORIGINS=http://localhost:3000 \
go run cmd/api/main.go
```

---

## Security Notes

⚠️ **For Development Only**

These credentials are for development and testing purposes only. In production:

1. Change all default passwords
2. Use strong, unique passwords
3. Enable HTTPS
4. Restrict database access
5. Use environment variables for sensitive data
6. Implement rate limiting
7. Enable proper CORS policies
8. Regular security audits

---

## Support

For issues or questions:
- Check logs in the backend terminal
- Review API documentation in README.md
- Check database for data integrity
- Verify all services are running

---

**Last Updated**: October 14, 2025  
**System Status**: Backend Running on Port 8081 ✅  
**Database**: Connected ✅



