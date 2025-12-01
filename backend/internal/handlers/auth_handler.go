package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	orgRepo  *repository.OrganizationRepository
	jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, orgRepo *repository.OrganizationRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		orgRepo:   orgRepo,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	OrganizationName string `json:"organization_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Generate API key for organization
	apiKey, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Create organization
	org := &models.Organization{
		Name:     req.OrganizationName,
		APIKey:   apiKey,
		PlanType: "free",
		Status:   "active",
	}

	err = h.orgRepo.Create(org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := &models.User{
		OrganizationID: org.ID,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		Role:           "admin",
	}

	err = h.userRepo.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.OrganizationID, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"organization_id": user.OrganizationID,
			"role":            user.Role,
		},
		"organization": gin.H{
			"id":      org.ID,
			"name":    org.Name,
			"api_key": org.APIKey,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get organization
	org, err := h.orgRepo.GetByID(user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.OrganizationID, user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"organization_id": user.OrganizationID,
			"role":            user.Role,
		},
		"organization": gin.H{
			"id":      org.ID,
			"name":    org.Name,
			"api_key": org.APIKey,
		},
	})
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}



