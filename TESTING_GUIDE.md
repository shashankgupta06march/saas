# 🧪 Chatbot Testing Guide

## Quick Start - 3 Ways to Test Your Chatbot

---

## Method 1: 🚀 Use the Beautiful Test Page (Recommended)

### Step 1: Open the Test Page
Simply open this file in your browser:
```
file:///var/www/html/chatbot/test-chatbot.html
```

Or if you're using a web server:
```
http://localhost/chatbot/test-chatbot.html
```

### Step 2: Look for the Chatbot Widget
- You'll see a chatbot icon in the **bottom-right corner**
- Click on it to open the chat window

### Step 3: Start Testing
- Type any question
- Try questions related to your knowledge base
- Test the conversation flow

---

## Method 2: 🎯 Copy Widget Code to Any Page

### Get Your Widget Code:
From your admin panel, go to:
**Chatbots → Select Your Chatbot → Widget Code**

### Current Widget Code for Chatbot ID 1:
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

### Paste it anywhere:
- Add this code before the closing `</body>` tag of any HTML page
- The chatbot will appear automatically

---

## Method 3: 🌐 Direct Browser Testing

### Option A: File Protocol
1. Navigate to: `/var/www/html/chatbot/test-chatbot.html`
2. Right-click → Open with → Your Browser

### Option B: Local Web Server
If you have Python installed:
```bash
cd /var/www/html/chatbot
python3 -m http.server 8000
```
Then open: `http://localhost:8000/test-chatbot.html`

### Option C: PHP Built-in Server
```bash
cd /var/www/html/chatbot
php -S localhost:8000
```
Then open: `http://localhost:8000/test-chatbot.html`

---

## 📝 Testing Checklist

### ✅ Visual Testing
- [ ] Chatbot widget appears in correct position
- [ ] Colors match your settings
- [ ] Widget size is appropriate
- [ ] Welcome message displays correctly
- [ ] Avatar shows (if configured)

### ✅ Functionality Testing
- [ ] Click on widget opens chat window
- [ ] Can type messages
- [ ] Chatbot responds to questions
- [ ] Responses are relevant to knowledge base
- [ ] Can close chat window
- [ ] Can minimize/maximize chat

### ✅ Knowledge Base Testing
- [ ] Ask questions from uploaded PDFs
- [ ] Ask about scraped website content
- [ ] Test with manual text entries
- [ ] Verify answers are accurate
- [ ] Check response time

### ✅ Conversation Flow
- [ ] Multi-turn conversations work
- [ ] Bot remembers context
- [ ] Handles unknown questions gracefully
- [ ] Fallback responses work

---

## 🎯 Sample Test Scenarios

### Scenario 1: Basic Interaction
```
User: "Hi"
Expected: Welcome message appears
```

### Scenario 2: Knowledge Base Query
```
User: "What services do you offer?"
Expected: Response based on uploaded knowledge
```

### Scenario 3: Complex Question
```
User: "Can you tell me about your pricing and support options?"
Expected: Detailed response combining multiple knowledge base entries
```

### Scenario 4: Unknown Question
```
User: "What's the weather today?"
Expected: Polite response indicating inability to answer
```

---

## 🔧 Test Different Scenarios

### Test with Different Content Types:
1. **PDF Content**: Ask about information from uploaded PDFs
2. **Website Content**: Query scraped website data
3. **Manual Text**: Test custom FAQ entries

### Test Chatbot Settings:
1. Change widget position → Refresh test page
2. Update colors → Check visual changes
3. Modify welcome message → Verify it displays
4. Change widget size → Test small/medium/large

---

## 🐛 Troubleshooting

### Widget Not Appearing?
**Check:**
1. Backend is running: `http://localhost:8081/widget.js` should load
2. Browser console for errors (F12)
3. Correct chatbot ID in the widget code
4. CORS is properly configured

**Solution:**
```bash
# Restart backend
/var/www/html/chatbot/stop-backend.sh
/var/www/html/chatbot/start-backend.sh
```

