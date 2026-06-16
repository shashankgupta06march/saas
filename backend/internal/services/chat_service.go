package services

import (
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
	"github.com/chatbot-saas/backend/pkg/openai"
	openailib "github.com/sashabaranov/go-openai"
)

type ChatService struct {
	convRepo         *repository.ConversationRepository
	knowledgeService *KnowledgeService
	chunkingService  *KBChunkingService
	openaiClient     *openai.Client
}

func NewChatService(convRepo *repository.ConversationRepository, knowledgeService *KnowledgeService, chunkingService *KBChunkingService, openaiClient *openai.Client) *ChatService {
	return &ChatService{
		convRepo:         convRepo,
		knowledgeService: knowledgeService,
		chunkingService:  chunkingService,
		openaiClient:     openaiClient,
	}
}

// getRelevantContext prefers chunk-level retrieval (precise matches within
// large pages) and falls back to whole-page retrieval for chatbots whose
// knowledge base entries haven't been chunked yet.
func (s *ChatService) getRelevantContext(chatbotID int64, message string) string {
	const chunkTopK = 5

	chunks, err := s.chunkingService.GetRelevantChunks(chatbotID, message, chunkTopK)
	if err == nil && len(chunks) > 0 {
		var parts []string
		for _, c := range chunks {
			parts = append(parts, c.Content)
		}
		return strings.Join(parts, "\n\n")
	}

	context, _, err := s.knowledgeService.GetRelevantContext(chatbotID, message, 3)
	if err != nil {
		return ""
	}
	return context
}

func (s *ChatService) HandleMessage(chatbotID int64, sessionID, message string) (string, error) {
	// Get or create conversation
	conv, err := s.convRepo.GetBySessionID(sessionID)
	if err != nil {
		return "", err
	}

	if conv == nil {
		conv = &models.Conversation{
			ChatbotID: chatbotID,
			SessionID: sessionID,
			VisitorID: sessionID,
		}
		err = s.convRepo.Create(conv)
		if err != nil {
			return "", err
		}
	}

	// Save user message
	userMsg := &models.Message{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        message,
	}
	err = s.convRepo.CreateMessage(userMsg)
	if err != nil {
		return "", err
	}

	// Get conversation history
	messages, err := s.convRepo.GetMessages(conv.ID)
	if err != nil {
		return "", err
	}

	// Convert to OpenAI format (exclude the last message we just added, we'll add it separately)
	var chatMessages []openailib.ChatCompletionMessage
	for i := 0; i < len(messages)-1; i++ {
		role := openailib.ChatMessageRoleUser
		if messages[i].Role == "assistant" {
			role = openailib.ChatMessageRoleAssistant
		}
		chatMessages = append(chatMessages, openailib.ChatCompletionMessage{
			Role:    role,
			Content: messages[i].Content,
		})
	}

	// Add current message
	chatMessages = append(chatMessages, openailib.ChatCompletionMessage{
		Role:    openailib.ChatMessageRoleUser,
		Content: message,
	})

	// Get relevant context from knowledge base
	context := s.getRelevantContext(chatbotID, message)

	// Generate response
	response, _, err := s.openaiClient.GenerateChatResponse(chatMessages, context)
	if err != nil {
		return "", err
	}

	// Save assistant message
	assistantMsg := &models.Message{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        response,
	}
	err = s.convRepo.CreateMessage(assistantMsg)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (s *ChatService) GetConversations(chatbotID int64, limit int) ([]models.Conversation, error) {
	return s.convRepo.GetByChatbot(chatbotID, limit)
}

func (s *ChatService) GetMessages(conversationID int64) ([]models.Message, error) {
	return s.convRepo.GetMessages(conversationID)
}
