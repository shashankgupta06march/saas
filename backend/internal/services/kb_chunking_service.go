package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/openai"
)

type KBChunkingService struct {
	chunkRepo    *repository.KBChunkRepository
	openaiClient *openai.Client
}

func NewKBChunkingService(chunkRepo *repository.KBChunkRepository, openaiClient *openai.Client) *KBChunkingService {
	return &KBChunkingService{
		chunkRepo:    chunkRepo,
		openaiClient: openaiClient,
	}
}

// ChunkContent splits content into chunks based on the specified method
func (s *KBChunkingService) ChunkContent(content string, method string, chunkSize, overlap int) ([]string, error) {
	switch method {
	case "fixed":
		return s.fixedChunking(content, chunkSize, overlap), nil
	case "sentence":
		return s.sentenceChunking(content, chunkSize), nil
	case "semantic":
		return s.semanticChunking(content, chunkSize), nil
	default:
		return s.fixedChunking(content, chunkSize, overlap), nil
	}
}

// fixedChunking splits content into fixed-size chunks with overlap
func (s *KBChunkingService) fixedChunking(content string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		chunkSize = 1000
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 4 // Max 25% overlap
	}

	var chunks []string
	contentLen := utf8.RuneCountInString(content)

	for start := 0; start < contentLen; {
		end := start + chunkSize
		if end > contentLen {
			end = contentLen
		}

		// Convert rune positions to byte positions
		chunk := string([]rune(content)[start:end])
		chunks = append(chunks, strings.TrimSpace(chunk))

		// Move start position with overlap
		start = end - overlap
		if start <= 0 {
			start = end
		}

		// Prevent infinite loop
		if start >= contentLen {
			break
		}
	}

	return chunks
}

// sentenceChunking splits content by sentences, grouping them to reach chunkSize
func (s *KBChunkingService) sentenceChunking(content string, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = 1000
	}

	// Split by common sentence endings
	sentences := s.splitIntoSentences(content)

	var chunks []string
	var currentChunk strings.Builder

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// If adding this sentence would exceed chunk size and we have content
		if currentChunk.Len() > 0 && currentChunk.Len()+len(sentence) > chunkSize {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(sentence)
	}

	// Add remaining content
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// semanticChunking attempts to split content at paragraph boundaries
func (s *KBChunkingService) semanticChunking(content string, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = 1000
	}

	// Split by paragraphs (double newlines)
	paragraphs := strings.Split(content, "\n\n")

	var chunks []string
	var currentChunk strings.Builder

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// If paragraph alone exceeds chunk size, fall back to sentence chunking
		if len(paragraph) > chunkSize*2 {
			sentenceChunks := s.sentenceChunking(paragraph, chunkSize)
			chunks = append(chunks, sentenceChunks...)
			continue
		}

		// If adding this paragraph would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+len(paragraph) > chunkSize {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(paragraph)
	}

	// Add remaining content
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// splitIntoSentences splits text into sentences
func (s *KBChunkingService) splitIntoSentences(text string) []string {
	// Simple sentence splitting by common endings
	replacer := strings.NewReplacer(
		". ", ".|",
		"! ", "!|",
		"? ", "?|",
		".\n", ".|",
		"!\n", "!|",
		"?\n", "?|",
	)
	marked := replacer.Replace(text)
	sentences := strings.Split(marked, "|")

	var result []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

// CreateChunksForKB creates and stores chunks for a knowledge base entry
func (s *KBChunkingService) CreateChunksForKB(kb *models.KnowledgeBase) error {
	// Delete existing chunks
	if err := s.chunkRepo.DeleteByKBID(kb.ID); err != nil {
		return fmt.Errorf("failed to delete existing chunks: %w", err)
	}

	// Split content into chunks
	chunks, err := s.ChunkContent(kb.Content, kb.ChunkingMethod, kb.ChunkSize, kb.ChunkOverlap)
	if err != nil {
		return fmt.Errorf("failed to chunk content: %w", err)
	}

	// Create and store chunks with embeddings
	for i, chunkContent := range chunks {
		// Generate embedding for chunk
		embedding, err := s.openaiClient.GenerateEmbedding(chunkContent)
		if err != nil {
			return fmt.Errorf("failed to generate embedding for chunk %d: %w", i, err)
		}

		embeddingJSON, err := openai.EmbeddingToJSON(embedding)
		if err != nil {
			return fmt.Errorf("failed to convert embedding to JSON: %w", err)
		}

		// Create metadata for chunk
		metadata := map[string]interface{}{
			"original_kb_id": kb.ID,
			"chunk_position": i,
			"total_chunks":   len(chunks),
		}
		metadataJSON, _ := json.Marshal(metadata)

		chunk := &models.KBChunk{
			KBID:            kb.ID,
			ChunkIndex:      i,
			Content:         chunkContent,
			EmbeddingVector: sql.NullString{String: embeddingJSON, Valid: true},
			TokenCount:      sql.NullInt32{Int32: int32(len(strings.Fields(chunkContent))), Valid: true},
			Metadata:        sql.NullString{String: string(metadataJSON), Valid: true},
		}

		if err := s.chunkRepo.Create(chunk); err != nil {
			return fmt.Errorf("failed to create chunk %d: %w", i, err)
		}
	}

	return nil
}

// GetRelevantChunks finds the most relevant chunks for a query
func (s *KBChunkingService) GetRelevantChunks(chatbotID int64, query string, topK int) ([]models.KBChunk, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.openaiClient.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Get all chunks for this chatbot
	allChunks, err := s.chunkRepo.GetAllChunks(chatbotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks: %w", err)
	}

	if len(allChunks) == 0 {
		return []models.KBChunk{}, nil
	}

	// Calculate similarity scores
	type scoredChunk struct {
		chunk models.KBChunk
		score float64
	}

	var scored []scoredChunk
	for _, chunk := range allChunks {
		if !chunk.EmbeddingVector.Valid || chunk.EmbeddingVector.String == "" {
			continue
		}

		embedding, err := openai.ParseEmbedding(chunk.EmbeddingVector.String)
		if err != nil {
			continue
		}

		similarity := openai.CosineSimilarity(queryEmbedding, embedding)
		scored = append(scored, scoredChunk{chunk: chunk, score: similarity})
	}

	if len(scored) == 0 {
		return []models.KBChunk{}, nil
	}

	// Sort by score (bubble sort for simplicity)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Return top K
	if topK > len(scored) {
		topK = len(scored)
	}

	var result []models.KBChunk
	for i := 0; i < topK; i++ {
		result = append(result, scored[i].chunk)
	}

	return result, nil
}

// GetChunksForKB retrieves all chunks for a knowledge base entry
func (s *KBChunkingService) GetChunksForKB(kbID int64) ([]models.KBChunk, error) {
	return s.chunkRepo.GetByKBID(kbID)
}

