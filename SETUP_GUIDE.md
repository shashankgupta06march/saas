# Quick Setup Guide

## Prerequisites Check

Before starting, ensure you have:
- [x] Go 1.21+ installed
- [x] Node.js 18+ installed  
- [x] MySQL 8.0+ running
- [x] OpenAI API key

## Step-by-Step Setup

### 1. Database Setup (5 minutes)

```bash
# Login to MySQL
mysql -u admin -pAdmin@123

# Run the migration
mysql -u admin -pAdmin@123 < backend/migrations/001_initial_schema.sql

# Verify database created
mysql -u admin -pAdmin@123 -e "USE chatbot_saas; SHOW TABLES;"
```

You should see 8 tables created.

### 2. Backend Setup (5 minutes)

```bash
cd backend

# Install Go dependencies
go mod download

# Update .env with your OpenAI API key
# Edit backend/.env and set OPENAI_API_KEY=sk-your-key-here

# Run the backend
go run cmd/api/main.go
```

You should see: "Server starting on port 8080"

### 3. Frontend Setup (5 minutes)

Open a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Start the development server
npm run dev
```

You should see: "Local: http://localhost:3000"

### 4. Access the Application

1. Open browser and go to: http://localhost:3000
2. Click "Register" to create your account
3. Fill in:
   - Organization Name: "My Company"
   - Email: your@email.com
   - Password: (min 6 characters)

### 5. Create Your First Chatbot

1. Click "Create Chatbot" button
2. Enter:
   - Name: "Support Bot"
   - Description: "Customer support chatbot"
3. Click "Create"

### 6. Train Your Chatbot

1. Click "Knowledge" on your chatbot
2. Click "Add Knowledge"
3. Enter:
   - Title: "Product Info"
   - Content: "We offer premium web hosting services..."
4. Click "Add"

Wait a few seconds for the embedding to be generated.

### 7. Customize Widget

1. Click "Manage" on your chatbot
2. Customize:
   - Theme Color: Choose your brand color
   - Position: bottom-right or bottom-left
   - Welcome Message: "Hello! How can I help?"
3. Click "Save Settings"
4. Copy the widget code

### 8. Test the Widget

Option A - Use Demo Page:
```bash
# Open widget/demo.html in browser
# Make sure to update the chatbot ID in the script tag
```

Option B - Test on Your Website:
```html
<!-- Paste this in your website's HTML -->
<script src="http://localhost:8080/widget.js" 
        data-chatbot-id="1" 
        data-api-url="http://localhost:8080/api">
</script>
```

### 9. Test the Chat

1. Click the chat widget button
2. Type a question related to your knowledge base
3. The chatbot should respond using the trained data

## Verification Checklist

- [ ] Backend running on http://localhost:8080
- [ ] Frontend running on http://localhost:3000
- [ ] Can register/login successfully
- [ ] Can create a chatbot
- [ ] Can add knowledge base entries
- [ ] Can customize widget settings
- [ ] Widget appears on demo page
- [ ] Can send messages and get AI responses

## Common Issues

### Backend won't start
```bash
# Check if MySQL is running
sudo systemctl status mysql

# Verify database exists
mysql -u admin -pAdmin@123 -e "SHOW DATABASES;"

# Check Go version
go version  # Should be 1.21+
```

### Frontend won't start
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Check Node version
node --version  # Should be 18+
```

### Widget not loading
- Check browser console for errors
- Verify API is running on port 8080
- Check CORS settings in backend/.env

### OpenAI API errors
- Verify API key is valid
- Check OpenAI account has credits
- Review backend logs for detailed errors

## Production Deployment

For production deployment:

1. Update `.env` with production values
2. Build frontend: `npm run build`
3. Build backend: `go build -o chatbot-api cmd/api/main.go`
4. Deploy to your server
5. Update widget URLs to production domain
6. Enable HTTPS
7. Set up proper database backups

## Next Steps

- Add more knowledge base entries
- Customize widget appearance
- View analytics and conversations
- Integrate with your website
- Monitor API usage

## Need Help?

Check the main README.md for detailed documentation.


