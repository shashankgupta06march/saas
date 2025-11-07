package repository

import (
	"database/sql"
	"errors"

	"github.com/chatbot-saas/backend/internal/models"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(org *models.Organization) error {
	query := `INSERT INTO organizations (name, api_key, plan_type, status) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, org.Name, org.APIKey, org.PlanType, org.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	org.ID = id
	return nil
}

func (r *OrganizationRepository) GetByID(id int64) (*models.Organization, error) {
	org := &models.Organization{}
	query := `SELECT id, name, api_key, plan_type, status, created_at, updated_at FROM organizations WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(
		&org.ID,
		&org.Name,
		&org.APIKey,
		&org.PlanType,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("organization not found")
	}

	return org, err
}

func (r *OrganizationRepository) GetByAPIKey(apiKey string) (*models.Organization, error) {
	org := &models.Organization{}
	query := `SELECT id, name, api_key, plan_type, status, created_at, updated_at FROM organizations WHERE api_key = ?`
	err := r.db.QueryRow(query, apiKey).Scan(
		&org.ID,
		&org.Name,
		&org.APIKey,
		&org.PlanType,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("organization not found")
	}

	return org, err
}

func (r *OrganizationRepository) Update(org *models.Organization) error {
	query := `UPDATE organizations SET name = ?, plan_type = ?, status = ? WHERE id = ?`
	_, err := r.db.Exec(query, org.Name, org.PlanType, org.Status, org.ID)
	return err
}

