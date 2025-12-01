package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/chatbot-saas/backend/pkg/auth"
	"crypto/rand"
	"encoding/hex"
)

type Organization struct {
	ID       int64
	Name     string
	APIKey   string
	PlanType string
	Status   string
}

type User struct {
	OrganizationID int64
	Email          string
	Password       string
	Role           string
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "admin"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "Admin@123"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "chatbot_saas"
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Connected to database successfully")

	// Check if data already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check existing data:", err)
	}

	if count > 0 {
		log.Printf("⚠️  Found %d existing organizations. Do you want to continue? (y/n): ", count)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			log.Println("Seeding cancelled.")
			return
		}
	}

	log.Println("🌱 Starting database seeding...")

	// Define seed data
	organizations := []Organization{
		{Name: "Test Company", PlanType: "free", Status: "active"},
		{Name: "Demo Corp", PlanType: "premium", Status: "active"},
		{Name: "Enterprise Solutions Ltd", PlanType: "enterprise", Status: "active"},
		{Name: "Startup Inc", PlanType: "free", Status: "active"},
		{Name: "Tech Innovations", PlanType: "premium", Status: "active"},
		{Name: "Global Services", PlanType: "enterprise", Status: "active"},
		{Name: "Small Business Co", PlanType: "free", Status: "trial"},
		{Name: "Medium Enterprise", PlanType: "premium", Status: "active"},
	}

	users := []User{
		{Email: "admin@test.com", Password: "password123", Role: "admin"},
		{Email: "demo@democorp.com", Password: "demo123456", Role: "admin"},
		{Email: "admin@enterprise.com", Password: "enterprise123", Role: "admin"},
		{Email: "contact@startup.com", Password: "startup123", Role: "admin"},
		{Email: "admin@techinnovations.com", Password: "tech123", Role: "admin"},
		{Email: "admin@globalservices.com", Password: "global123", Role: "admin"},
		{Email: "owner@smallbiz.com", Password: "smallbiz123", Role: "admin"},
		{Email: "admin@mediumenterprise.com", Password: "medium123", Role: "admin"},
		// Additional users for some organizations
		{Email: "user@test.com", Password: "user123", Role: "user"},
		{Email: "manager@democorp.com", Password: "manager123", Role: "manager"},
		{Email: "support@enterprise.com", Password: "support123", Role: "support"},
	}

	// Seed organizations
	log.Println("📦 Seeding organizations...")
	orgIDs := make(map[int]int64)
	
	for i, org := range organizations {
		// Check if organization already exists
		var existingID int64
		err := db.QueryRow("SELECT id FROM organizations WHERE name = ?", org.Name).Scan(&existingID)
		
		if err == sql.ErrNoRows {
			// Generate API key
			apiKey, err := generateAPIKey()
			if err != nil {
				log.Printf("❌ Failed to generate API key for %s: %v", org.Name, err)
				continue
			}

			// Insert organization
			result, err := db.Exec(
				"INSERT INTO organizations (name, api_key, plan_type, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				org.Name, apiKey, org.PlanType, org.Status, time.Now(), time.Now(),
			)
			if err != nil {
				log.Printf("❌ Failed to create organization %s: %v", org.Name, err)
				continue
			}

			id, _ := result.LastInsertId()
			orgIDs[i] = id
			log.Printf("  ✓ Created organization: %s (ID: %d, Plan: %s, API Key: %s)", 
				org.Name, id, org.PlanType, apiKey[:20]+"...")
		} else if err == nil {
			orgIDs[i] = existingID
			log.Printf("  ⚠️  Organization already exists: %s (ID: %d)", org.Name, existingID)
		} else {
			log.Printf("❌ Error checking organization %s: %v", org.Name, err)
		}
	}

	// Seed users
	log.Println("\n👥 Seeding users...")
	userCount := 0
	
	for i, user := range users {
		// Determine organization ID
		var orgID int64
		if i < len(organizations) {
			orgID = orgIDs[i]
		} else {
			// Additional users for existing organizations
			orgID = orgIDs[i-len(organizations)]
		}

		if orgID == 0 {
			log.Printf("  ⚠️  Skipping user %s - no valid organization", user.Email)
			continue
		}

		// Check if user already exists
		var existingID int64
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingID)
		
		if err == sql.ErrNoRows {
			// Hash password
			hashedPassword, err := auth.HashPassword(user.Password)
			if err != nil {
				log.Printf("❌ Failed to hash password for %s: %v", user.Email, err)
				continue
			}

			// Insert user
			result, err := db.Exec(
				"INSERT INTO users (organization_id, email, password_hash, role, created_at) VALUES (?, ?, ?, ?, ?)",
				orgID, user.Email, hashedPassword, user.Role, time.Now(),
			)
			if err != nil {
				log.Printf("❌ Failed to create user %s: %v", user.Email, err)
				continue
			}

			id, _ := result.LastInsertId()
			userCount++
			log.Printf("  ✓ Created user: %s (ID: %d, Org ID: %d, Role: %s, Password: %s)", 
				user.Email, id, orgID, user.Role, user.Password)
		} else if err == nil {
			log.Printf("  ⚠️  User already exists: %s (ID: %d)", user.Email, existingID)
		} else {
			log.Printf("❌ Error checking user %s: %v", user.Email, err)
		}
	}

	// Seed sample chatbots
	log.Println("\n🤖 Seeding sample chatbots...")
	chatbotCount := 0
	
	for i, orgID := range orgIDs {
		if i < 3 { // Only create chatbots for first 3 organizations
			orgName := organizations[i].Name
			
			result, err := db.Exec(
				"INSERT INTO chatbots (organization_id, name, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				orgID, 
				fmt.Sprintf("%s Support Bot", orgName),
				fmt.Sprintf("Customer support chatbot for %s", orgName),
				"active",
				time.Now(),
				time.Now(),
			)
			if err != nil {
				log.Printf("❌ Failed to create chatbot for %s: %v", orgName, err)
				continue
			}

			chatbotID, _ := result.LastInsertId()
			chatbotCount++
			
			// Create chatbot settings
			_, err = db.Exec(
				`INSERT INTO chatbot_settings (chatbot_id, theme_color, position, welcome_message, widget_size) 
				 VALUES (?, ?, ?, ?, ?)`,
				chatbotID,
				"#007bff",
				"bottom-right",
				fmt.Sprintf("Hi! Welcome to %s. How can I help you today?", orgName),
				"medium",
			)
			if err != nil {
				log.Printf("⚠️  Failed to create chatbot settings for chatbot %d: %v", chatbotID, err)
			}

			log.Printf("  ✓ Created chatbot: %s Support Bot (ID: %d)", orgName, chatbotID)
		}
	}

	// Print summary
	log.Println("\n" + repeat("=", 60))
	log.Println("✅ DATABASE SEEDING COMPLETED!")
	log.Println(repeat("=", 60))
	log.Printf("📊 Summary:")
	log.Printf("  - Organizations: %d", len(orgIDs))
	log.Printf("  - Users: %d", userCount)
	log.Printf("  - Chatbots: %d", chatbotCount)
	log.Println("\n📝 Sample Login Credentials:")
	log.Println("  " + repeat("-", 58))
	
	// Display credentials for each organization
	for i, org := range organizations {
		if i < len(users) {
			user := users[i]
			log.Printf("  %s:", org.Name)
			log.Printf("    Email: %s", user.Email)
			log.Printf("    Password: %s", user.Password)
			log.Printf("    Plan: %s", org.PlanType)
			log.Println()
		}
	}
	
	log.Println("  " + repeat("-", 58))
	log.Println("\n🚀 You can now start the application and login with these credentials!")
}

// Helper function to repeat strings (Go doesn't have built-in repeat for strings in older versions)
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

