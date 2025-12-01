package services

import (
	"errors"
	"strings"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
)

type KBTagService struct {
	repo *repository.KBTagRepository
}

func NewKBTagService(repo *repository.KBTagRepository) *KBTagService {
	return &KBTagService{repo: repo}
}

func (s *KBTagService) CreateTag(tag *models.KBTag) error {
	// Validate tag name
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return errors.New("tag name is required")
	}

	// Normalize tag name (lowercase for consistency)
	tag.Name = strings.ToLower(tag.Name)

	// Set default color
	if tag.Color == "" {
		tag.Color = "#10B981"
	}

	return s.repo.Create(tag)
}

func (s *KBTagService) GetTag(id int64) (*models.KBTag, error) {
	return s.repo.GetByID(id)
}

func (s *KBTagService) GetAllTags(chatbotID int64) ([]models.KBTag, error) {
	return s.repo.GetByChatbot(chatbotID)
}

func (s *KBTagService) UpdateTag(tag *models.KBTag) error {
	// Validate tag exists
	existing, err := s.repo.GetByID(tag.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("tag not found")
	}

	// Validate name
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	tag.Name = strings.ToLower(tag.Name)

	return s.repo.Update(tag)
}

func (s *KBTagService) DeleteTag(id int64) error {
	// Note: Tag assignments will be automatically deleted due to CASCADE constraint
	return s.repo.Delete(id)
}

func (s *KBTagService) AssignTags(kbID int64, tagIDs []int64) error {
	// First, clear existing tags
	if err := s.repo.ClearKBTags(kbID); err != nil {
		return err
	}

	// Assign new tags
	for _, tagID := range tagIDs {
		if err := s.repo.AssignTagToKB(kbID, tagID); err != nil {
			return err
		}
	}

	return nil
}

func (s *KBTagService) AddTag(kbID, tagID int64) error {
	return s.repo.AssignTagToKB(kbID, tagID)
}

func (s *KBTagService) RemoveTag(kbID, tagID int64) error {
	return s.repo.RemoveTagFromKB(kbID, tagID)
}

func (s *KBTagService) GetKBTags(kbID int64) ([]models.KBTag, error) {
	return s.repo.GetTagsByKB(kbID)
}

func (s *KBTagService) GetKBsByTag(tagID int64) ([]int64, error) {
	return s.repo.GetKBsByTag(tagID)
}

func (s *KBTagService) FindOrCreateTag(chatbotID int64, tagName string) (*models.KBTag, error) {
	tagName = strings.ToLower(strings.TrimSpace(tagName))
	if tagName == "" {
		return nil, errors.New("tag name is required")
	}

	// Try to find existing tag
	tags, err := s.repo.GetByChatbot(chatbotID)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.Name == tagName {
			return &tag, nil
		}
	}

	// Create new tag
	newTag := &models.KBTag{
		ChatbotID: chatbotID,
		Name:      tagName,
		Color:     "#10B981",
	}

	if err := s.repo.Create(newTag); err != nil {
		return nil, err
	}

	return newTag, nil
}

func (s *KBTagService) BulkAssignTags(kbID int64, tagNames []string, chatbotID int64) error {
	var tagIDs []int64

	for _, tagName := range tagNames {
		tag, err := s.FindOrCreateTag(chatbotID, tagName)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	return s.AssignTags(kbID, tagIDs)
}


