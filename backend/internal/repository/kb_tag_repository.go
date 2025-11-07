package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBTagRepository struct {
	db *sql.DB
}

func NewKBTagRepository(db *sql.DB) *KBTagRepository {
	return &KBTagRepository{db: db}
}

func (r *KBTagRepository) Create(tag *models.KBTag) error {
	query := `INSERT INTO kb_tags (chatbot_id, name, color) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, tag.ChatbotID, tag.Name, tag.Color)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tag.ID = id
	return nil
}

func (r *KBTagRepository) GetByID(id int64) (*models.KBTag, error) {
	query := `SELECT id, chatbot_id, name, color, created_at FROM kb_tags WHERE id = ?`

	tag := &models.KBTag{}
	err := r.db.QueryRow(query, id).Scan(&tag.ID, &tag.ChatbotID, &tag.Name, &tag.Color, &tag.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *KBTagRepository) GetByChatbot(chatbotID int64) ([]models.KBTag, error) {
	query := `SELECT id, chatbot_id, name, color, created_at 
	          FROM kb_tags WHERE chatbot_id = ? ORDER BY name ASC`

	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.KBTag
	for rows.Next() {
		var tag models.KBTag
		err := rows.Scan(&tag.ID, &tag.ChatbotID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *KBTagRepository) Update(tag *models.KBTag) error {
	query := `UPDATE kb_tags SET name = ?, color = ? WHERE id = ?`
	_, err := r.db.Exec(query, tag.Name, tag.Color, tag.ID)
	return err
}

func (r *KBTagRepository) Delete(id int64) error {
	query := `DELETE FROM kb_tags WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Tag assignment operations
func (r *KBTagRepository) AssignTagToKB(kbID, tagID int64) error {
	query := `INSERT INTO kb_entry_tags (kb_id, tag_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, kbID, tagID)
	return err
}

func (r *KBTagRepository) RemoveTagFromKB(kbID, tagID int64) error {
	query := `DELETE FROM kb_entry_tags WHERE kb_id = ? AND tag_id = ?`
	_, err := r.db.Exec(query, kbID, tagID)
	return err
}

func (r *KBTagRepository) GetTagsByKB(kbID int64) ([]models.KBTag, error) {
	query := `SELECT t.id, t.chatbot_id, t.name, t.color, t.created_at 
	          FROM kb_tags t
	          INNER JOIN kb_entry_tags et ON t.id = et.tag_id
	          WHERE et.kb_id = ?
	          ORDER BY t.name ASC`

	rows, err := r.db.Query(query, kbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.KBTag
	for rows.Next() {
		var tag models.KBTag
		err := rows.Scan(&tag.ID, &tag.ChatbotID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *KBTagRepository) GetKBsByTag(tagID int64) ([]int64, error) {
	query := `SELECT kb_id FROM kb_entry_tags WHERE tag_id = ?`

	rows, err := r.db.Query(query, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kbIDs []int64
	for rows.Next() {
		var kbID int64
		if err := rows.Scan(&kbID); err != nil {
			return nil, err
		}
		kbIDs = append(kbIDs, kbID)
	}

	return kbIDs, nil
}

func (r *KBTagRepository) ClearKBTags(kbID int64) error {
	query := `DELETE FROM kb_entry_tags WHERE kb_id = ?`
	_, err := r.db.Exec(query, kbID)
	return err
}