### Chatbot Not Responding?
**Check:**
1. Knowledge base has entries
2. OpenAI API key is valid
3. Backend logs: `tail -f /tmp/backend.log`

**Solution:**
```bash
# Check backend logs
tail -f /tmp/backend.log

# Verify OpenAI key
cat /var/www/html/chatbot/backend/.env | grep OPENAI
```

### CORS Errors?
**Check browser console for:**
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**Solution:**
```bash
# Edit .env to add your domain
nano /var/www/html/chatbot/backend/.env

# Add to ALLOWED_ORIGINS:
ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3001,http://localhost:8000

# Restart backend
/var/www/html/chatbot/stop-backend.sh && /var/www/html/chatbot/start-backend.sh
```

---

## 📊 Monitor Your Tests

### Backend Logs:
```bash
tail -f /tmp/backend.log
```

### Watch for:
- Chat requests: `POST /api/chat/:chatbot_id`
- OpenAI API calls
- Error messages
- Response times

### Browser Developer Tools (F12):
- **Console**: JavaScript errors
- **Network**: API calls and responses
- **Storage**: Check for stored data

---

## 🚀 Advanced Testing

### Test API Directly:
```bash
# Get chatbot settings
curl http://localhost:8081/api/chatbots/1/settings

# Send a test message
curl -X POST http://localhost:8081/api/chat/1 \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello, how can you help me?",
    "conversation_id": ""
  }'
```

### Test Widget Loading:
```bash
# Should return JavaScript code
curl http://localhost:8081/widget.js
```

---

## 📱 Test on Different Devices

### Desktop:
- Chrome
- Firefox
- Safari
- Edge

### Mobile:
- iOS Safari
- Android Chrome
- Responsive mode in desktop browsers (F12 → Device toolbar)

### Screen Sizes:
- Desktop (1920x1080)
- Laptop (1366x768)
- Tablet (768x1024)
- Mobile (375x667)

---

## ✨ Tips for Better Testing

1. **Clear Browser Cache**: Hard refresh with Ctrl+Shift+R (or Cmd+Shift+R on Mac)
2. **Use Incognito Mode**: Test without cached data
3. **Test Incrementally**: Make one change at a time
4. **Document Issues**: Keep track of bugs found
5. **Test Edge Cases**: Try unusual inputs and questions

---

## 🎓 What to Test

### Must Test:
✅ Widget loads and displays
✅ Chat opens when clicked
✅ Messages send successfully
✅ Bot responds with relevant answers
✅ Widget is mobile-responsive

### Should Test:
✅ Multiple conversations
✅ Long messages
✅ Special characters
✅ Emoji support
✅ Different languages (if supported)

### Nice to Test:
✅ Performance with large knowledge base
✅ Concurrent users (if possible)
✅ Offline behavior
✅ Session persistence

---

## 📋 Test Report Template

After testing, document your results:

```
Date: [Date]
Tester: [Your Name]
Chatbot ID: 1
Knowledge Base Entries: [Number]

✅ Working:
- Widget loads correctly
- Responses are accurate
- [Add more...]

❌ Issues Found:
- [Describe issue]
- [Steps to reproduce]

💡 Improvements Needed:
- [Suggestions]
```

---

## 🎉 Ready to Deploy?

Once testing is complete:
1. ✅ All features work as expected
2. ✅ No console errors
3. ✅ Responses are accurate
4. ✅ Visual design is correct
5. ✅ Mobile responsive

**Next Steps:**
1. Update widget URL for production
2. Update CORS settings for your domain
3. Deploy to your live website
4. Monitor analytics

---

## 📞 Need Help?

Check the logs:
```bash
# Backend logs
tail -f /tmp/backend.log

# Frontend logs (if running dev server)
tail -f /tmp/frontend.log
```

Restart services:
```bash
# Restart everything
/var/www/html/chatbot/stop-backend.sh && sleep 2 && /var/www/html/chatbot/start-backend.sh
```

**Happy Testing! 🚀**


