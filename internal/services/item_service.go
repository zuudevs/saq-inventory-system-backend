package services

import (
	"errors"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type ItemService struct {
	ItemRepository *repositories.ItemRepository
}

func (s *ItemService) Create(req dto.CreateItemRequest) (*dto.ItemResponse, error) {
	item := req.ToModel()

	if err := validateItem(item); err != nil {
		return nil, err
	}

	if err := s.ItemRepository.Create(item); err != nil {
		return nil, err
	}

	return dto.ToItemResponse(item), nil
}

func (s *ItemService) FindAll() ([]dto.ItemResponse, error) {
	items, err := s.ItemRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ItemResponse, len(items))
	for i := range items {
		responses[i] = *dto.ToItemResponse(&items[i])
	}

	return responses, nil
}

func (s *ItemService) FindByID(id uint64) (*dto.ItemResponse, error) {
	item, err := s.ItemRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	return dto.ToItemResponse(item), nil
}

func (s *ItemService) Update(id uint64, req dto.UpdateItemRequest) (*dto.ItemResponse, error) {
	item, err := s.ItemRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	req.Apply(item)

	if err := validateItem(item); err != nil {
		return nil, err
	}

	if err := s.ItemRepository.Update(item); err != nil {
		return nil, err
	}

	return dto.ToItemResponse(item), nil
}

func (s *ItemService) Delete(id uint64) error {
	return s.ItemRepository.Delete(id)
}

func validateItem(item *models.Item) error {
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required")
	}

	if item.CategoryID == 0 {
		return errors.New("category_id is required")
	}

	if strings.TrimSpace(item.AssetCode) == "" {
		return errors.New("asset_code is required")
	}

	switch item.ItemCondition {
	case 
		models.ItemConditionGood,
		models.ItemConditionMinorDamage,
		models.ItemConditionMajorDamage,
		models.ItemConditionLost:
	default:
		return errors.New("invalid item_condition")
	}

	switch item.ItemStatus {
	case
		models.ItemStatusActive,
		models.ItemStatusInactive,
		models.ItemStatusMaintenance,
		models.ItemStatusBorrowed:
	default:
		return errors.New("invalid item_status")
	}

	return nil
}
