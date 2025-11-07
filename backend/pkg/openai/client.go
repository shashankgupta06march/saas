package openai

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
	}
}

func (c *Client) GenerateEmbedding(text string) ([]float64, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	resp, err := c.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequestStrings{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embedding generated")
	}

	// Convert []float32 to []float64
	embedding32 := resp.Data[0].Embedding
	embedding64 := make([]float64, len(embedding32))
	for i, v := range embedding32 {
		embedding64[i] = float64(v)
	}

	return embedding64, nil
}

func (c *Client) GenerateChatResponse(messages []openai.ChatCompletionMessage, knowledgeContext string) (string, int, error) {
	// Prepend knowledge base context if available
	if knowledgeContext != "" {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant. Use the following context to answer questions:\n\n" + knowledgeContext,
		}
		messages = append([]openai.ChatCompletionMessage{systemMessage}, messages...)
	} else {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		}
		messages = append([]openai.ChatCompletionMessage{systemMessage}, messages...)
	}

	resp, err := c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})

	if err != nil {
		return "", 0, err
	}

	if len(resp.Choices) == 0 {
		return "", 0, errors.New("no response generated")
	}

	return resp.Choices[0].Message.Content, resp.Usage.TotalTokens, nil
}

// Cosine similarity between two vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Parse embedding vector from JSON string
func ParseEmbedding(embeddingJSON string) ([]float64, error) {
	var embedding []float64
	err := json.Unmarshal([]byte(embeddingJSON), &embedding)
	return embedding, err
}

// Convert embedding to JSON string
func EmbeddingToJSON(embedding []float64) (string, error) {
	bytes, err := json.Marshal(embedding)
	return string(bytes), err
}

