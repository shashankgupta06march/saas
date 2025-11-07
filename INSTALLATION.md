# Installation Guide

This guide will help you install all prerequisites and set up the chatbot platform.

## Prerequisites Installation

### 1. Install Go (1.21 or higher)

**Ubuntu/Debian:**
```bash
# Remove old Go installation if exists
sudo rm -rf /usr/local/go

# Download Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

**Alternative (using snap):**
```bash
sudo snap install go --classic
```

### 2. Install Node.js (18 or higher)

**Using NodeSource repository:**
```bash
# Install Node.js 18.x
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

**Alternative (using nvm):**
```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc

# Install Node.js
nvm install 18
nvm use 18
```

### 3. MySQL Setup

MySQL should already be installed. Verify:
```bash
mysql --version
sudo systemctl status mysql
```

If not installed:
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql
```

### 4. Verify All Prerequisites

Run this command to check all prerequisites:
```bash
echo "Go: $(go version 2>/dev/null || echo 'NOT INSTALLED')"
echo "Node: $(node --version 2>/dev/null || echo 'NOT INSTALLED')"
echo "NPM: $(npm --version 2>/dev/null || echo 'NOT INSTALLED')"
echo "MySQL: $(mysql --version 2>/dev/null || echo 'NOT INSTALLED')"
```

## Quick Start

Once all prerequisites are installed:

### 1. Run Automated Setup

```bash
cd /var/www/html/chatbot
chmod +x setup.sh
./setup.sh
```

This script will:
- Check all prerequisites
- Create MySQL database and tables
- Install backend Go dependencies
- Install frontend npm packages
- Create .env file if needed

### 2. Configure OpenAI API Key

Edit the backend configuration:
```bash
nano /var/www/html/chatbot/backend/.env
```

Add your OpenAI API key:
```
OPENAI_API_KEY=sk-your-actual-api-key-here
```

Get your API key from: https://platform.openai.com/api-keys

### 3. Start Backend

Terminal 1:
```bash
cd /var/www/html/chatbot
./start-backend.sh
```

Or manually:
```bash
cd /var/www/html/chatbot/backend
go run cmd/api/main.go
```

You should see:
```
Server starting on port 8080
Database connection established successfully
```

### 4. Start Frontend

Terminal 2:
```bash
cd /var/www/html/chatbot
./start-frontend.sh
```

Or manually:
```bash
cd /var/www/html/chatbot/frontend
npm run dev
```

You should see:
```
Local: http://localhost:3000
```

### 5. Access the Application

Open your browser and navigate to:
- **Admin Dashboard**: http://localhost:3000
- **API**: http://localhost:8080

### 6. Create Your First Account

1. Click "Register"
2. Fill in:
   - Organization Name: Your company name
   - Email: Your email
   - Password: At least 6 characters
3. Click "Register"

You'll be automatically logged in to the dashboard.

## Manual Setup (If Automated Setup Fails)

### Step 1: Database Setup

```bash
# Login to MySQL
mysql -u admin -pAdmin@123

# Run migration
mysql -u admin -pAdmin@123 < /var/www/html/chatbot/backend/migrations/001_initial_schema.sql

# Verify database created
mysql -u admin -pAdmin@123 -e "USE chatbot_saas; SHOW TABLES;"
```

### Step 2: Backend Setup

```bash
cd /var/www/html/chatbot/backend

# Copy environment file
cp .env.example .env

# Edit .env and add your OpenAI API key
nano .env

# Install Go dependencies
go mod download
go mod tidy

# Run backend
go run cmd/api/main.go
```

### Step 3: Frontend Setup

```bash
cd /var/www/html/chatbot/frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## Troubleshooting

### Go Not Found After Installation

Add Go to your PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### MySQL Connection Error

Check MySQL credentials:
```bash
mysql -u admin -pAdmin@123 -e "SELECT 1;"
```

If connection fails, verify user exists:
```bash
sudo mysql
CREATE USER 'admin'@'localhost' IDENTIFIED BY 'Admin@123';
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### Port Already in Use

If port 8080 or 3000 is already in use:

For backend (port 8080):
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill the process
sudo kill -9 <PID>
```

Or change port in `backend/.env`:
```
SERVER_PORT=8081
```

For frontend (port 3000):
```bash
# Kill process on port 3000
sudo lsof -i :3000 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

Or set different port:
```bash
PORT=3001 npm run dev
```

### OpenAI API Errors

Common issues:
1. **Invalid API key**: Verify key in backend/.env
2. **No credits**: Check billing at https://platform.openai.com/account/billing
3. **Rate limit**: Wait a moment and try again

Test API key:
```bash
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Frontend Build Errors

Clear cache and reinstall:
```bash
cd /var/www/html/chatbot/frontend
rm -rf node_modules package-lock.json
npm install
```

### Database Migration Errors

Drop and recreate database:
```bash
mysql -u admin -pAdmin@123 -e "DROP DATABASE IF EXISTS chatbot_saas;"
mysql -u admin -pAdmin@123 < /var/www/html/chatbot/backend/migrations/001_initial_schema.sql
```

## Testing the Installation

### 1. Test Backend API

```bash
# Health check (should return 404 but confirms server is running)
curl http://localhost:8080/

# Test registration
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Test Org",
    "email": "test@example.com",
    "password": "test123"
  }'
```

### 2. Test Frontend

Open http://localhost:3000 in your browser. You should see the login page.

### 3. Test Widget

1. Create a chatbot in the admin dashboard
2. Add some knowledge base entries
3. Open `/var/www/html/chatbot/widget/demo.html` in browser
4. Update the chatbot ID in the script tag
5. Test the chat functionality

## Production Deployment

For production deployment:

### 1. Build Frontend

```bash
cd /var/www/html/chatbot/frontend
npm run build
```

This creates a `dist/` folder with optimized files.

### 2. Build Backend

```bash
cd /var/www/html/chatbot/backend
go build -o chatbot-api cmd/api/main.go
```

This creates an executable `chatbot-api`.

### 3. Configure Production Environment

Update `backend/.env`:
```
DB_HOST=your-production-db-host
SERVER_PORT=8080
JWT_SECRET=very-strong-secret-key
OPENAI_API_KEY=your-api-key
ALLOWED_ORIGINS=https://yourdomain.com
```

### 4. Set Up Reverse Proxy (nginx)

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        root /var/www/html/chatbot/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    location /widget.js {
        proxy_pass http://localhost:8080;
    }
}
```

### 5. Set Up System Service

Create `/etc/systemd/system/chatbot-api.service`:
```ini
[Unit]
Description=Chatbot API Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/html/chatbot/backend
ExecStart=/var/www/html/chatbot/backend/chatbot-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable chatbot-api
sudo systemctl start chatbot-api
```

## Support

For issues or questions:
1. Check the TROUBLESHOOTING section above
2. Review logs: Backend logs appear in terminal
3. Check browser console for frontend errors
4. Verify all prerequisites are installed correctly

## Next Steps

After successful installation:
1. Read SETUP_GUIDE.md for usage instructions
2. Explore the admin dashboard
3. Create your first chatbot
4. Add knowledge base entries
5. Customize widget appearance
6. Embed widget on your website

Congratulations! Your chatbot platform is now installed and ready to use.


