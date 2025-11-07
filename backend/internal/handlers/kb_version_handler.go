package handlers

import (
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBVersionHandler struct {
	service *services.KBVersionService
}

func NewKBVersionHandler(service *services.KBVersionService) *KBVersionHandler {
	return &KBVersionHandler{service: service}
}

// GetVersionHistory handles GET /api/kb/:id/versions
func (h *KBVersionHandler) GetVersionHistory(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	versions, err := h.service.GetVersionHistory(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// GetVersion handles GET /api/knowledge/:kb_id/versions/:version
func (h *KBVersionHandler) GetVersion(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
		return
	}

	versionData, err := h.service.GetVersion(kbID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if versionData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	c.JSON(http.StatusOK, versionData)
}

// RestoreVersion handles POST /api/knowledge/:kb_id/versions/:version/restore
func (h *KBVersionHandler) RestoreVersion(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := c.GetInt64("user_id")

	if err := h.service.RestoreVersion(kbID, version, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Version restored successfully"})
}

// CompareVersions handles GET /api/knowledge/:kb_id/versions/compare
func (h *KBVersionHandler) CompareVersions(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	version1Str := c.Query("version1")
	version2Str := c.Query("version2")

	if version1Str == "" || version2Str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both version1 and version2 query parameters are required"})
		return
	}

	version1, err := strconv.Atoi(version1Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version1"})
		return
	}

	version2, err := strconv.Atoi(version2Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version2"})
		return
	}

	comparison, err := h.service.CompareVersions(kbID, version1, version2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparison)
}
