package services

import (
	"errors"
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/openai"
)

type KnowledgeService struct {
	repo         *repository.KnowledgeRepository
	openaiClient *openai.Client
}

func NewKnowledgeService(repo *repository.KnowledgeRepository, openaiClient *openai.Client) *KnowledgeService {
	return &KnowledgeService{
		repo:         repo,
		openaiClient: openaiClient,
	}
}

func (s *KnowledgeService) AddKnowledge(kb *models.KnowledgeBase) error {
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

	return s.repo.Create(kb)
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

	var contextParts []string
	var sources []SourceInfo

	for i := 0; i < topK; i++ {
		contextParts = append(contextParts, scored[i].kb.Content)

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
