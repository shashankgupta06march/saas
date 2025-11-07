# Backend Server Commands

## 📝 Configuration

The backend uses a `.env` file for configuration. The file is located at:
```
/var/www/html/chatbot/backend/.env
```

## 🚀 Start Backend Server

### Option 1: Using Helper Script (Recommended)
```bash
/var/www/html/chatbot/start-backend.sh
```

### Option 2: Manual Command
```bash
cd /var/www/html/chatbot/backend
export PATH=$HOME/go/bin:$PATH
go run cmd/api/main.go
```

## 🛑 Stop Backend Server

### Option 1: Using Helper Script
```bash
/var/www/html/chatbot/stop-backend.sh
```

### Option 2: Manual Command
```bash
pkill -f "go run cmd/api/main.go"
```

## 📊 View Server Logs

### Live logs (follow mode):
```bash
tail -f /tmp/backend.log
```

### Last 50 lines:
```bash
tail -50 /tmp/backend.log
```

### All logs:
```bash
cat /tmp/backend.log
```

## 🔧 Backend Configuration (.env)

Current configuration:
- **Server Port:** 8081
- **Database:** chatbot_saas (MySQL)
- **Frontend URL:** http://localhost:3000
- **OpenAI API:** Configured ✅

## 📍 Server Status

Check if backend is running:
```bash
ps aux | grep "go run cmd/api/main.go" | grep -v grep
```

Check port 8081:
```bash
ss -tulpn | grep :8081
```

## 🔄 Restart Backend

```bash
/var/www/html/chatbot/stop-backend.sh && sleep 2 && /var/www/html/chatbot/start-backend.sh
```

## 🌐 API Base URL

```
http://localhost:8081/api
```

## 📌 Important Files

- **Configuration:** `/var/www/html/chatbot/backend/.env`
- **Template:** `/var/www/html/chatbot/backend/.env.example`
- **Server Logs:** `/tmp/backend.log`
- **Main Entry:** `/var/www/html/chatbot/backend/cmd/api/main.go`

---

**Note:** The `.env` file contains your OpenAI API key and is not committed to Git for security reasons.


