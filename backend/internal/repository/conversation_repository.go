package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	query := `INSERT INTO conversations (chatbot_id, session_id, visitor_id) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, conv.ChatbotID, conv.SessionID, conv.VisitorID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	conv.ID = id
	return nil
}

func (r *ConversationRepository) GetBySessionID(sessionID string) (*models.Conversation, error) {
	conv := &models.Conversation{}
	query := `SELECT id, chatbot_id, session_id, visitor_id, started_at, ended_at FROM conversations WHERE session_id = ? ORDER BY started_at DESC LIMIT 1`
	err := r.db.QueryRow(query, sessionID).Scan(
		&conv.ID,
		&conv.ChatbotID,
		&conv.SessionID,
		&conv.VisitorID,
		&conv.StartedAt,
		&conv.EndedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return conv, err
}

func (r *ConversationRepository) GetByChatbot(chatbotID int64, limit int) ([]models.Conversation, error) {
	query := `SELECT id, chatbot_id, session_id, visitor_id, started_at, ended_at FROM conversations WHERE chatbot_id = ? ORDER BY started_at DESC LIMIT ?`
	rows, err := r.db.Query(query, chatbotID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.ChatbotID,
			&conv.SessionID,
			&conv.VisitorID,
			&conv.StartedAt,
			&conv.EndedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *ConversationRepository) CreateMessage(msg *models.Message) error {
	query := `INSERT INTO messages (conversation_id, role, content) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, msg.ConversationID, msg.Role, msg.Content)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	msg.ID = id
	return nil
}

func (r *ConversationRepository) GetMessages(conversationID int64) ([]models.Message, error) {
	query := `SELECT id, conversation_id, role, content, timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp ASC`
	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.Role,
			&msg.Content,
			&msg.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

