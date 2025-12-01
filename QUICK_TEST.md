# 🚀 Quick Test Guide - 3 Easy Steps

## Fastest Way to Test Your Chatbot

---

## Option 1: 🎯 One Command Test (Easiest!)

### Run this command:
```bash
/var/www/html/chatbot/start-test-server.sh
```

Then open your browser to:
```
http://localhost:8000/test-chatbot.html
```

**That's it!** Your chatbot widget will appear in the bottom-right corner.

---

## Option 2: 📝 Open Test File Directly

### Simply open in your browser:
```bash
# Linux/Mac
xdg-open /var/www/html/chatbot/test-chatbot.html

# Or drag and drop this file into your browser:
/var/www/html/chatbot/test-chatbot.html
```

---

## Option 3: 🔗 Copy Widget Code

### Get your widget code from the admin panel:

1. Login: `http://localhost:3001` (or 3000)
2. Go to **Chatbots** → Select your chatbot
3. Click **"COPY CODE"** button
4. Paste into any HTML page before `</body>`

### Or use this code directly:
```html
<script>
  (function() {
    var script = document.createElement('script');
    script.src = 'http://localhost:8081/widget.js';
    script.setAttribute('data-chatbot-id', '1');
    document.body.appendChild(script);
  })();
</script>
```

---

## ✅ What to Check

### Visual Check:
- ✅ Widget icon appears in bottom corner
- ✅ Click icon opens chat window
- ✅ Your colors and settings applied

### Functionality Check:
- ✅ Type a message
- ✅ Chatbot responds
- ✅ Answers are relevant to your knowledge base

---

## 🎯 Test Questions

Try asking:
- "Hi" or "Hello"
- Questions about your uploaded documents
- Questions from your scraped websites
- Your custom FAQ topics

---

## 🐛 If Something's Wrong

### Widget Not Showing?
```bash
# Check if backend is running
curl http://localhost:8081/widget.js

# If not, restart:
/var/www/html/chatbot/start-backend.sh
```

### No Response from Bot?
```bash
# Check backend logs
tail -f /tmp/backend.log

# Verify knowledge base has entries
# Login to admin → Chatbots → Knowledge Base
```

### CORS Error?
```bash
# Check browser console (F12)
# Then check .env file:
cat /var/www/html/chatbot/backend/.env | grep ALLOWED_ORIGINS
```

---

## 📊 Current Setup

### Your Services:
- ✅ **Backend**: http://localhost:8081
- ✅ **Admin Panel**: http://localhost:3001 (or 3000)
- ✅ **Widget Script**: http://localhost:8081/widget.js
- ✅ **Database**: MySQL (chatbot_saas)

### Login Credentials:
- **Email**: admin@test.com
- **Password**: password123

---

## 🎉 Quick Commands

```bash
# Start test server
/var/www/html/chatbot/start-test-server.sh

# Restart backend
/var/www/html/chatbot/stop-backend.sh && \
/var/www/html/chatbot/start-backend.sh

# View logs
tail -f /tmp/backend.log

# Check if services are running
ss -tulpn | grep -E ":(8081|3000|3001)"
```

---

## 📱 Test Pages Available

1. **Beautiful Test Page**: `/var/www/html/chatbot/test-chatbot.html`
   - Full-featured with styling
   - Instructions and examples
   - Best for demonstration

2. **Simple Demo**: `/var/www/html/chatbot/widget/demo.html`
   - Minimal example
   - Good for developers

---

## 🚀 Ready?

**Just run:**
```bash
/var/www/html/chatbot/start-test-server.sh
```

**Then open:**
```
http://localhost:8000/test-chatbot.html
```

**Click the chat icon and start testing!** 🎯



