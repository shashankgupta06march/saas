package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type WidgetHandler struct {
	widgetPath string
}

func NewWidgetHandler(widgetPath string) *WidgetHandler {
	return &WidgetHandler{widgetPath: widgetPath}
}

func (h *WidgetHandler) ServeWidget(c *gin.Context) {
	widgetFile := filepath.Join(h.widgetPath, "chatbot-widget.js")
	
	// Check if file exists
	if _, err := os.Stat(widgetFile); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Widget file not found"})
		return
	}

	// Set appropriate headers
	c.Header("Content-Type", "application/javascript")
	c.Header("Cache-Control", "public, max-age=3600")
	
	// Serve the file
	c.File(widgetFile)
}

