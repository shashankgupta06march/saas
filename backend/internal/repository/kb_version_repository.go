package repository

import (
	"database/sql"

	"github.com/chatbot-saas/backend/internal/models"
)

type KBVersionRepository struct {
	db *sql.DB
}

func NewKBVersionRepository(db *sql.DB) *KBVersionRepository {
	return &KBVersionRepository{db: db}
}

func (r *KBVersionRepository) Create(version *models.KBVersion) error {
	query := `INSERT INTO kb_versions (kb_id, version, title, content, content_type, source_url, changed_by, change_summary) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, version.KBID, version.Version, version.Title, version.Content,
		version.ContentType, version.SourceURL, version.ChangedBy, version.ChangeSummary)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	version.ID = id
	return nil
}

func (r *KBVersionRepository) GetByKBID(kbID int64) ([]models.KBVersion, error) {
	query := `SELECT id, kb_id, version, title, content, content_type, source_url, changed_by, change_summary, created_at 
	          FROM kb_versions WHERE kb_id = ? ORDER BY version DESC`

	rows, err := r.db.Query(query, kbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.KBVersion
	for rows.Next() {
		var v models.KBVersion
		err := rows.Scan(
			&v.ID, &v.KBID, &v.Version, &v.Title, &v.Content,
			&v.ContentType, &v.SourceURL, &v.ChangedBy, &v.ChangeSummary, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (r *KBVersionRepository) GetByVersion(kbID int64, version int) (*models.KBVersion, error) {
	query := `SELECT id, kb_id, version, title, content, content_type, source_url, changed_by, change_summary, created_at 
	          FROM kb_versions WHERE kb_id = ? AND version = ?`

	v := &models.KBVersion{}
	err := r.db.QueryRow(query, kbID, version).Scan(
		&v.ID, &v.KBID, &v.Version, &v.Title, &v.Content,
		&v.ContentType, &v.SourceURL, &v.ChangedBy, &v.ChangeSummary, &v.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *KBVersionRepository) GetLatestVersion(kbID int64) (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) FROM kb_versions WHERE kb_id = ?`

	var version int
	err := r.db.QueryRow(query, kbID).Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (r *KBVersionRepository) DeleteOldVersions(kbID int64, keepLast int) error {
	query := `DELETE FROM kb_versions 
	          WHERE kb_id = ? 
	          AND version NOT IN (
	              SELECT version FROM (
	                  SELECT version FROM kb_versions 
	                  WHERE kb_id = ? 
	                  ORDER BY version DESC 
	                  LIMIT ?
	              ) AS keep_versions
	          )`

	_, err := r.db.Exec(query, kbID, kbID, keepLast)
	return err
}


