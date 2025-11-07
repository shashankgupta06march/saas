package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBFolderHandler struct {
	service *services.KBFolderService
}

func NewKBFolderHandler(service *services.KBFolderService) *KBFolderHandler {
	return &KBFolderHandler{service: service}
}

// CreateFolder handles POST /api/folders
func (h *KBFolderHandler) CreateFolder(c *gin.Context) {
	var req struct {
		ChatbotID   int64  `json:"chatbot_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    *int64 `json:"parent_id"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder := &models.KBFolder{
		ChatbotID: req.ChatbotID,
		Name:      req.Name,
		Color:     req.Color,
		Icon:      req.Icon,
	}

	if req.Description != "" {
		folder.Description = sql.NullString{String: req.Description, Valid: true}
	}

	if req.ParentID != nil {
		folder.ParentID = sql.NullInt64{Int64: *req.ParentID, Valid: true}
	}

	if err := h.service.CreateFolder(folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

// GetFolder handles GET /api/folders/:id
func (h *KBFolderHandler) GetFolder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	folder, err := h.service.GetFolder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if folder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// UpdateFolder handles PUT /api/folders/:id
func (h *KBFolderHandler) UpdateFolder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    *int64 `json:"parent_id"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder := &models.KBFolder{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
		Icon:  req.Icon,
	}

	if req.Description != "" {
		folder.Description = sql.NullString{String: req.Description, Valid: true}
	}

	if req.ParentID != nil {
		folder.ParentID = sql.NullInt64{Int64: *req.ParentID, Valid: true}
	}

	if err := h.service.UpdateFolder(folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// DeleteFolder handles DELETE /api/folders/:id
func (h *KBFolderHandler) DeleteFolder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	if err := h.service.DeleteFolder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}

// GetFolderTree handles GET /api/folders/tree/:chatbot_id
func (h *KBFolderHandler) GetFolderTree(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	tree, err := h.service.GetFolderTree(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetRootFolders handles GET /api/folders/roots/:chatbot_id
func (h *KBFolderHandler) GetRootFolders(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	folders, err := h.service.GetRootFolders(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, folders)
}

// MoveFolder handles POST /api/folders/:id/move
func (h *KBFolderHandler) MoveFolder(c *gin.Context) {
	folderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req struct {
		NewParentID *int64 `json:"new_parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newParentID := int64(0)
	if req.NewParentID != nil {
		newParentID = *req.NewParentID
	}

	if err := h.service.MoveFolder(folderID, newParentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder moved successfully"})
}

// GetFolderPath handles GET /api/folders/:id/path
func (h *KBFolderHandler) GetFolderPath(c *gin.Context) {
	folderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	path, err := h.service.GetFolderPath(folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, path)
}
