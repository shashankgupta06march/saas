# Quick Start Guide

Get your chatbot platform running in 15 minutes!

## Before You Begin

You need:
- Ubuntu/Linux system
- MySQL installed (username: admin, password: Admin@123)
- OpenAI API key (get from https://platform.openai.com/api-keys)

## Step 1: Install Prerequisites (5 minutes)

### Install Go
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Install Node.js
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### Verify Installation
```bash
go version    # Should show go1.21.x
node --version  # Should show v18.x.x
```

## Step 2: Setup Project (5 minutes)

```bash
cd /var/www/html/chatbot

# Run automated setup
chmod +x setup.sh
./setup.sh

# Add your OpenAI API key
nano backend/.env
# Set: OPENAI_API_KEY=sk-your-key-here
```

## Step 3: Seed Demo Data (Optional - 1 minute)

**Option A: Use Demo Accounts (Recommended for Testing)**
```bash
cd /var/www/html/chatbot
./seed.sh go
```

This creates:
- 8 Organizations (Free, Premium, Enterprise)
- 11 Users with various roles
- 3 Sample Chatbots
- 4 Knowledge Base Entries

Login with: `admin@test.com` / `password123`

**Option B: Create Your Own Account**
Skip to Step 4 and register manually.

## Step 4: Start Services (2 minutes)

### Terminal 1 - Start Backend
```bash
cd /var/www/html/chatbot/backend
go run cmd/api/main.go
```

### Terminal 2 - Start Frontend
```bash
cd /var/www/html/chatbot/frontend
npm run dev
```

## Step 5: Login/Create Account (2 minutes)

**If you seeded demo data (Step 3):**
1. Open http://localhost:3000
2. Login with: `admin@test.com` / `password123`
3. Skip to Step 6

**If creating new account:**
1. Open http://localhost:3000
2. Click "Register"
3. Fill in:
   - Organization: "My Company"
   - Email: your@email.com
   - Password: password123
4. Click "Register"

## Step 6: Create Chatbot (5 minutes)

**If you used demo data:** You already have 3 chatbots! Click on "Test Company Support Bot" to see it.

**If creating new chatbot:**

1. Click "Create Chatbot"
2. Name: "Support Bot"
3. Description: "Customer support"
4. Click "Create"

5. Click "Knowledge" button
6. Click "Add Knowledge"
7. Add training data:
   ```
   Title: Company Info
   Content: We are a software company that provides amazing products. 
   Our business hours are 9 AM to 5 PM. Contact us at support@company.com
   ```
8. Click "Add" and wait for processing

9. Click back, then "Manage"
10. Customize widget:
    - Pick a color
    - Set welcome message
    - Click "Save Settings"

11. Copy the widget code

## Step 7: Test Widget (1 minute)

```bash
# Open demo page
cd /var/www/html/chatbot/widget
# Edit demo.html and update chatbot ID to 1
# Open demo.html in browser

# Or paste widget code on your website
```

## Done! 🎉

Your chatbot is now:
- ✅ Trained with your data
- ✅ Customized with your branding
- ✅ Ready to embed on your website

## Quick Commands Reference

```bash
# Start backend
cd /var/www/html/chatbot/backend && go run cmd/api/main.go

# Start frontend
cd /var/www/html/chatbot/frontend && npm run dev

# View logs
# Backend logs show in terminal
# Frontend: Check browser console (F12)

# Restart database
mysql -u admin -pAdmin@123 -e "DROP DATABASE chatbot_saas;"
mysql -u admin -pAdmin@123 < backend/migrations/001_initial_schema.sql
```

## Common Issues

**Go not found?**
```bash
export PATH=$PATH:/usr/local/go/bin
source ~/.bashrc
```

**Port 8080 in use?**
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

**OpenAI errors?**
- Check API key in backend/.env
- Verify credits at https://platform.openai.com/account/billing

## Next Steps

- Add more knowledge to train your bot
- Embed widget on your website
- Check analytics for conversations
- Invite team members

Need help? Check INSTALLATION.md or SETUP_GUIDE.md



