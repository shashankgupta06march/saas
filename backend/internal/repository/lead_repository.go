package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/chatbot-saas/backend/internal/models"
)

type LeadRepository struct {
	db *sql.DB
}

func NewLeadRepository(db *sql.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) GetConfig(chatbotID int64) (*models.LeadCaptureConfig, error) {
	cfg := &models.LeadCaptureConfig{}
	query := `SELECT id, chatbot_id, enabled, title, subtitle, fields, created_at, updated_at
	          FROM lead_capture_config WHERE chatbot_id = ?`
	err := r.db.QueryRow(query, chatbotID).Scan(
		&cfg.ID, &cfg.ChatbotID, &cfg.Enabled,
		&cfg.Title, &cfg.Subtitle, &cfg.Fields,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		// Return disabled default so widget always gets a valid response.
		defaultFields, _ := json.Marshal([]models.LeadCaptureField{
			{Name: "name", Label: "Your Name", Type: "text", Required: true, Placeholder: "Enter your name"},
			{Name: "email", Label: "Email Address", Type: "email", Required: true, Placeholder: "Enter your email"},
		})
		return &models.LeadCaptureConfig{
			ChatbotID: chatbotID,
			Enabled:   false,
			Title:     "Before we begin...",
			Subtitle:  "Please share a few details so we can assist you better.",
			Fields:    string(defaultFields),
		}, nil
	}
	return cfg, err
}

func (r *LeadRepository) UpsertConfig(cfg *models.LeadCaptureConfig) error {
	query := `INSERT INTO lead_capture_config (chatbot_id, enabled, title, subtitle, fields)
	          VALUES (?, ?, ?, ?, ?)
	          ON DUPLICATE KEY UPDATE
	            enabled = VALUES(enabled),
	            title   = VALUES(title),
	            subtitle = VALUES(subtitle),
	            fields  = VALUES(fields),
	            updated_at = CURRENT_TIMESTAMP`
	result, err := r.db.Exec(query, cfg.ChatbotID, cfg.Enabled, cfg.Title, cfg.Subtitle, cfg.Fields)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	if id > 0 {
		cfg.ID = id
	}
	return nil
}

func (r *LeadRepository) CreateLead(lead *models.Lead) error {
	query := `INSERT INTO leads (chatbot_id, session_id, field_values) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, lead.ChatbotID, lead.SessionID, lead.FieldValues)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	lead.ID = id
	return nil
}

func (r *LeadRepository) GetLeads(chatbotID int64) ([]models.Lead, error) {
	query := `SELECT id, chatbot_id, session_id, field_values, created_at
	          FROM leads WHERE chatbot_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []models.Lead
	for rows.Next() {
		var l models.Lead
		if err := rows.Scan(&l.ID, &l.ChatbotID, &l.SessionID, &l.FieldValues, &l.CreatedAt); err != nil {
			return nil, err
		}
		leads = append(leads, l)
	}
	return leads, nil
}
