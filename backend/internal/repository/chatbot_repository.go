package repository

import (
	"database/sql"
	"errors"

	"github.com/chatbot-saas/backend/internal/models"
)

type ChatbotRepository struct {
	db *sql.DB
}

func NewChatbotRepository(db *sql.DB) *ChatbotRepository {
	return &ChatbotRepository{db: db}
}

func (r *ChatbotRepository) Create(chatbot *models.Chatbot) error {
	query := `INSERT INTO chatbots (organization_id, name, description, status) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, chatbot.OrganizationID, chatbot.Name, chatbot.Description, chatbot.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	chatbot.ID = id

	// Create default settings for this chatbot
	settingsQuery := `INSERT INTO chatbot_settings (chatbot_id) VALUES (?)`
	_, err = r.db.Exec(settingsQuery, id)
	return err
}

func (r *ChatbotRepository) GetByID(id int64) (*models.Chatbot, error) {
	chatbot := &models.Chatbot{}
	query := `SELECT id, organization_id, name, description, status, created_at, updated_at FROM chatbots WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(
		&chatbot.ID,
		&chatbot.OrganizationID,
		&chatbot.Name,
		&chatbot.Description,
		&chatbot.Status,
		&chatbot.CreatedAt,
		&chatbot.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("chatbot not found")
	}

	return chatbot, err
}

func (r *ChatbotRepository) GetByOrganization(orgID int64) ([]models.Chatbot, error) {
	query := `SELECT id, organization_id, name, description, status, created_at, updated_at FROM chatbots WHERE organization_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatbots []models.Chatbot
	for rows.Next() {
		var chatbot models.Chatbot
		err := rows.Scan(
			&chatbot.ID,
			&chatbot.OrganizationID,
			&chatbot.Name,
			&chatbot.Description,
			&chatbot.Status,
			&chatbot.CreatedAt,
			&chatbot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		chatbots = append(chatbots, chatbot)
	}

	return chatbots, nil
}

func (r *ChatbotRepository) Update(chatbot *models.Chatbot) error {
	query := `UPDATE chatbots SET name = ?, description = ?, status = ? WHERE id = ?`
	_, err := r.db.Exec(query, chatbot.Name, chatbot.Description, chatbot.Status, chatbot.ID)
	return err
}

func (r *ChatbotRepository) Delete(id int64) error {
	query := `DELETE FROM chatbots WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ChatbotRepository) GetSettings(chatbotID int64) (*models.ChatbotSettings, error) {
	settings := &models.ChatbotSettings{}
	query := `SELECT id, chatbot_id, theme_color, position, welcome_message, avatar_url, custom_css, widget_size FROM chatbot_settings WHERE chatbot_id = ?`
	err := r.db.QueryRow(query, chatbotID).Scan(
		&settings.ID,
		&settings.ChatbotID,
		&settings.ThemeColor,
		&settings.Position,
		&settings.WelcomeMessage,
		&settings.AvatarURL,
		&settings.CustomCSS,
		&settings.WidgetSize,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("settings not found")
	}

	return settings, err
}

func (r *ChatbotRepository) UpdateSettings(settings *models.ChatbotSettings) error {
	query := `UPDATE chatbot_settings SET theme_color = ?, position = ?, welcome_message = ?, avatar_url = ?, custom_css = ?, widget_size = ? WHERE chatbot_id = ?`
	_, err := r.db.Exec(query, settings.ThemeColor, settings.Position, settings.WelcomeMessage, settings.AvatarURL, settings.CustomCSS, settings.WidgetSize, settings.ChatbotID)
	return err
}


