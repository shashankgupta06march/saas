package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	repo     *repository.ChatbotRepository
	leadRepo *repository.LeadRepository
}

func NewChatbotHandler(repo *repository.ChatbotRepository, leadRepo *repository.LeadRepository) *ChatbotHandler {
	return &ChatbotHandler{repo: repo, leadRepo: leadRepo}
}

type CreateChatbotRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *ChatbotHandler) Create(c *gin.Context) {
	var req CreateChatbotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organizationID := c.GetInt64("organization_id")

	chatbot := &models.Chatbot{
		OrganizationID: organizationID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         "active",
	}

	err := h.repo.Create(chatbot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chatbot"})
		return
	}

	c.JSON(http.StatusCreated, chatbot)
}

func (h *ChatbotHandler) GetAll(c *gin.Context) {
	organizationID := c.GetInt64("organization_id")

	chatbots, err := h.repo.GetByOrganization(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chatbots"})
		return
	}

	c.JSON(http.StatusOK, chatbots)
}

func (h *ChatbotHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	chatbot, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatbot not found"})
		return
	}

	// Verify ownership
	organizationID := c.GetInt64("organization_id")
	if chatbot.OrganizationID != organizationID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, chatbot)
}

func (h *ChatbotHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	var req CreateChatbotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing chatbot
	chatbot, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatbot not found"})
		return
	}

	// Verify ownership
	organizationID := c.GetInt64("organization_id")
	if chatbot.OrganizationID != organizationID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	chatbot.Name = req.Name
	chatbot.Description = req.Description

	err = h.repo.Update(chatbot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chatbot"})
		return
	}

	c.JSON(http.StatusOK, chatbot)
}

func (h *ChatbotHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	// Get existing chatbot
	chatbot, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatbot not found"})
		return
	}

	// Verify ownership
	organizationID := c.GetInt64("organization_id")
	if chatbot.OrganizationID != organizationID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chatbot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chatbot deleted successfully"})
}

func (h *ChatbotHandler) GetSettings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	settings, err := h.repo.GetSettings(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settings not found"})
		return
	}

	// Include lead capture config so the widget only needs one settings call.
	leadCfg, err := h.leadRepo.GetConfig(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lead capture config"})
		return
	}

	var fields []models.LeadCaptureField
	json.Unmarshal([]byte(leadCfg.Fields), &fields)
	if fields == nil {
		fields = []models.LeadCaptureField{}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              settings.ID,
		"chatbot_id":      settings.ChatbotID,
		"theme_color":     settings.ThemeColor,
		"position":        settings.Position,
		"welcome_message": settings.WelcomeMessage,
		"avatar_url":      settings.AvatarURL,
		"custom_css":      settings.CustomCSS,
		"widget_size":     settings.WidgetSize,
		"lead_capture": gin.H{
			"enabled":  leadCfg.Enabled,
			"title":    leadCfg.Title,
			"subtitle": leadCfg.Subtitle,
			"fields":   fields,
		},
	})
}

func (h *ChatbotHandler) UpdateSettings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	var settings models.ChatbotSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings.ChatbotID = id

	err = h.repo.UpdateSettings(&settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}



