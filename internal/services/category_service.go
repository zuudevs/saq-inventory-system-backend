package services

import (
	"errors"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type CategoryService struct {
	CategoryRepository *repositories.CategoryRepository
}

func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := req.ToModel()

	if strings.TrimSpace(category.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.CategoryRepository.Create(category); err != nil {
		return nil, err
	}

	return dto.ToCategoryResponse(category), nil
}

func (s *CategoryService) FindAll() ([]dto.CategoryResponse, error) {
	categories, err := s.CategoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = *dto.ToCategoryResponse(&categories[i])
	}

	return responses, nil
}

func (s *CategoryService) FindByID(id uint64) (*dto.CategoryResponse, error) {
	category, err := s.CategoryRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	return dto.ToCategoryResponse(category), nil
}

func (s *CategoryService) Update(id uint64, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.CategoryRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	req.Apply(category)

	if strings.TrimSpace(category.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.CategoryRepository.Update(category); err != nil {
		return nil, err
	}

	return dto.ToCategoryResponse(category), nil
}

func (s *CategoryService) Delete(id uint64) error {
	return s.CategoryRepository.Delete(id)
}
