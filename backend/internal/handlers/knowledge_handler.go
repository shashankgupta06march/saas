package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/chatbot-saas/backend/pkg/parser"
	"github.com/gin-gonic/gin"
)

type KnowledgeHandler struct {
	service *services.KnowledgeService
}

func NewKnowledgeHandler(service *services.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{service: service}
}

type AddKnowledgeRequest struct {
	ChatbotID   int64  `json:"chatbot_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	ContentType string `json:"content_type"`
	SourceURL   string `json:"source_url"`
}

func (h *KnowledgeHandler) Add(c *gin.Context) {
	var req AddKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organizationID := c.GetInt64("organization_id")

	kb := &models.KnowledgeBase{
		OrganizationID: organizationID,
		ChatbotID:      req.ChatbotID,
		Title:          req.Title,
		Content:        req.Content,
		ContentType:    req.ContentType,
		SourceURL:      req.SourceURL,
	}

	err := h.service.AddKnowledge(kb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add knowledge: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, kb)
}

func (h *KnowledgeHandler) GetByChatbot(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	items, err := h.service.GetByChatbot(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get knowledge base"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *KnowledgeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete knowledge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge deleted successfully"})
}

// UploadFile handles file uploads (PDF, DOCX, TXT)
func (h *KnowledgeHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	chatbotIDStr := c.PostForm("chatbot_id")
	chatbotID, err := strconv.ParseInt(chatbotIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	organizationID := c.GetInt64("organization_id")

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Determine file type and parse accordingly
	ext := strings.ToLower(filepath.Ext(header.Filename))
	var content string
	var contentType string

	switch ext {
	case ".pdf":
		reader := bytes.NewReader(fileData)
		content, err = parser.ParsePDF(reader, int64(len(fileData)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse PDF: " + err.Error()})
			return
		}
		contentType = "pdf"

	case ".docx":
		reader := bytes.NewReader(fileData)
		content, err = parser.ParseDOCX(reader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse DOCX: " + err.Error()})
			return
		}
		contentType = "docx"

	case ".txt", ".text":
		reader := bytes.NewReader(fileData)
		content, err = parser.ParsePlainText(reader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse text file: " + err.Error()})
			return
		}
		contentType = "text"

	case ".xlsx", ".xls":
		reader := bytes.NewReader(fileData)
		content, err = parser.ParseExcel(reader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Excel file: " + err.Error()})
			return
		}
		contentType = "xlsx"

	case ".csv":
		reader := bytes.NewReader(fileData)
		content, err = parser.ParseCSV(reader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV file: " + err.Error()})
			return
		}
		contentType = "csv"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported file type: %s. Supported types: PDF, DOCX, TXT, XLSX, XLS, CSV", ext)})
		return
	}

	// Limit content size
	if len(content) > 50000 {
		content = content[:50000]
	}

	// Create knowledge base entry
	kb := &models.KnowledgeBase{
		OrganizationID: organizationID,
		ChatbotID:      chatbotID,
		Title:          title,
		Content:        content,
		ContentType:    contentType,
		SourceURL:      header.Filename,
	}

	err = h.service.AddKnowledge(kb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add knowledge: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "File uploaded and processed successfully",
		"knowledge": kb,
	})
}

// ScrapeURL handles website URL scraping
type ScrapeURLRequest struct {
	ChatbotID int64  `json:"chatbot_id" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Title     string `json:"title"`
	Depth     int    `json:"depth"`
}

func (h *KnowledgeHandler) ScrapeURL(c *gin.Context) {
	var req ScrapeURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate URL
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL must start with http:// or https://"})
		return
	}

	organizationID := c.GetInt64("organization_id")

	// Set default depth if not provided
	depth := req.Depth
	if depth < 0 {
		depth = 0
	}
	if depth > 5 {
		depth = 5
	}

	// Scrape the website
	content, err := parser.ScrapeWebsite(req.URL, depth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape website: " + err.Error()})
		return
	}

	// Limit content size
	if len(content) > 50000 {
		content = content[:50000]
	}

	// Use URL as title if not provided
	title := req.Title
	if title == "" {
		title = req.URL
	}

	// Create knowledge base entry
	kb := &models.KnowledgeBase{
		OrganizationID: organizationID,
		ChatbotID:      req.ChatbotID,
		Title:          title,
		Content:        content,
		ContentType:    "webpage",
		SourceURL:      req.URL,
	}

	err = h.service.AddKnowledge(kb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add knowledge: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Website scraped and processed successfully",
		"knowledge": kb,
	})
}
