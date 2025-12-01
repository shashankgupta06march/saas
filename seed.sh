#!/bin/bash

# Seed Database Script
# This script seeds the database with demo organizations and users

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  Chatbot SaaS - Database Seeder${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# Load environment variables from .env file if it exists
if [ -f "backend/.env" ]; then
    echo -e "${BLUE}📄 Loading configuration from backend/.env${NC}"
    export $(grep -v '^#' backend/.env | xargs)
elif [ -f ".env" ]; then
    echo -e "${BLUE}📄 Loading configuration from .env${NC}"
    export $(grep -v '^#' .env | xargs)
else
    echo -e "${YELLOW}⚠️  No .env file found, using environment variables or defaults${NC}"
fi

# Get database credentials from environment (now loaded from .env if available)
DB_USER="${DB_USER:-admin}"
DB_PASSWORD="${DB_PASSWORD:-Admin@123}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-chatbot_saas}"

echo -e "${BLUE}🔧 Database Configuration:${NC}"
echo -e "   Host: ${DB_HOST}:${DB_PORT}"
echo -e "   Database: ${DB_NAME}"
echo -e "   User: ${DB_USER}"
echo ""

# Check if tables exist (organizations table is required)
TABLE_CHECK=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES LIKE 'organizations';" 2>/dev/null | tail -n +2 | wc -l)

if [ -z "$TABLE_CHECK" ] || [ "$TABLE_CHECK" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Database tables not found. Running migrations automatically...${NC}"
    echo ""
    echo -e "${BLUE}📦 Running database migrations...${NC}"
    
    # Run migration 001 (replace hardcoded database name)
    echo -e "${BLUE}  → Running 001_initial_schema.sql...${NC}"
    sed "s/chatbot_saas/$DB_NAME/g" backend/migrations/001_initial_schema.sql | \
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" 2>&1 | grep -v "Using a password"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✓ Migration 001 completed${NC}"
    else
        echo -e "${RED}  ✗ Migration 001 failed${NC}"
        echo -e "${YELLOW}  Manual: sed 's/chatbot_saas/$DB_NAME/g' backend/migrations/001_initial_schema.sql | mysql -u $DB_USER -p$DB_PASSWORD${NC}"
        exit 1
    fi
    
    # Run migration 002 (already uses current database)
    echo -e "${BLUE}  → Running 002_enhanced_knowledge_base.sql...${NC}"
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < backend/migrations/002_enhanced_knowledge_base.sql 2>&1 | grep -v "Using a password"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✓ Migration 002 completed${NC}"
    else
        echo -e "${RED}  ✗ Migration 002 failed${NC}"
        echo -e "${YELLOW}  Manual: mysql -u $DB_USER -p$DB_PASSWORD $DB_NAME < backend/migrations/002_enhanced_knowledge_base.sql${NC}"
        exit 1
    fi
    
    echo ""
    echo -e "${GREEN}✅ All migrations completed successfully!${NC}"
    echo ""
else
    echo -e "${GREEN}✅ Database tables found, skipping migrations${NC}"
    echo ""
fi

# Check if Go seeder or SQL seeder should be used
SEEDER_TYPE="${1:-go}"

if [ "$SEEDER_TYPE" = "sql" ]; then
    echo -e "${YELLOW}Using SQL seeder...${NC}"
    echo ""
    
    # Check if MySQL client is available
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}❌ MySQL client not found. Please install it first.${NC}"
        exit 1
    fi
    
    # Run SQL seeder
    echo -e "${BLUE}🌱 Running SQL seeder...${NC}"
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < backend/migrations/003_seed_data.sql 2>&1 | grep -v "Using a password"
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✅ Database seeded successfully using SQL!${NC}"
        echo ""
        echo -e "${BLUE}📝 Sample Login Credentials:${NC}"
        echo -e "${BLUE}──────────────────────────────────────────${NC}"
        echo -e "  Test Company (Free Plan):"
        echo -e "    Email: admin@test.com"
        echo -e "    Password: password123"
        echo ""
        echo -e "  Demo Corp (Premium Plan):"
        echo -e "    Email: demo@democorp.com"
        echo -e "    Password: password123"
        echo ""
        echo -e "  Enterprise Solutions (Enterprise Plan):"
        echo -e "    Email: admin@enterprise.com"
        echo -e "    Password: password123"
        echo -e "${BLUE}──────────────────────────────────────────${NC}"
        echo ""
        echo -e "See backend/migrations/003_seed_data.sql for complete list"
    else
        echo -e "${RED}❌ Failed to seed database${NC}"
        exit 1
    fi
    
elif [ "$SEEDER_TYPE" = "go" ]; then
    echo -e "${YELLOW}Using Go seeder...${NC}"
    echo ""
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go not found. Please install Go first.${NC}"
        exit 1
    fi
    
    # Navigate to backend directory
    cd backend
    
    # Set up Go path
    export PATH=$HOME/go/bin:$PATH
    
    # Export database credentials
    export DB_USER="$DB_USER"
    export DB_PASSWORD="$DB_PASSWORD"
    export DB_HOST="$DB_HOST"
    export DB_PORT="$DB_PORT"
    export DB_NAME="$DB_NAME"
    
    # Run Go seeder
    echo -e "${BLUE}🌱 Running Go seeder...${NC}"
    echo ""
    go run cmd/seeder/main.go
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✅ Database seeded successfully using Go!${NC}"
    else
        echo -e "${RED}❌ Failed to seed database${NC}"
        exit 1
    fi
    
    cd ..
else
    echo -e "${RED}❌ Invalid seeder type. Use 'go' or 'sql'${NC}"
    echo -e "Usage: ./seed.sh [go|sql]"
    exit 1
fi

echo ""
echo -e "${GREEN}🚀 You can now start the application and login with the seeded credentials!${NC}"
echo ""

