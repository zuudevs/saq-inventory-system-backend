package services

import (
	"errors"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type LocationService struct {
	LocationRepository *repositories.LocationRepository
}

func (s *LocationService) Create(req dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	location := req.ToModel()

	if strings.TrimSpace(location.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.LocationRepository.Create(location); err != nil {
		return nil, err
	}

	return dto.ToLocationResponse(location), nil
}

func (s *LocationService) FindAll() ([]dto.LocationResponse, error) {
	locations, err := s.LocationRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LocationResponse, len(locations))
	for i := range locations {
		responses[i] = *dto.ToLocationResponse(&locations[i])
	}

	return responses, nil
}

func (s *LocationService) FindByID(id uint64) (*dto.LocationResponse, error) {
	location, err := s.LocationRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, nil
	}

	return dto.ToLocationResponse(location), nil
}

func (s *LocationService) Update(id uint64, req dto.UpdateLocationRequest) (*dto.LocationResponse, error) {
	location, err := s.LocationRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, nil
	}

	req.Apply(location)

	if strings.TrimSpace(location.Name) == "" {
		return nil, errors.New("name is required")
	}

	if err := s.LocationRepository.Update(location); err != nil {
		return nil, err
	}

	return dto.ToLocationResponse(location), nil
}

func (s *LocationService) Delete(id uint64) error {
	return s.LocationRepository.Delete(id)
}
