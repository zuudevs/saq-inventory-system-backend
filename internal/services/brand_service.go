package services

import (
	"errors"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type BrandService struct {
	BrandRepository *repositories.BrandRepository
}

func (s *BrandService) Create(req dto.CreateBrandRequest) (*dto.BrandResponse, error) {
	brand := req.ToModel()

	if strings.TrimSpace(brand.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.BrandRepository.Create(brand); err != nil {
		return nil, err
	}

	return dto.ToBrandResponse(brand), nil
}

func (s *BrandService) FindAll() ([]dto.BrandResponse, error) {
	brands, err := s.BrandRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BrandResponse, len(brands))
	for i := range brands {
		responses[i] = *dto.ToBrandResponse(&brands[i])
	}

	return responses, nil
}

func (s *BrandService) FindByID(id uint64) (*dto.BrandResponse, error) {
	brand, err := s.BrandRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if brand == nil {
		return nil, nil
	}

	return dto.ToBrandResponse(brand), nil
}

func (s *BrandService) Update(id uint64, req dto.UpdateBrandRequest) (*dto.BrandResponse, error) {
	brand, err := s.BrandRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if brand == nil {
		return nil, nil
	}

	req.Apply(brand)

	if strings.TrimSpace(brand.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.BrandRepository.Update(brand); err != nil {
		return nil, err
	}

	return dto.ToBrandResponse(brand), nil
}

func (s *BrandService) Delete(id uint64) error {
	return s.BrandRepository.Delete(id)
}