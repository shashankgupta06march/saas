package repository

import (
	"database/sql"
	"time"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBSyncRepository struct {
	db *sql.DB
}

func NewKBSyncRepository(db *sql.DB) *KBSyncRepository {
	return &KBSyncRepository{db: db}
}

func (r *KBSyncRepository) Create(sync *models.KBSyncSource) error {
	query := `INSERT INTO kb_sync_sources 
	          (chatbot_id, source_type, source_identifier, auth_token, sync_frequency, is_active, sync_settings) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, sync.ChatbotID, sync.SourceType, sync.SourceIdentifier,
		sync.AuthToken, sync.SyncFrequency, sync.IsActive, sync.SyncSettings)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	sync.ID = id
	return nil
}

func (r *KBSyncRepository) GetByID(id int64) (*models.KBSyncSource, error) {
	query := `SELECT id, chatbot_id, source_type, source_identifier, auth_token, sync_frequency, 
	          last_sync_at, next_sync_at, is_active, sync_settings, created_at, updated_at 
	          FROM kb_sync_sources WHERE id = ?`

	sync := &models.KBSyncSource{}
	err := r.db.QueryRow(query, id).Scan(
		&sync.ID, &sync.ChatbotID, &sync.SourceType, &sync.SourceIdentifier,
		&sync.AuthToken, &sync.SyncFrequency, &sync.LastSyncAt, &sync.NextSyncAt,
		&sync.IsActive, &sync.SyncSettings, &sync.CreatedAt, &sync.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return sync, nil
}

func (r *KBSyncRepository) GetByChatbot(chatbotID int64) ([]models.KBSyncSource, error) {
	query := `SELECT id, chatbot_id, source_type, source_identifier, auth_token, sync_frequency, 
	          last_sync_at, next_sync_at, is_active, sync_settings, created_at, updated_at 
	          FROM kb_sync_sources WHERE chatbot_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var syncs []models.KBSyncSource
	for rows.Next() {
		var sync models.KBSyncSource
		err := rows.Scan(
			&sync.ID, &sync.ChatbotID, &sync.SourceType, &sync.SourceIdentifier,
			&sync.AuthToken, &sync.SyncFrequency, &sync.LastSyncAt, &sync.NextSyncAt,
			&sync.IsActive, &sync.SyncSettings, &sync.CreatedAt, &sync.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		syncs = append(syncs, sync)
	}

	return syncs, nil
}

func (r *KBSyncRepository) Update(sync *models.KBSyncSource) error {
	query := `UPDATE kb_sync_sources 
	          SET source_type = ?, source_identifier = ?, auth_token = ?, 
	              sync_frequency = ?, is_active = ?, sync_settings = ? 
	          WHERE id = ?`

	_, err := r.db.Exec(query, sync.SourceType, sync.SourceIdentifier, sync.AuthToken,
		sync.SyncFrequency, sync.IsActive, sync.SyncSettings, sync.ID)
	return err
}

func (r *KBSyncRepository) UpdateSyncTime(id int64, lastSync, nextSync time.Time) error {
	query := `UPDATE kb_sync_sources SET last_sync_at = ?, next_sync_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, lastSync, nextSync, id)
	return err
}

func (r *KBSyncRepository) Delete(id int64) error {
	query := `DELETE FROM kb_sync_sources WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *KBSyncRepository) GetDueSyncs() ([]models.KBSyncSource, error) {
	query := `SELECT id, chatbot_id, source_type, source_identifier, auth_token, sync_frequency, 
	          last_sync_at, next_sync_at, is_active, sync_settings, created_at, updated_at 
	          FROM kb_sync_sources 
	          WHERE is_active = true AND (next_sync_at IS NULL OR next_sync_at <= NOW())
	          ORDER BY next_sync_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var syncs []models.KBSyncSource
	for rows.Next() {
		var sync models.KBSyncSource
		err := rows.Scan(
			&sync.ID, &sync.ChatbotID, &sync.SourceType, &sync.SourceIdentifier,
			&sync.AuthToken, &sync.SyncFrequency, &sync.LastSyncAt, &sync.NextSyncAt,
			&sync.IsActive, &sync.SyncSettings, &sync.CreatedAt, &sync.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		syncs = append(syncs, sync)
	}

	return syncs, nil
}

func (r *KBSyncRepository) ToggleActive(id int64, active bool) error {
	query := `UPDATE kb_sync_sources SET is_active = ? WHERE id = ?`
	_, err := r.db.Exec(query, active, id)
	return err
}
