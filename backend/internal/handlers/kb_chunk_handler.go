package handlers

import (
	"net/http"
	"strconv"

	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KBChunkHandler struct {
	chunkingService *services.KBChunkingService
	kbRepo          *repository.KnowledgeRepository
}

func NewKBChunkHandler(chunkingService *services.KBChunkingService, kbRepo *repository.KnowledgeRepository) *KBChunkHandler {
	return &KBChunkHandler{
		chunkingService: chunkingService,
		kbRepo:          kbRepo,
	}
}

// GetChunks handles GET /api/knowledge/:kb_id/chunks
func (h *KBChunkHandler) GetChunks(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	chunks, err := h.chunkingService.GetChunksForKB(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chunks)
}

// RechunkKB handles POST /api/knowledge/:kb_id/rechunk
func (h *KBChunkHandler) RechunkKB(c *gin.Context) {
	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	// Get KB entry
	kb, err := h.kbRepo.GetByID(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if kb == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Knowledge base entry not found"})
		return
	}

	// Optional: Accept new chunking parameters
	var req struct {
		ChunkSize      *int    `json:"chunk_size"`
		ChunkOverlap   *int    `json:"chunk_overlap"`
		ChunkingMethod *string `json:"chunking_method"`
	}

	if err := c.ShouldBindJSON(&req); err == nil {
		// Update KB with new chunking parameters if provided
		if req.ChunkSize != nil {
			kb.ChunkSize = *req.ChunkSize
		}
		if req.ChunkOverlap != nil {
			kb.ChunkOverlap = *req.ChunkOverlap
		}
		if req.ChunkingMethod != nil {
			kb.ChunkingMethod = *req.ChunkingMethod
		}
	}

	// Rechunk the content
	if err := h.chunkingService.CreateChunksForKB(kb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the new chunks
	chunks, err := h.chunkingService.GetChunksForKB(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Knowledge base rechunked successfully",
		"chunk_count": len(chunks),
		"chunks":      chunks,
	})
}

// PreviewChunking handles POST /api/knowledge/preview-chunks
func (h *KBChunkHandler) PreviewChunking(c *gin.Context) {
	var req struct {
		Content        string `json:"content" binding:"required"`
		ChunkSize      int    `json:"chunk_size"`
		ChunkOverlap   int    `json:"chunk_overlap"`
		ChunkingMethod string `json:"chunking_method"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.ChunkSize == 0 {
		req.ChunkSize = 1000
	}
	if req.ChunkOverlap == 0 {
		req.ChunkOverlap = 200
	}
	if req.ChunkingMethod == "" {
		req.ChunkingMethod = "fixed"
	}

	// Preview chunks without saving
	chunks, err := h.chunkingService.ChunkContent(req.Content, req.ChunkingMethod, req.ChunkSize, req.ChunkOverlap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chunk_count":     len(chunks),
		"chunks":          chunks,
		"chunk_size":      req.ChunkSize,
		"chunk_overlap":   req.ChunkOverlap,
		"chunking_method": req.ChunkingMethod,
	})
}
