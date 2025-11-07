package handlers

import (
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBQualityHandler struct {
	service *services.KBQualityService
}

func NewKBQualityHandler(service *services.KBQualityService) *KBQualityHandler {
	return &KBQualityHandler{service: service}
}

// RunTest handles POST /api/quality/test
func (h *KBQualityHandler) RunTest(c *gin.Context) {
	var req struct {
		ChatbotID       int64  `json:"chatbot_id" binding:"required"`
		TestQuery       string `json:"test_query" binding:"required"`
		ExpectedContent string `json:"expected_content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	test, err := h.service.RunTest(req.ChatbotID, req.TestQuery, req.ExpectedContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

// RunBatchTests handles POST /api/quality/batch-test
func (h *KBQualityHandler) RunBatchTests(c *gin.Context) {
	var req struct {
		ChatbotID int64 `json:"chatbot_id" binding:"required"`
		Tests     []struct {
			Query    string `json:"query" binding:"required"`
			Expected string `json:"expected"`
		} `json:"tests" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert tests to map format
	var testCases []map[string]string
	for _, t := range req.Tests {
		testCases = append(testCases, map[string]string{
			"query":    t.Query,
			"expected": t.Expected,
		})
	}

	results, err := h.service.RunBatchTests(req.ChatbotID, testCases)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_tests": len(results),
		"results":     results,
	})
}

// GetTestHistory handles GET /api/quality/tests/:chatbot_id
func (h *KBQualityHandler) GetTestHistory(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	tests, err := h.service.GetTestHistory(chatbotID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tests)
}

// GetTestStats handles GET /api/quality/stats/:chatbot_id
func (h *KBQualityHandler) GetTestStats(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	stats, err := h.service.GetTestStats(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetFailedTests handles GET /api/quality/failed/:chatbot_id
func (h *KBQualityHandler) GetFailedTests(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	tests, err := h.service.GetFailedTests(chatbotID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tests)
}

// IdentifyGaps handles GET /api/quality/gaps/:chatbot_id
func (h *KBQualityHandler) IdentifyGaps(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	gaps, err := h.service.IdentifyGaps(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gaps": gaps,
	})
}

// GetQualityScore handles GET /api/quality/score/:chatbot_id
func (h *KBQualityHandler) GetQualityScore(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	score, err := h.service.GenerateQualityScore(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chatbot_id":    chatbotID,
		"quality_score": score,
		"rating":        getQualityRating(score),
	})
}

// GetSuggestions handles GET /api/quality/suggestions/:chatbot_id
func (h *KBQualityHandler) GetSuggestions(c *gin.Context) {
	chatbotID, err := strconv.ParseInt(c.Param("chatbot_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	suggestions, err := h.service.SuggestImprovements(chatbotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}

// Helper function to rate quality score
func getQualityRating(score float64) string {
	if score >= 0.9 {
		return "Excellent"
	} else if score >= 0.75 {
		return "Good"
	} else if score >= 0.6 {
		return "Fair"
	} else if score >= 0.4 {
		return "Poor"
	}
	return "Very Poor"
}
