#!/bin/bash

echo "=========================================="
echo "Chatbot SaaS Platform Setup"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
echo "Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
else
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓ Go is installed ($GO_VERSION)${NC}"
fi

# Check Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}✗ Node.js is not installed${NC}"
    echo "Please install Node.js 18 or higher from https://nodejs.org/"
    exit 1
else
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✓ Node.js is installed ($NODE_VERSION)${NC}"
fi

# Check MySQL
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}✗ MySQL is not installed${NC}"
    echo "Please install MySQL 8.0 or higher"
    exit 1
else
    echo -e "${GREEN}✓ MySQL is installed${NC}"
fi

echo ""
echo "=========================================="
echo "Step 1: Database Setup"
echo "=========================================="

# Check if database exists
DB_EXISTS=$(mysql -u admin -pAdmin@123 -e "SHOW DATABASES LIKE 'chatbot_saas';" 2>/dev/null | grep chatbot_saas)

if [ -z "$DB_EXISTS" ]; then
    echo "Creating database..."
    mysql -u admin -pAdmin@123 < backend/migrations/001_initial_schema.sql 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Database created successfully${NC}"
    else
        echo -e "${RED}✗ Failed to create database${NC}"
        echo "Please check your MySQL credentials and try again"
        exit 1
    fi
else
    echo -e "${YELLOW}Database already exists, skipping creation${NC}"
fi

echo ""
echo "=========================================="
echo "Step 2: Backend Setup"
echo "=========================================="

cd backend

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from .env.example...${NC}"
    cp .env.example .env
    echo -e "${YELLOW}⚠ Please update the OPENAI_API_KEY in backend/.env${NC}"
fi

echo "Installing Go dependencies..."
go mod download

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend dependencies installed${NC}"
else
    echo -e "${RED}✗ Failed to install backend dependencies${NC}"
    exit 1
fi

cd ..

echo ""
echo "=========================================="
echo "Step 3: Frontend Setup"
echo "=========================================="

cd frontend

if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Frontend dependencies installed${NC}"
    else
        echo -e "${RED}✗ Failed to install frontend dependencies${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Frontend dependencies already installed${NC}"
fi

cd ..

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "To start the application:"
echo ""
echo "1. Start the backend (Terminal 1):"
echo "   cd backend && go run cmd/api/main.go"
echo ""
echo "2. Start the frontend (Terminal 2):"
echo "   cd frontend && npm run dev"
echo ""
echo "3. Access the application:"
echo "   Admin Dashboard: http://localhost:3000"
echo "   API: http://localhost:8080"
echo ""
echo -e "${YELLOW}⚠ Don't forget to set your OpenAI API key in backend/.env${NC}"
echo ""

