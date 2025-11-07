# Multi-Tenant Chatbot SaaS Platform

A complete SaaS platform for creating and managing AI-powered chatbots with embeddable widgets, knowledge base training, and customizable UI.

## Features

- **Multi-tenant Architecture**: Support multiple organizations with complete data isolation
- **AI-Powered Chatbots**: OpenAI integration for intelligent responses
- **Knowledge Base Training**: Upload and train chatbots with custom data
- **Customizable Widget**: Fully customizable embeddable chat widget
- **Admin Dashboard**: Comprehensive React-based management interface
- **Analytics**: Conversation history and usage tracking
- **RESTful API**: Complete Go-based REST API

## Technology Stack

- **Backend**: Go (Golang) with Gin framework
- **Frontend**: React.js with Material-UI
- **Database**: MySQL
- **AI**: OpenAI API (GPT-3.5-turbo & embeddings)
- **Widget**: Vanilla JavaScript

## Project Structure

```
chatbot/
├── backend/               # Go backend API
│   ├── cmd/api/          # Main application entry point
│   ├── internal/         # Internal packages
│   │   ├── config/       # Configuration management
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Middleware components
│   │   ├── models/       # Data models
│   │   ├── repository/   # Database operations
│   │   └── services/     # Business logic
│   ├── pkg/              # Public packages
│   │   ├── auth/         # Authentication
│   │   └── openai/       # OpenAI client
│   └── migrations/       # Database migrations
├── frontend/             # React admin dashboard
│   └── src/
│       ├── components/   # React components
│       ├── contexts/     # React contexts
│       └── services/     # API services
└── widget/               # Embeddable JavaScript widget
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- MySQL 8.0 or higher
- OpenAI API key

### Database Setup

1. Create MySQL database:
```bash
mysql -u admin -p
```

2. Run the migration script:
```bash
mysql -u admin -pAdmin@123 < backend/migrations/001_initial_schema.sql
```

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables in `.env`:
```
DB_HOST=localhost
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=Admin@123
DB_NAME=chatbot_saas
SERVER_PORT=8080
JWT_SECRET=your-secret-key
OPENAI_API_KEY=your-openai-api-key
ALLOWED_ORIGINS=http://localhost:3000
```

4. Run the backend:
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

The admin dashboard will be available at `http://localhost:3000`

### Widget Setup

The widget can be embedded on any website using:

```html
<script src="http://your-domain/widget/chatbot-widget.js" 
        data-chatbot-id="YOUR_CHATBOT_ID" 
        data-api-url="http://your-api-domain/api"></script>
```

Test the widget locally by opening `widget/demo.html` in a browser.

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new organization
- `POST /api/auth/login` - Login user

### Chatbots
- `GET /api/chatbots` - List all chatbots
- `POST /api/chatbots` - Create new chatbot
- `GET /api/chatbots/:id` - Get chatbot details
- `PUT /api/chatbots/:id` - Update chatbot
- `DELETE /api/chatbots/:id` - Delete chatbot
- `GET /api/chatbots/:id/settings` - Get chatbot settings
- `PUT /api/chatbots/:id/settings` - Update chatbot settings

### Knowledge Base
- `POST /api/knowledge` - Add knowledge entry
- `GET /api/knowledge/chatbot/:chatbot_id` - Get knowledge base
- `DELETE /api/knowledge/:id` - Delete knowledge entry

### Chat (Public)
- `POST /api/chat/:chatbot_id` - Send message to chatbot

### Analytics
- `GET /api/analytics/conversations/:chatbot_id` - Get conversations
- `GET /api/analytics/messages/:conversation_id` - Get messages

## Usage Guide

### 1. Create an Account
- Navigate to the admin dashboard
- Click "Register" and create your organization

### 2. Create a Chatbot
- Go to "Chatbots" section
- Click "Create Chatbot"
- Enter name and description

### 3. Train Your Chatbot
- Click "Knowledge" on your chatbot
- Add training data (text content)
- The system will automatically generate embeddings

### 4. Customize Widget
- Click "Manage" on your chatbot
- Adjust theme color, position, welcome message
- Copy the widget code

### 5. Embed on Website
- Paste the widget code in your website's HTML
- The chatbot will appear on your site

### 6. View Analytics
- Click "View Analytics" to see conversations
- View message history and usage stats

## Key Features Explained

### Multi-Tenancy
- Complete data isolation per organization
- Unique API keys for each organization
- Organization-scoped queries

### Knowledge Base with Embeddings
- Text content is converted to vector embeddings using OpenAI
- Cosine similarity search finds relevant context
- Context is injected into chat prompts for accurate responses

### Customizable Widget
- Theme colors
- Widget position (bottom-left/right)
- Welcome message
- Avatar images
- Widget size options

### Security
- JWT authentication for admin users
- Password hashing with bcrypt
- CORS protection
- Input validation
- Organization-level data isolation

## Development

### Running in Development Mode

Backend:
```bash
cd backend
go run cmd/api/main.go
```

Frontend:
```bash
cd frontend
npm run dev
```

### Building for Production

Backend:
```bash
cd backend
go build -o chatbot-api cmd/api/main.go
```

Frontend:
```bash
cd frontend
npm run build
```

## Configuration

### Environment Variables

Backend `.env`:
- `DB_HOST` - MySQL host
- `DB_PORT` - MySQL port
- `DB_USER` - MySQL username
- `DB_PASSWORD` - MySQL password
- `DB_NAME` - Database name
- `SERVER_PORT` - API server port
- `JWT_SECRET` - Secret for JWT tokens
- `OPENAI_API_KEY` - OpenAI API key
- `ALLOWED_ORIGINS` - CORS allowed origins

## Troubleshooting

### Database Connection Issues
- Verify MySQL is running
- Check database credentials in `.env`
- Ensure database exists

### OpenAI API Issues
- Verify API key is correct
- Check API quota and billing
- Ensure network connectivity

### Widget Not Loading
- Check CORS settings
- Verify API URL in widget code
- Check browser console for errors

## License

This project is proprietary software.

## Support

For support, contact your system administrator.

# saas
