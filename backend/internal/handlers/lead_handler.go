package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type LeadHandler struct {
	repo         *repository.LeadRepository
	chatbotRepo  *repository.ChatbotRepository
	convRepo     *repository.ConversationRepository
}

func NewLeadHandler(repo *repository.LeadRepository, chatbotRepo *repository.ChatbotRepository, convRepo *repository.ConversationRepository) *LeadHandler {
	return &LeadHandler{repo: repo, chatbotRepo: chatbotRepo, convRepo: convRepo}
}

// GetConfig returns lead capture config for a chatbot (protected).
func (h *LeadHandler) GetConfig(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	cfg, err := h.repo.GetConfig(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lead capture config"})
		return
	}

	c.JSON(http.StatusOK, h.configResponse(cfg))
}

// UpsertConfig saves lead capture config for a chatbot (protected).
func (h *LeadHandler) UpsertConfig(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	// Verify ownership
	chatbot, err := h.chatbotRepo.GetByID(chatbotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatbot not found"})
		return
	}
	orgID := c.GetInt64("organization_id")
	if chatbot.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var body struct {
		Enabled  bool                       `json:"enabled"`
		Title    string                     `json:"title"`
		Subtitle string                     `json:"subtitle"`
		Fields   []models.LeadCaptureField  `json:"fields"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldsJSON, err := json.Marshal(body.Fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode fields"})
		return
	}

	cfg := &models.LeadCaptureConfig{
		ChatbotID: chatbotID,
		Enabled:   body.Enabled,
		Title:     body.Title,
		Subtitle:  body.Subtitle,
		Fields:    string(fieldsJSON),
	}

	if err := h.repo.UpsertConfig(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save lead capture config"})
		return
	}

	c.JSON(http.StatusOK, h.configResponse(cfg))
}

// SubmitLead accepts a lead from the widget (public).
func (h *LeadHandler) SubmitLead(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	var body struct {
		SessionID   string                 `json:"session_id" binding:"required"`
		FieldValues map[string]string      `json:"field_values" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valuesJSON, err := json.Marshal(body.FieldValues)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode field values"})
		return
	}

	lead := &models.Lead{
		ChatbotID:   chatbotID,
		SessionID:   body.SessionID,
		FieldValues: string(valuesJSON),
	}

	if err := h.repo.CreateLead(lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save lead"})
		return
	}

	c.JSON(http.StatusCreated, lead)
}

// GetLeads returns all leads for a chatbot (protected).
func (h *LeadHandler) GetLeads(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	// Verify ownership
	chatbot, err := h.chatbotRepo.GetByID(chatbotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatbot not found"})
		return
	}
	orgID := c.GetInt64("organization_id")
	if chatbot.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	leads, err := h.repo.GetLeads(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leads"})
		return
	}

	// Parse field_values JSON for each lead so the response is usable without extra parsing.
	type LeadResponse struct {
		ID          int64                  `json:"id"`
		ChatbotID   int64                  `json:"chatbot_id"`
		SessionID   string                 `json:"session_id"`
		FieldValues map[string]string      `json:"field_values"`
		CreatedAt   interface{}            `json:"created_at"`
	}

	var result []LeadResponse
	for _, l := range leads {
		var vals map[string]string
		json.Unmarshal([]byte(l.FieldValues), &vals)
		result = append(result, LeadResponse{
			ID:          l.ID,
			ChatbotID:   l.ChatbotID,
			SessionID:   l.SessionID,
			FieldValues: vals,
			CreatedAt:   l.CreatedAt,
		})
	}

	if result == nil {
		result = []LeadResponse{}
	}

	c.JSON(http.StatusOK, result)
}

// GetSessionMessages returns the chat messages for a given session (protected).
func (h *LeadHandler) GetSessionMessages(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	conv, err := h.convRepo.GetBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find conversation"})
		return
	}
	if conv == nil {
		c.JSON(http.StatusOK, gin.H{"messages": []interface{}{}, "conversation": nil})
		return
	}

	messages, err := h.convRepo.GetMessages(conv.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}

	c.JSON(http.StatusOK, gin.H{
		"conversation": conv,
		"messages":     messages,
	})
}

// configResponse converts stored JSON fields back to a parsed struct for the API response.
func (h *LeadHandler) configResponse(cfg *models.LeadCaptureConfig) gin.H {
	var fields []models.LeadCaptureField
	json.Unmarshal([]byte(cfg.Fields), &fields)
	if fields == nil {
		fields = []models.LeadCaptureField{}
	}
	return gin.H{
		"id":        cfg.ID,
		"chatbot_id": cfg.ChatbotID,
		"enabled":   cfg.Enabled,
		"title":     cfg.Title,
		"subtitle":  cfg.Subtitle,
		"fields":    fields,
	}
}
