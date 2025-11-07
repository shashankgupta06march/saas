package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBFolderRepository struct {
	db *sql.DB
}

func NewKBFolderRepository(db *sql.DB) *KBFolderRepository {
	return &KBFolderRepository{db: db}
}

func (r *KBFolderRepository) Create(folder *models.KBFolder) error {
	query := `INSERT INTO kb_folders (chatbot_id, name, description, parent_id, color, icon) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, folder.ChatbotID, folder.Name, folder.Description,
		folder.ParentID, folder.Color, folder.Icon)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	folder.ID = id
	return nil
}

func (r *KBFolderRepository) GetByID(id int64) (*models.KBFolder, error) {
	query := `SELECT id, chatbot_id, name, description, parent_id, color, icon, created_at, updated_at 
	          FROM kb_folders WHERE id = ?`

	folder := &models.KBFolder{}
	err := r.db.QueryRow(query, id).Scan(
		&folder.ID, &folder.ChatbotID, &folder.Name, &folder.Description,
		&folder.ParentID, &folder.Color, &folder.Icon, &folder.CreatedAt, &folder.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *KBFolderRepository) GetByChatbot(chatbotID int64) ([]models.KBFolder, error) {
	query := `SELECT id, chatbot_id, name, description, parent_id, color, icon, created_at, updated_at 
	          FROM kb_folders WHERE chatbot_id = ? ORDER BY name ASC`

	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.KBFolder
	for rows.Next() {
		var folder models.KBFolder
		err := rows.Scan(
			&folder.ID, &folder.ChatbotID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.Color, &folder.Icon, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *KBFolderRepository) Update(folder *models.KBFolder) error {
	query := `UPDATE kb_folders 
	          SET name = ?, description = ?, parent_id = ?, color = ?, icon = ? 
	          WHERE id = ?`

	_, err := r.db.Exec(query, folder.Name, folder.Description, folder.ParentID,
		folder.Color, folder.Icon, folder.ID)
	return err
}

func (r *KBFolderRepository) Delete(id int64) error {
	query := `DELETE FROM kb_folders WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *KBFolderRepository) GetChildren(parentID int64) ([]models.KBFolder, error) {
	query := `SELECT id, chatbot_id, name, description, parent_id, color, icon, created_at, updated_at 
	          FROM kb_folders WHERE parent_id = ? ORDER BY name ASC`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.KBFolder
	for rows.Next() {
		var folder models.KBFolder
		err := rows.Scan(
			&folder.ID, &folder.ChatbotID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.Color, &folder.Icon, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *KBFolderRepository) GetRootFolders(chatbotID int64) ([]models.KBFolder, error) {
	query := `SELECT id, chatbot_id, name, description, parent_id, color, icon, created_at, updated_at 
	          FROM kb_folders WHERE chatbot_id = ? AND parent_id IS NULL ORDER BY name ASC`

	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.KBFolder
	for rows.Next() {
		var folder models.KBFolder
		err := rows.Scan(
			&folder.ID, &folder.ChatbotID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.Color, &folder.Icon, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

