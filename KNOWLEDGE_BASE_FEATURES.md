# Knowledge Base - New Features

## 🎉 Enhanced Knowledge Base with File Upload & Web Scraping

Your chatbot knowledge base now supports **3 ways** to add training data:

---

## ✨ Features

### 1. 📝 **Text Input** (Original)
- Manually type or paste content
- Great for custom information and FAQs

### 2. 📄 **File Upload** (NEW!)
Supported file formats:
- **PDF files** (.pdf)
- **Word documents** (.docx)
- **Text files** (.txt)

**How it works:**
- Upload your document
- Backend automatically extracts text
- Content is processed and stored
- AI embeddings are generated for smart responses

### 3. 🌐 **Website Scraping** (NEW!)
- Enter any website URL
- Automatically scrapes and extracts content
- Supports articles, blog posts, documentation pages
- Cleans HTML and extracts readable text

---

## 🚀 How to Use

### Access the Knowledge Base:
1. Login to your admin panel
2. Navigate to **Chatbots**
3. Select your chatbot
4. Click **Knowledge Base**
5. Click **"Add Knowledge"**

### Choose Your Method:

#### **Text Tab:**
- Enter a title
- Type or paste content
- Click "Add"

#### **Upload File Tab:**
- Click "Choose File"
- Select PDF, DOCX, or TXT file
- Optionally add a custom title
- Click "Add"

#### **Website URL Tab:**
- Enter the full URL (e.g., https://example.com)
- Optionally add a custom title
- Click "Add"

---

## 📊 API Endpoints

### Upload File:
```
POST /api/knowledge/upload
Content-Type: multipart/form-data

Form Data:
- file: (binary file)
- chatbot_id: (integer)
- title: (string, optional)
```

### Scrape URL:
```
POST /api/knowledge/scrape
Content-Type: application/json

Body:
{
  "chatbot_id": 1,
  "url": "https://example.com",
  "title": "Optional Title"
}
```

### Text Input (Original):
```
POST /api/knowledge
Content-Type: application/json

Body:
{
  "chatbot_id": 1,
  "title": "Title",
  "content": "Your content here",
  "content_type": "text"
}
```

---

## 🔧 Technical Details

### Backend Libraries:
- **PDF Parsing**: `github.com/ledongthuc/pdf`
- **DOCX Parsing**: `github.com/nguyenthenguyen/docx`
- **Web Scraping**: `github.com/gocolly/colly/v2`

### Content Processing:
- Maximum content size: **50,000 characters**
- Text extraction removes formatting
- Web scraper removes scripts, styles, and ads
- All content is cleaned and normalized

### Content Types:
- `text` - Manual text input
- `pdf` - PDF documents
- `docx` - Word documents
- `webpage` - Scraped websites

---

## 🎨 UI Features

### Visual Indicators:
- Each knowledge entry shows its type with a colored chip:
  - 🔴 **PDF** - Red chip
  - 🔵 **DOCX** - Blue chip
  - 🟢 **Text** - Green chip
  - 🔵 **Webpage** - Info blue chip

### Tabs Interface:
- Clean tab-based dialog
- Icons for each input method
- Real-time validation
- Loading states during processing

---

## 📝 Examples

### Example 1: Upload Product Manual
```
1. Click "Add Knowledge"
2. Go to "Upload File" tab
3. Choose "product-manual.pdf"
4. Click "Add"
✅ PDF content extracted and stored
```

### Example 2: Scrape Company Blog
```
1. Click "Add Knowledge"
2. Go to "Website URL" tab
3. Enter: https://yourblog.com/article
4. Click "Add"
✅ Website content scraped and stored
```

### Example 3: Add FAQ
```
1. Click "Add Knowledge"
2. Stay on "Text" tab
3. Enter title: "Shipping Policy"
4. Paste FAQ content
5. Click "Add"
✅ Text content stored
```

---

## ⚙️ Server Status

### Backend: ✅ Running
- Port: **8081**
- API: `http://localhost:8081/api`
- Logs: `/tmp/backend.log`

### Frontend: ✅ Running
- Port 3000: `http://localhost:3000`
- Port 3001: `http://localhost:3001`
- Logs: `/tmp/frontend.log`

### Database: ✅ Connected
- MySQL: `chatbot_saas`
- Host: localhost:3306

---

## 🔒 Security Notes

- File uploads are limited to safe formats (PDF, DOCX, TXT)
- Maximum file size handling in place
- Content is sanitized before storage
- Web scraping respects robots.txt
- CORS enabled for localhost origins

---

## 🐛 Troubleshooting

### File Upload Not Working:
- Check file format is supported
- Ensure file is not corrupt
- Check browser console for errors

### Website Scraping Failed:
- Verify URL is accessible
- Check if website blocks scraping
- Try URL with https:// prefix
- Some sites may have anti-scraping protection

### OpenAI Quota Error:
- Add billing to your OpenAI account
- Check usage at: https://platform.openai.com/usage
- Verify API key has credits

---

## 📚 Next Steps

1. Upload your product documentation
2. Scrape your company blog posts
3. Add FAQs and support content
4. Test the chatbot responses
5. Monitor in the Analytics section

---

**Need Help?** Check the logs:
```bash
# Backend logs
tail -f /tmp/backend.log

# Frontend logs
tail -f /tmp/frontend.log
```

