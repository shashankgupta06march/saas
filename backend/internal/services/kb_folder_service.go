package services

import (
	"errors"
	"fmt"

	"github.com/chatbot-saas/backend/internal/models"
	"github.com/chatbot-saas/backend/internal/repository"
)

type KBFolderService struct {
	repo *repository.KBFolderRepository
}

func NewKBFolderService(repo *repository.KBFolderRepository) *KBFolderService {
	return &KBFolderService{repo: repo}
}

type FolderTreeNode struct {
	*models.KBFolder
	Children []*FolderTreeNode `json:"children"`
}

func (s *KBFolderService) CreateFolder(folder *models.KBFolder) error {
	// Validate folder name
	if folder.Name == "" {
		return errors.New("folder name is required")
	}

	// Validate parent exists if specified
	if folder.ParentID.Valid {
		parent, err := s.repo.GetByID(folder.ParentID.Int64)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent folder not found")
		}
		// Ensure parent belongs to same chatbot
		if parent.ChatbotID != folder.ChatbotID {
			return errors.New("parent folder belongs to different chatbot")
		}
	}

	// Set defaults
	if folder.Color == "" {
		folder.Color = "#3B82F6"
	}
	if folder.Icon == "" {
		folder.Icon = "folder"
	}

	return s.repo.Create(folder)
}

func (s *KBFolderService) GetFolder(id int64) (*models.KBFolder, error) {
	return s.repo.GetByID(id)
}

func (s *KBFolderService) UpdateFolder(folder *models.KBFolder) error {
	// Validate folder exists
	existing, err := s.repo.GetByID(folder.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("folder not found")
	}

	// Prevent circular references
	if folder.ParentID.Valid {
		if folder.ParentID.Int64 == folder.ID {
			return errors.New("folder cannot be its own parent")
		}

		// Check if new parent would create a circular reference
		if err := s.validateNoCircularReference(folder.ID, folder.ParentID.Int64); err != nil {
			return err
		}
	}

	return s.repo.Update(folder)
}

func (s *KBFolderService) DeleteFolder(id int64) error {
	// Check if folder has children
	children, err := s.repo.GetChildren(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete folder with children")
	}

	return s.repo.Delete(id)
}

func (s *KBFolderService) GetFolderTree(chatbotID int64) ([]*FolderTreeNode, error) {
	// Get all folders for the chatbot
	folders, err := s.repo.GetByChatbot(chatbotID)
	if err != nil {
		return nil, err
	}

	// Build folder map
	folderMap := make(map[int64]*FolderTreeNode)
	for i := range folders {
		folderMap[folders[i].ID] = &FolderTreeNode{
			KBFolder: &folders[i],
			Children: []*FolderTreeNode{},
		}
	}

	// Build tree structure
	var rootNodes []*FolderTreeNode
	for _, node := range folderMap {
		if node.ParentID.Valid {
			// Add to parent's children
			if parent, exists := folderMap[node.ParentID.Int64]; exists {
				parent.Children = append(parent.Children, node)
			}
		} else {
			// Root level folder
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes, nil
}

func (s *KBFolderService) GetFolderPath(folderID int64) ([]models.KBFolder, error) {
	var path []models.KBFolder
	currentID := folderID

	for currentID != 0 {
		folder, err := s.repo.GetByID(currentID)
		if err != nil {
			return nil, err
		}
		if folder == nil {
			return nil, errors.New("folder not found in path")
		}

		// Prepend to maintain order from root to current
		path = append([]models.KBFolder{*folder}, path...)

		if folder.ParentID.Valid {
			currentID = folder.ParentID.Int64
		} else {
			currentID = 0
		}
	}

	return path, nil
}

func (s *KBFolderService) MoveFolder(folderID, newParentID int64) error {
	folder, err := s.repo.GetByID(folderID)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("folder not found")
	}

	// Prevent moving to self
	if folderID == newParentID {
		return errors.New("cannot move folder into itself")
	}

	// Validate new parent
	if newParentID != 0 {
		parent, err := s.repo.GetByID(newParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent folder not found")
		}
		if parent.ChatbotID != folder.ChatbotID {
			return errors.New("cannot move to folder in different chatbot")
		}

		// Check for circular reference
		if err := s.validateNoCircularReference(folderID, newParentID); err != nil {
			return err
		}

		folder.ParentID.Int64 = newParentID
		folder.ParentID.Valid = true
	} else {
		folder.ParentID.Valid = false
	}

	return s.repo.Update(folder)
}

func (s *KBFolderService) validateNoCircularReference(folderID, newParentID int64) error {
	currentID := newParentID
	depth := 0
	maxDepth := 100 // Prevent infinite loops

	for currentID != 0 && depth < maxDepth {
		if currentID == folderID {
			return errors.New("circular reference detected")
		}

		folder, err := s.repo.GetByID(currentID)
		if err != nil {
			return err
		}
		if folder == nil {
			break
		}

		if folder.ParentID.Valid {
			currentID = folder.ParentID.Int64
		} else {
			break
		}
		depth++
	}

	if depth >= maxDepth {
		return fmt.Errorf("folder hierarchy too deep (max %d levels)", maxDepth)
	}

	return nil
}

func (s *KBFolderService) GetRootFolders(chatbotID int64) ([]models.KBFolder, error) {
	return s.repo.GetRootFolders(chatbotID)
}

func (s *KBFolderService) GetSubfolders(folderID int64) ([]models.KBFolder, error) {
	return s.repo.GetChildren(folderID)
}
