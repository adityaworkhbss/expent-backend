package service

import (
	"context"
	"expent-backend/internal/category/model"
	"expent-backend/internal/category/repository"
	"fmt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// ListCategories returns all categories for a user.
func (s *Service) ListCategories(ctx context.Context, userID string) ([]model.Category, error) {
	return s.repo.ListCategories(userID)
}

// CreateCategory creates a new category for the user.
func (s *Service) CreateCategory(ctx context.Context, userID, name, typ, color, icon string) (*model.Category, error) {
	// Simple validation: name and type required (already validated in handler)
	// Ensure uniqueness per user could be added later.
	cat := model.Category{UserID: userID, Name: name, Type: typ, Color: color, Icon: icon}
	return s.repo.CreateCategory(cat)
}

// UpdateCategory updates a category belonging to the user.
func (s *Service) UpdateCategory(ctx context.Context, userID, categoryID, name, typ, color, icon string) (*model.Category, error) {
	cat, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, fmt.Errorf("category not found")
	}
	if cat.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	updated := model.Category{UserID: userID, Name: name, Type: typ, Color: color, Icon: icon}
	return s.repo.UpdateCategory(categoryID, updated)
}

// DeleteCategory removes a category belonging to the user.
func (s *Service) DeleteCategory(ctx context.Context, userID, categoryID string) error {
	// Verify ownership (optional)
	cat, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return fmt.Errorf("category not found")
	}
	if cat.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return s.repo.DeleteCategory(categoryID)
}
