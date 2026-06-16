package services

import (
	"errors"
	"log"
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/openai"
)

type KnowledgeService struct {
	repo            *repository.KnowledgeRepository
	openaiClient    *openai.Client
	chunkingService *KBChunkingService
}

func NewKnowledgeService(repo *repository.KnowledgeRepository, openaiClient *openai.Client, chunkingService *KBChunkingService) *KnowledgeService {
	return &KnowledgeService{
		repo:            repo,
		openaiClient:    openaiClient,
		chunkingService: chunkingService,
	}
}

// maxEmbeddableChars keeps content safely under text-embedding-3-small's
// 8191-token input limit (~32-35K chars for English text) — content beyond
// this would make the OpenAI embedding call fail outright.
const maxEmbeddableChars = 24000

func (s *KnowledgeService) AddKnowledge(kb *models.KnowledgeBase) error {
	if len(kb.Content) > maxEmbeddableChars {
		kb.Content = kb.Content[:maxEmbeddableChars]
	}

	// Generate embedding for the content
	embedding, err := s.openaiClient.GenerateEmbedding(kb.Content)
	if err != nil {
		return err
	}

	// Convert embedding to JSON
	embeddingJSON, err := openai.EmbeddingToJSON(embedding)
	if err != nil {
		return err
	}

	kb.EmbeddingVector = embeddingJSON

	if err := s.repo.Create(kb); err != nil {
		return err
	}

	// Chunk the content so retrieval can match specific facts within large
	// pages instead of only ever comparing against one whole-page embedding.
	// A chunking failure shouldn't fail the whole add — whole-page retrieval
	// still works as a fallback.
	if err := s.chunkingService.CreateChunksForKB(kb); err != nil {
		log.Printf("warning: failed to chunk knowledge base entry %d: %v", kb.ID, err)
	}

	return nil
}

type SourceInfo struct {
	Title string
	URL   string
	Type  string
}

func (s *KnowledgeService) GetRelevantContext(chatbotID int64, query string, topK int) (string, []SourceInfo, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.openaiClient.GenerateEmbedding(query)
	if err != nil {
		return "", nil, err
	}

	// Get all knowledge base entries for this chatbot
	allKB, err := s.repo.GetAll(chatbotID)
	if err != nil {
		return "", nil, err
	}

	if len(allKB) == 0 {
		return "", nil, nil
	}

	// Calculate similarity scores
	type scoredKB struct {
		kb    models.KnowledgeBase
		score float64
	}

	var scored []scoredKB
	for _, kb := range allKB {
		if kb.EmbeddingVector == "" {
			continue
		}

		embedding, err := openai.ParseEmbedding(kb.EmbeddingVector)
		if err != nil {
			continue
		}

		similarity := openai.CosineSimilarity(queryEmbedding, embedding)
		scored = append(scored, scoredKB{kb: kb, score: similarity})
	}

	if len(scored) == 0 {
		return "", nil, nil
	}

	// Sort by score (simple bubble sort for small datasets)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Get top K results
	if topK > len(scored) {
		topK = len(scored)
	}

	// Cap each snippet so a handful of large pages (e.g. a multi-page site
	// crawl) can't blow past the chat model's context window.
	const maxContextCharsPerEntry = 3000

	var contextParts []string
	var sources []SourceInfo

	for i := 0; i < topK; i++ {
		content := scored[i].kb.Content
		if len(content) > maxContextCharsPerEntry {
			content = content[:maxContextCharsPerEntry]
		}
		contextParts = append(contextParts, content)

		// Add source information if available
		if scored[i].kb.SourceURL != "" {
			sources = append(sources, SourceInfo{
				Title: scored[i].kb.Title,
				URL:   scored[i].kb.SourceURL,
				Type:  scored[i].kb.ContentType,
			})
		}
	}

	return strings.Join(contextParts, "\n\n"), sources, nil
}

func (s *KnowledgeService) GetByChatbot(chatbotID int64) ([]models.KnowledgeBase, error) {
	return s.repo.GetByChatbot(chatbotID)
}

func (s *KnowledgeService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid knowledge base ID")
	}
	return s.repo.Delete(id)
}
