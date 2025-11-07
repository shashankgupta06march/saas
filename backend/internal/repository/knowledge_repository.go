package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KnowledgeRepository struct {
	db *sql.DB
}

func NewKnowledgeRepository(db *sql.DB) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

func (r *KnowledgeRepository) Create(kb *models.KnowledgeBase) error {
	query := `INSERT INTO knowledge_base (organization_id, chatbot_id, title, content, content_type, source_url, embedding_vector) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, kb.OrganizationID, kb.ChatbotID, kb.Title, kb.Content, kb.ContentType, kb.SourceURL, kb.EmbeddingVector)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	kb.ID = id
	return nil
}

func (r *KnowledgeRepository) GetByChatbot(chatbotID int64) ([]models.KnowledgeBase, error) {
	query := `SELECT id, organization_id, chatbot_id, title, content, content_type, source_url, embedding_vector, created_at FROM knowledge_base WHERE chatbot_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.KnowledgeBase
	for rows.Next() {
		var kb models.KnowledgeBase
		var embedVector sql.NullString
		err := rows.Scan(
			&kb.ID,
			&kb.OrganizationID,
			&kb.ChatbotID,
			&kb.Title,
			&kb.Content,
			&kb.ContentType,
			&kb.SourceURL,
			&embedVector,
			&kb.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if embedVector.Valid {
			kb.EmbeddingVector = embedVector.String
		}
		items = append(items, kb)
	}

	return items, nil
}

func (r *KnowledgeRepository) GetAll(chatbotID int64) ([]models.KnowledgeBase, error) {
	query := `SELECT id, organization_id, chatbot_id, title, content, content_type, source_url, embedding_vector, created_at FROM knowledge_base WHERE chatbot_id = ?`
	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.KnowledgeBase
	for rows.Next() {
		var kb models.KnowledgeBase
		var embedVector sql.NullString
		err := rows.Scan(
			&kb.ID,
			&kb.OrganizationID,
			&kb.ChatbotID,
			&kb.Title,
			&kb.Content,
			&kb.ContentType,
			&kb.SourceURL,
			&embedVector,
			&kb.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if embedVector.Valid {
			kb.EmbeddingVector = embedVector.String
		}
		items = append(items, kb)
	}

	return items, nil
}

func (r *KnowledgeRepository) GetByID(id int64) (*models.KnowledgeBase, error) {
	query := `SELECT id, organization_id, chatbot_id, title, content, content_type, source_url, 
	          embedding_vector, version, status, created_at 
	          FROM knowledge_base WHERE id = ?`

	kb := &models.KnowledgeBase{}
	var embedVector sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&kb.ID,
		&kb.OrganizationID,
		&kb.ChatbotID,
		&kb.Title,
		&kb.Content,
		&kb.ContentType,
		&kb.SourceURL,
		&embedVector,
		&kb.Version,
		&kb.Status,
		&kb.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if embedVector.Valid {
		kb.EmbeddingVector = embedVector.String
	}

	return kb, nil
}

func (r *KnowledgeRepository) Update(kb *models.KnowledgeBase) error {
	query := `UPDATE knowledge_base 
	          SET title = ?, content = ?, content_type = ?, source_url = ?, 
	              embedding_vector = ?, version = ?, status = ?
	          WHERE id = ?`

	_, err := r.db.Exec(query,
		kb.Title, kb.Content, kb.ContentType, kb.SourceURL,
		kb.EmbeddingVector, kb.Version, kb.Status, kb.ID)

	return err
}

func (r *KnowledgeRepository) Delete(id int64) error {
	query := `DELETE FROM knowledge_base WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
