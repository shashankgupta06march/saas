package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/openai"
	openailib "github.com/sashabaranov/go-openai"
)

type KBQualityService struct {
	qualityRepo      *repository.KBQualityRepository
	knowledgeService *KnowledgeService
	openaiClient     *openai.Client
}

func NewKBQualityService(
	qualityRepo *repository.KBQualityRepository,
	knowledgeService *KnowledgeService,
	openaiClient *openai.Client,
) *KBQualityService {
	return &KBQualityService{
		qualityRepo:      qualityRepo,
		knowledgeService: knowledgeService,
		openaiClient:     openaiClient,
	}
}

// RunTest executes a quality test for a chatbot's knowledge base
func (s *KBQualityService) RunTest(chatbotID int64, testQuery string, expectedContent string) (*models.KBQualityTest, error) {
	// Get relevant context from knowledge base
	context, _, err := s.knowledgeService.GetRelevantContext(chatbotID, testQuery, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get relevant context: %w", err)
	}

	// Generate response using the knowledge base
	messages := []openailib.ChatCompletionMessage{
		{
			Role:    openailib.ChatMessageRoleUser,
			Content: testQuery,
		},
	}

	response, _, err := s.openaiClient.GenerateChatResponse(messages, context)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Calculate relevance score
	relevanceScore := s.calculateRelevance(response, expectedContent)

	// Determine if test passed (threshold: 0.6)
	passed := relevanceScore >= 0.6

	// Create test record
	test := &models.KBQualityTest{
		ChatbotID:       chatbotID,
		TestQuery:       testQuery,
		ExpectedContent: sql.NullString{String: expectedContent, Valid: true},
		ActualResponse:  sql.NullString{String: response, Valid: true},
		RelevanceScore:  sql.NullFloat64{Float64: relevanceScore, Valid: true},
		Passed:          sql.NullBool{Bool: passed, Valid: true},
	}

	if err := s.qualityRepo.Create(test); err != nil {
		return nil, fmt.Errorf("failed to save test: %w", err)
	}

	return test, nil
}

// calculateRelevance calculates how relevant the response is to expected content
func (s *KBQualityService) calculateRelevance(response, expected string) float64 {
	if expected == "" {
		// If no expected content, just check if we got a response
		if response != "" {
			return 0.8
		}
		return 0.0
	}

	// Normalize strings
	response = strings.ToLower(strings.TrimSpace(response))
	expected = strings.ToLower(strings.TrimSpace(expected))

	// Check if response contains expected content
	if strings.Contains(response, expected) {
		return 1.0
	}

	// Check for keyword overlap
	expectedWords := strings.Fields(expected)
	responseWords := strings.Fields(response)

	// Count matching keywords
	matches := 0
	for _, expectedWord := range expectedWords {
		if len(expectedWord) < 3 {
			continue // Skip short words
		}
		for _, responseWord := range responseWords {
			if expectedWord == responseWord {
				matches++
				break
			}
		}
	}

	if len(expectedWords) == 0 {
		return 0.0
	}

	// Calculate score based on keyword overlap
	score := float64(matches) / float64(len(expectedWords))
	return score
}

// RunBatchTests runs multiple tests
func (s *KBQualityService) RunBatchTests(chatbotID int64, tests []map[string]string) ([]models.KBQualityTest, error) {
	var results []models.KBQualityTest

	for _, testCase := range tests {
		query := testCase["query"]
		expected := testCase["expected"]

		test, err := s.RunTest(chatbotID, query, expected)
		if err != nil {
			// Log error but continue with other tests
			continue
		}

		results = append(results, *test)
	}

	return results, nil
}

// GetTestHistory retrieves test history for a chatbot
func (s *KBQualityService) GetTestHistory(chatbotID int64, limit int) ([]models.KBQualityTest, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.qualityRepo.GetByChatbot(chatbotID, limit)
}

// GetTestStats calculates statistics for a chatbot's tests
func (s *KBQualityService) GetTestStats(chatbotID int64) (map[string]interface{}, error) {
	return s.qualityRepo.GetStats(chatbotID)
}

// GetFailedTests retrieves failed tests
func (s *KBQualityService) GetFailedTests(chatbotID int64, limit int) ([]models.KBQualityTest, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.qualityRepo.GetFailedTests(chatbotID, limit)
}

// IdentifyGaps analyzes failed tests to identify knowledge gaps
func (s *KBQualityService) IdentifyGaps(chatbotID int64) ([]string, error) {
	failedTests, err := s.qualityRepo.GetFailedTests(chatbotID, 50)
	if err != nil {
		return nil, err
	}

	var gaps []string
	gapMap := make(map[string]bool)

	for _, test := range failedTests {
		// Extract keywords from failed test queries
		words := strings.Fields(strings.ToLower(test.TestQuery))
		for _, word := range words {
			if len(word) > 4 && !gapMap[word] {
				gapMap[word] = true
				gaps = append(gaps, word)
			}
		}
	}

	return gaps, nil
}

// GenerateQualityScore calculates an overall quality score for the knowledge base
func (s *KBQualityService) GenerateQualityScore(chatbotID int64) (float64, error) {
	stats, err := s.qualityRepo.GetStats(chatbotID)
	if err != nil {
		return 0.0, err
	}

	totalTests := stats["total_tests"].(int)
	if totalTests == 0 {
		return 0.0, nil
	}

	avgScore := stats["average_score"].(float64)
	passRate := stats["pass_rate"].(float64) / 100.0

	// Combined score: 70% based on average relevance score, 30% on pass rate
	qualityScore := (avgScore * 0.7) + (passRate * 0.3)

	return qualityScore, nil
}

// SuggestImprovements provides suggestions based on test results
func (s *KBQualityService) SuggestImprovements(chatbotID int64) ([]string, error) {
	stats, err := s.qualityRepo.GetStats(chatbotID)
	if err != nil {
		return nil, err
	}

	var suggestions []string

	totalTests := stats["total_tests"].(int)
	if totalTests == 0 {
		suggestions = append(suggestions, "No tests have been run yet. Start by creating test queries.")
		return suggestions, nil
	}

	passRate := stats["pass_rate"].(float64)
	avgScore := stats["average_score"].(float64)

	if passRate < 50 {
		suggestions = append(suggestions, "Pass rate is low. Review failed tests and add missing content to knowledge base.")
	}

	if avgScore < 0.5 {
		suggestions = append(suggestions, "Average relevance score is low. Improve content quality and organization.")
	}

	if passRate >= 50 && passRate < 75 {
		suggestions = append(suggestions, "Consider adding more detailed information for common queries.")
	}

	if avgScore >= 0.7 && passRate >= 75 {
		suggestions = append(suggestions, "Knowledge base quality is good! Continue monitoring and testing.")
	}

	// Identify gaps
	gaps, err := s.IdentifyGaps(chatbotID)
	if err == nil && len(gaps) > 0 {
		topGaps := gaps
		if len(topGaps) > 5 {
			topGaps = topGaps[:5]
		}
		suggestions = append(suggestions, fmt.Sprintf("Common topics in failed tests: %s", strings.Join(topGaps, ", ")))
	}

	return suggestions, nil
}


