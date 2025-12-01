package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
)

type KBVersionService struct {
	versionRepo *repository.KBVersionRepository
	kbRepo      *repository.KnowledgeRepository
}

func NewKBVersionService(versionRepo *repository.KBVersionRepository, kbRepo *repository.KnowledgeRepository) *KBVersionService {
	return &KBVersionService{
		versionRepo: versionRepo,
		kbRepo:      kbRepo,
	}
}

// SaveVersion creates a new version of a knowledge base entry
func (s *KBVersionService) SaveVersion(kb *models.KnowledgeBase, userID int64, changeSummary string) error {
	// Get the next version number
	latestVersion, err := s.versionRepo.GetLatestVersion(kb.ID)
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	nextVersion := latestVersion + 1

	// Create version record
	version := &models.KBVersion{
		KBID:          kb.ID,
		Version:       nextVersion,
		Title:         sql.NullString{String: kb.Title, Valid: true},
		Content:       kb.Content,
		ContentType:   sql.NullString{String: kb.ContentType, Valid: true},
		SourceURL:     sql.NullString{String: kb.SourceURL, Valid: true},
		ChangedBy:     sql.NullInt64{Int64: userID, Valid: true},
		ChangeSummary: sql.NullString{String: changeSummary, Valid: true},
	}

	return s.versionRepo.Create(version)
}

// GetVersionHistory retrieves all versions for a knowledge base entry
func (s *KBVersionService) GetVersionHistory(kbID int64) ([]models.KBVersion, error) {
	return s.versionRepo.GetByKBID(kbID)
}

// GetVersion retrieves a specific version
func (s *KBVersionService) GetVersion(kbID int64, version int) (*models.KBVersion, error) {
	return s.versionRepo.GetByVersion(kbID, version)
}

// RestoreVersion restores a knowledge base entry to a previous version
func (s *KBVersionService) RestoreVersion(kbID int64, version int, userID int64) error {
	// Get the version to restore
	oldVersion, err := s.versionRepo.GetByVersion(kbID, version)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}
	if oldVersion == nil {
		return errors.New("version not found")
	}

	// Get current KB entry
	kb, err := s.kbRepo.GetByID(kbID)
	if err != nil {
		return fmt.Errorf("failed to get knowledge base entry: %w", err)
	}
	if kb == nil {
		return errors.New("knowledge base entry not found")
	}

	// Save current state as a version before restoring
	if err := s.SaveVersion(kb, userID, fmt.Sprintf("Before restoring to version %d", version)); err != nil {
		return fmt.Errorf("failed to save current state: %w", err)
	}

	// Update KB with old version content
	kb.Title = oldVersion.Title.String
	kb.Content = oldVersion.Content
	if oldVersion.ContentType.Valid {
		kb.ContentType = oldVersion.ContentType.String
	}
	if oldVersion.SourceURL.Valid {
		kb.SourceURL = oldVersion.SourceURL.String
	}

	// Increment version number
	kb.Version++

	// Update in database
	if err := s.kbRepo.Update(kb); err != nil {
		return fmt.Errorf("failed to update knowledge base: %w", err)
	}

	// Create new version record for the restore
	if err := s.SaveVersion(kb, userID, fmt.Sprintf("Restored from version %d", version)); err != nil {
		return fmt.Errorf("failed to save restored version: %w", err)
	}

	return nil
}

// CompareVersions compares two versions and returns differences
func (s *KBVersionService) CompareVersions(kbID int64, version1, version2 int) (map[string]interface{}, error) {
	v1, err := s.versionRepo.GetByVersion(kbID, version1)
	if err != nil || v1 == nil {
		return nil, errors.New("version 1 not found")
	}

	v2, err := s.versionRepo.GetByVersion(kbID, version2)
	if err != nil || v2 == nil {
		return nil, errors.New("version 2 not found")
	}

	comparison := map[string]interface{}{
		"version1": version1,
		"version2": version2,
		"changes":  []string{},
	}

	var changes []string

	if v1.Title.String != v2.Title.String {
		changes = append(changes, "title")
	}
	if v1.Content != v2.Content {
		changes = append(changes, "content")
	}
	if v1.ContentType.String != v2.ContentType.String {
		changes = append(changes, "content_type")
	}
	if v1.SourceURL.String != v2.SourceURL.String {
		changes = append(changes, "source_url")
	}

	comparison["changes"] = changes
	comparison["has_changes"] = len(changes) > 0

	return comparison, nil
}

// CleanupOldVersions keeps only the latest N versions
func (s *KBVersionService) CleanupOldVersions(kbID int64, keepLast int) error {
	if keepLast <= 0 {
		keepLast = 10 // Default to keeping last 10 versions
	}

	return s.versionRepo.DeleteOldVersions(kbID, keepLast)
}

// GetLatestVersion gets the latest version number
func (s *KBVersionService) GetLatestVersion(kbID int64) (int, error) {
	return s.versionRepo.GetLatestVersion(kbID)
}


