package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBSyncHandler struct {
	service *services.KBSyncService
}

func NewKBSyncHandler(service *services.KBSyncService) *KBSyncHandler {
	return &KBSyncHandler{service: service}
}

// CreateSyncSource handles POST /api/sync-sources
func (h *KBSyncHandler) CreateSyncSource(c *gin.Context) {
	var req struct {
		ChatbotID        int64  `json:"chatbot_id" binding:"required"`
		SourceType       string `json:"source_type" binding:"required"`
		SourceIdentifier string `json:"source_identifier" binding:"required"`
		AuthToken        string `json:"auth_token"`
		SyncFrequency    string `json:"sync_frequency" binding:"required"`
		SyncSettings     string `json:"sync_settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syncSource := &models.KBSyncSource{
		ChatbotID:        req.ChatbotID,
		SourceType:       req.SourceType,
		SourceIdentifier: req.SourceIdentifier,
		SyncFrequency:    req.SyncFrequency,
		IsActive:         true,
	}

	if req.AuthToken != "" {
		syncSource.AuthToken = sql.NullString{String: req.AuthToken, Valid: true}
	}

	if req.SyncSettings != "" {
		syncSource.SyncSettings = sql.NullString{String: req.SyncSettings, Valid: true}
	}

	if err := h.service.CreateSyncSource(syncSource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, syncSource)
}

// GetSyncSource handles GET /api/sync-sources/:id
func (h *KBSyncHandler) GetSyncSource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	syncSource, err := h.service.GetSyncSource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if syncSource == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sync source not found"})
		return
	}

	c.JSON(http.StatusOK, syncSource)
}

// GetSyncSources handles GET /api/sync-sources/chatbot/:chatbot_id
func (h *KBSyncHandler) GetSyncSources(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	syncSources, err := h.service.GetSyncSources(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, syncSources)
}

// UpdateSyncSource handles PUT /api/sync-sources/:id
func (h *KBSyncHandler) UpdateSyncSource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	var req struct {
		SourceType       string `json:"source_type" binding:"required"`
		SourceIdentifier string `json:"source_identifier" binding:"required"`
		AuthToken        string `json:"auth_token"`
		SyncFrequency    string `json:"sync_frequency" binding:"required"`
		SyncSettings     string `json:"sync_settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syncSource := &models.KBSyncSource{
		ID:               id,
		SourceType:       req.SourceType,
		SourceIdentifier: req.SourceIdentifier,
		SyncFrequency:    req.SyncFrequency,
	}

	if req.AuthToken != "" {
		syncSource.AuthToken = sql.NullString{String: req.AuthToken, Valid: true}
	}

	if req.SyncSettings != "" {
		syncSource.SyncSettings = sql.NullString{String: req.SyncSettings, Valid: true}
	}

	if err := h.service.UpdateSyncSource(syncSource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, syncSource)
}

// DeleteSyncSource handles DELETE /api/sync-sources/:id
func (h *KBSyncHandler) DeleteSyncSource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	if err := h.service.DeleteSyncSource(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync source deleted successfully"})
}

// ToggleSyncSource handles POST /api/sync-sources/:id/toggle
func (h *KBSyncHandler) ToggleSyncSource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	var req struct {
		Active bool `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ToggleSyncSource(id, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := "activated"
	if !req.Active {
		status = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync source " + status + " successfully"})
}

// TriggerSync handles POST /api/sync-sources/:id/trigger
func (h *KBSyncHandler) TriggerSync(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	if err := h.service.TriggerManualSync(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync triggered successfully"})
}

// GetSyncStatus handles GET /api/sync-sources/:id/status
func (h *KBSyncHandler) GetSyncStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sync source ID"})
		return
	}

	status, err := h.service.GetSyncStatus(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
