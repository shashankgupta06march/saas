package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBChunkRepository struct {
	db *sql.DB
}

func NewKBChunkRepository(db *sql.DB) *KBChunkRepository {
	return &KBChunkRepository{db: db}
}

func (r *KBChunkRepository) Create(chunk *models.KBChunk) error {
	query := `INSERT INTO kb_chunks (kb_id, chunk_index, content, embedding_vector, token_count, metadata) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, chunk.KBID, chunk.ChunkIndex, chunk.Content,
		chunk.EmbeddingVector, chunk.TokenCount, chunk.Metadata)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	chunk.ID = id
	return nil
}

func (r *KBChunkRepository) GetByKBID(kbID int64) ([]models.KBChunk, error) {
	query := `SELECT id, kb_id, chunk_index, content, embedding_vector, token_count, metadata, created_at 
	          FROM kb_chunks WHERE kb_id = ? ORDER BY chunk_index ASC`

	rows, err := r.db.Query(query, kbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []models.KBChunk
	for rows.Next() {
		var chunk models.KBChunk
		err := rows.Scan(
			&chunk.ID, &chunk.KBID, &chunk.ChunkIndex, &chunk.Content,
			&chunk.EmbeddingVector, &chunk.TokenCount, &chunk.Metadata, &chunk.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (r *KBChunkRepository) GetByID(id int64) (*models.KBChunk, error) {
	query := `SELECT id, kb_id, chunk_index, content, embedding_vector, token_count, metadata, created_at 
	          FROM kb_chunks WHERE id = ?`

	chunk := &models.KBChunk{}
	err := r.db.QueryRow(query, id).Scan(
		&chunk.ID, &chunk.KBID, &chunk.ChunkIndex, &chunk.Content,
		&chunk.EmbeddingVector, &chunk.TokenCount, &chunk.Metadata, &chunk.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

func (r *KBChunkRepository) DeleteByKBID(kbID int64) error {
	query := `DELETE FROM kb_chunks WHERE kb_id = ?`
	_, err := r.db.Exec(query, kbID)
	return err
}

func (r *KBChunkRepository) GetAllChunks(chatbotID int64) ([]models.KBChunk, error) {
	query := `SELECT c.id, c.kb_id, c.chunk_index, c.content, c.embedding_vector, c.token_count, c.metadata, c.created_at 
	          FROM kb_chunks c
	          INNER JOIN knowledge_base kb ON c.kb_id = kb.id
	          WHERE kb.chatbot_id = ? AND kb.status = 'active'
	          ORDER BY c.kb_id, c.chunk_index`

	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []models.KBChunk
	for rows.Next() {
		var chunk models.KBChunk
		err := rows.Scan(
			&chunk.ID, &chunk.KBID, &chunk.ChunkIndex, &chunk.Content,
			&chunk.EmbeddingVector, &chunk.TokenCount, &chunk.Metadata, &chunk.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (r *KBChunkRepository) UpdateEmbedding(chunkID int64, embedding string) error {
	query := `UPDATE kb_chunks SET embedding_vector = ? WHERE id = ?`
	_, err := r.db.Exec(query, embedding, chunkID)
	return err
}

func (r *KBChunkRepository) CountByKB(kbID int64) (int, error) {
	query := `SELECT COUNT(*) FROM kb_chunks WHERE kb_id = ?`

	var count int
	err := r.db.QueryRow(query, kbID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}


