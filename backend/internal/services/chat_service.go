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
	chatbotRepo      *repository.ChatbotRepository
	openaiClient     *openai.Client
}

func NewChatService(convRepo *repository.ConversationRepository, knowledgeService *KnowledgeService, chunkingService *KBChunkingService, chatbotRepo *repository.ChatbotRepository, openaiClient *openai.Client) *ChatService {
	return &ChatService{
		convRepo:         convRepo,
		knowledgeService: knowledgeService,
		chunkingService:  chunkingService,
		chatbotRepo:      chatbotRepo,
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

// noAnswerToken is what we instruct the model to emit when the retrieved
// context can't answer the user's question. Embedding similarity alone can't
// distinguish "tell me about X college" from a different, on-file college, so
// we let the model itself judge whether the context actually answers the
// question and map that judgment onto the admin's fallback message.
const noAnswerToken = "NO_ANSWER"

// guardedSystemPrompt builds a context-only system prompt that asks the model
// to emit noAnswerToken when the context doesn't cover the question. Used only
// when the admin has configured a fallback message.
func guardedSystemPrompt(context string) string {
	var b strings.Builder
	b.WriteString("You are a helpful assistant for a specific organization. ")
	b.WriteString("Answer the user ONLY using the context provided below. ")
	b.WriteString("Do not use outside or general knowledge to answer factual questions.\n\n")
	if strings.TrimSpace(context) != "" {
		b.WriteString("Context:\n")
		b.WriteString(context)
		b.WriteString("\n\n")
	} else {
		b.WriteString("Context: (none provided)\n\n")
	}
	b.WriteString("Rules:\n")
	b.WriteString("- If the user's message is only a greeting, thanks, or small talk, reply briefly and politely and do NOT output the token below.\n")
	b.WriteString("- If the user asks a question and the context above does not clearly contain the information needed to answer it, reply with exactly this and nothing else: " + noAnswerToken + "\n")
	b.WriteString("- Never invent facts that are not present in the context.")
	return b.String()
}

// isNoAnswerResponse reports whether the model signalled it couldn't answer
// from the context.
func isNoAnswerResponse(resp string) bool {
	return strings.Contains(strings.ToUpper(strings.TrimSpace(resp)), noAnswerToken)
}

// fallbackMessage returns the admin-configured "nothing found" message for the
// chatbot, or an empty string if none is set.
func (s *ChatService) fallbackMessage(chatbotID int64) string {
	if s.chatbotRepo == nil {
		return ""
	}
	settings, err := s.chatbotRepo.GetSettings(chatbotID)
	if err != nil || settings == nil {
		return ""
	}
	return strings.TrimSpace(settings.FallbackMessage)
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

	// When the admin has configured a fallback message, answer strictly from
	// the knowledge base context and let the model signal (via noAnswerToken)
	// when the context can't answer the question — then show the fallback.
	// Greetings/small talk still get a normal reply. When no fallback is
	// configured we keep the original behavior.
	var response string
	if fallback := s.fallbackMessage(chatbotID); fallback != "" {
		resp, _, err := s.openaiClient.GenerateChatResponseWithSystem(chatMessages, guardedSystemPrompt(context))
		if err != nil {
			return "", err
		}
		if isNoAnswerResponse(resp) {
			response = fallback
		} else {
			response = resp
		}
	} else {
		resp, _, err := s.openaiClient.GenerateChatResponse(chatMessages, context)
		if err != nil {
			return "", err
		}
		response = resp
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
