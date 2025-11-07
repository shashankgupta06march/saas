package handlers

import (
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBTagHandler struct {
	service *services.KBTagService
}

func NewKBTagHandler(service *services.KBTagService) *KBTagHandler {
	return &KBTagHandler{service: service}
}

// CreateTag handles POST /api/tags
func (h *KBTagHandler) CreateTag(c *gin.Context) {
	var req struct {
		ChatbotID int64  `json:"chatbot_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Color     string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.KBTag{
		ChatbotID: req.ChatbotID,
		Name:      req.Name,
		Color:     req.Color,
	}

	if err := h.service.CreateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTag handles GET /api/tags/:id
func (h *KBTagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := h.service.GetTag(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tag == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// GetAllTags handles GET /api/tags/chatbot/:chatbot_id
func (h *KBTagHandler) GetAllTags(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	tags, err := h.service.GetAllTags(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag handles PUT /api/tags/:id
func (h *KBTagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.KBTag{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
	}

	if err := h.service.UpdateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag handles DELETE /api/tags/:id
func (h *KBTagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.service.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

// AssignTags handles POST /api/kb/:id/tags
func (h *KBTagHandler) AssignTags(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	var req struct {
		TagIDs []int64 `json:"tag_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignTags(kbID, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags assigned successfully"})
}

// GetKBTags handles GET /api/kb/:id/tags
func (h *KBTagHandler) GetKBTags(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	tags, err := h.service.GetKBTags(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// AddTagToKB handles POST /api/kb/:id/tags/:tag_id
func (h *KBTagHandler) AddTagToKB(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	tagID, err := strconv.ParseInt(c.Param("tag_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.service.AddTag(kbID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added successfully"})
}

// RemoveTagFromKB handles DELETE /api/kb/:id/tags/:tag_id
func (h *KBTagHandler) RemoveTagFromKB(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	tagID, err := strconv.ParseInt(c.Param("tag_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.service.RemoveTag(kbID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed successfully"})
}
