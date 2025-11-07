package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
)

type KBSyncService struct {
	syncRepo *repository.KBSyncRepository
	kbRepo   *repository.KnowledgeRepository
}

func NewKBSyncService(syncRepo *repository.KBSyncRepository, kbRepo *repository.KnowledgeRepository) *KBSyncService {
	return &KBSyncService{
		syncRepo: syncRepo,
		kbRepo:   kbRepo,
	}
}

// CreateSyncSource creates a new sync source configuration
func (s *KBSyncService) CreateSyncSource(sync *models.KBSyncSource) error {
	// Validate source type
	validTypes := map[string]bool{
		"google_drive": true,
		"dropbox":      true,
		"notion":       true,
		"wordpress":    true,
		"confluence":   true,
		"webhook":      true,
	}

	if !validTypes[sync.SourceType] {
		return fmt.Errorf("invalid source type: %s", sync.SourceType)
	}

	// Validate sync frequency
	validFrequencies := map[string]bool{
		"realtime": true,
		"hourly":   true,
		"daily":    true,
		"weekly":   true,
	}

	if !validFrequencies[sync.SyncFrequency] {
		return errors.New("invalid sync frequency")
	}

	// Calculate next sync time
	sync.NextSyncAt.Time = s.calculateNextSync(time.Now(), sync.SyncFrequency)
	sync.NextSyncAt.Valid = true

	return s.syncRepo.Create(sync)
}

// GetSyncSource retrieves a sync source by ID
func (s *KBSyncService) GetSyncSource(id int64) (*models.KBSyncSource, error) {
	return s.syncRepo.GetByID(id)
}

// GetSyncSources retrieves all sync sources for a chatbot
func (s *KBSyncService) GetSyncSources(chatbotID int64) ([]models.KBSyncSource, error) {
	return s.syncRepo.GetByChatbot(chatbotID)
}

// UpdateSyncSource updates a sync source configuration
func (s *KBSyncService) UpdateSyncSource(sync *models.KBSyncSource) error {
	// Validate sync source exists
	existing, err := s.syncRepo.GetByID(sync.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("sync source not found")
	}

	return s.syncRepo.Update(sync)
}

// DeleteSyncSource deletes a sync source
func (s *KBSyncService) DeleteSyncSource(id int64) error {
	return s.syncRepo.Delete(id)
}

// ToggleSyncSource enables or disables a sync source
func (s *KBSyncService) ToggleSyncSource(id int64, active bool) error {
	return s.syncRepo.ToggleActive(id, active)
}

// ExecuteSync performs synchronization for a sync source
func (s *KBSyncService) ExecuteSync(syncID int64) error {
	sync, err := s.syncRepo.GetByID(syncID)
	if err != nil {
		return err
	}
	if sync == nil {
		return errors.New("sync source not found")
	}

	if !sync.IsActive {
		return errors.New("sync source is not active")
	}

	// Execute sync based on source type
	switch sync.SourceType {
	case "google_drive":
		return s.syncGoogleDrive(sync)
	case "dropbox":
		return s.syncDropbox(sync)
	case "notion":
		return s.syncNotion(sync)
	case "wordpress":
		return s.syncWordPress(sync)
	case "confluence":
		return s.syncConfluence(sync)
	case "webhook":
		return s.syncWebhook(sync)
	default:
		return fmt.Errorf("unsupported source type: %s", sync.SourceType)
	}
}

// Placeholder implementations for different sync sources
func (s *KBSyncService) syncGoogleDrive(sync *models.KBSyncSource) error {
	// TODO: Implement Google Drive OAuth and file sync
	// 1. Authenticate with Google Drive API
	// 2. List files in specified folder
	// 3. Download new/modified files
	// 4. Create/update knowledge base entries
	// 5. Update last_sync_at

	// For now, return a placeholder message
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

func (s *KBSyncService) syncDropbox(sync *models.KBSyncSource) error {
	// TODO: Implement Dropbox API integration
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

func (s *KBSyncService) syncNotion(sync *models.KBSyncSource) error {
	// TODO: Implement Notion API integration
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

func (s *KBSyncService) syncWordPress(sync *models.KBSyncSource) error {
	// TODO: Implement WordPress REST API integration
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

func (s *KBSyncService) syncConfluence(sync *models.KBSyncSource) error {
	// TODO: Implement Confluence API integration
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

func (s *KBSyncService) syncWebhook(sync *models.KBSyncSource) error {
	// TODO: Implement webhook-based sync
	now := time.Now()
	nextSync := s.calculateNextSync(now, sync.SyncFrequency)
	return s.syncRepo.UpdateSyncTime(sync.ID, now, nextSync)
}

// ProcessDueSyncs finds and executes all syncs that are due
func (s *KBSyncService) ProcessDueSyncs() error {
	syncs, err := s.syncRepo.GetDueSyncs()
	if err != nil {
		return err
	}

	for _, sync := range syncs {
		// Execute sync in a goroutine to prevent blocking
		go func(syncID int64) {
			if err := s.ExecuteSync(syncID); err != nil {
				// Log error (in production, use proper logging)
				fmt.Printf("Sync failed for source %d: %v\n", syncID, err)
			}
		}(sync.ID)
	}

	return nil
}

// calculateNextSync calculates the next sync time based on frequency
func (s *KBSyncService) calculateNextSync(from time.Time, frequency string) time.Time {
	switch frequency {
	case "realtime":
		return from.Add(5 * time.Minute) // Check every 5 minutes for realtime
	case "hourly":
		return from.Add(1 * time.Hour)
	case "daily":
		return from.Add(24 * time.Hour)
	case "weekly":
		return from.Add(7 * 24 * time.Hour)
	default:
		return from.Add(24 * time.Hour) // Default to daily
	}
}

// TriggerManualSync manually triggers a sync regardless of schedule
func (s *KBSyncService) TriggerManualSync(syncID int64) error {
	return s.ExecuteSync(syncID)
}

// GetSyncStatus returns the status of a sync source
func (s *KBSyncService) GetSyncStatus(syncID int64) (map[string]interface{}, error) {
	sync, err := s.syncRepo.GetByID(syncID)
	if err != nil {
		return nil, err
	}
	if sync == nil {
		return nil, errors.New("sync source not found")
	}

	status := map[string]interface{}{
		"id":             sync.ID,
		"source_type":    sync.SourceType,
		"is_active":      sync.IsActive,
		"sync_frequency": sync.SyncFrequency,
		"last_synced_at": nil,
		"next_sync_at":   nil,
		"sync_overdue":   false,
	}

	if sync.LastSyncAt.Valid {
		status["last_synced_at"] = sync.LastSyncAt.Time
	}

	if sync.NextSyncAt.Valid {
		status["next_sync_at"] = sync.NextSyncAt.Time
		status["sync_overdue"] = time.Now().After(sync.NextSyncAt.Time)
	}

	return status, nil
}

