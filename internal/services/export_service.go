package services

import (
	"archive/zip"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/exporters"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type ExportService struct {
	DB                 *sqlx.DB
	BrandRepository    *repositories.BrandRepository
	CategoryRepository *repositories.CategoryRepository
	ItemRepository     *repositories.ItemRepository
	LocationRepository *repositories.LocationRepository
	ImageRepository    *repositories.ImageRepository
}

// ExportCSV fetches all exportable resources and exports them in CSV format bundled into a ZIP archive.
func (s *ExportService) ExportCSV(writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)

	exportResource := func(filename string, getDTOs func() (interface{}, error)) error {
		fw, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}
		data, err := getDTOs()
		if err != nil {
			return err
		}
		return exporters.ExportCSV(fw, data)
	}

	if err := exportResource("brands.csv", func() (interface{}, error) { return s.getBrandResponses() }); err != nil {
		return err
	}
	if err := exportResource("categories.csv", func() (interface{}, error) { return s.getCategoryResponses() }); err != nil {
		return err
	}
	if err := exportResource("locations.csv", func() (interface{}, error) { return s.getLocationResponses() }); err != nil {
		return err
	}
	if err := exportResource("items.csv", func() (interface{}, error) { return s.getItemResponses() }); err != nil {
		return err
	}
	if err := exportResource("images.csv", func() (interface{}, error) { return s.getImageResponses() }); err != nil {
		return err
	}

	return zipWriter.Close()
}

// ExportXLSX fetches all exportable resources and exports them in XLSX format to the given writer.
func (s *ExportService) ExportXLSX(writer io.Writer) error {
	brands, err := s.getBrandResponses()
	if err != nil {
		return err
	}
	categories, err := s.getCategoryResponses()
	if err != nil {
		return err
	}
	locations, err := s.getLocationResponses()
	if err != nil {
		return err
	}
	items, err := s.getItemResponses()
	if err != nil {
		return err
	}
	images, err := s.getImageResponses()
	if err != nil {
		return err
	}

	sheets := []exporters.SheetData{
		{Name: "Brands", Data: brands},
		{Name: "Categories", Data: categories},
		{Name: "Locations", Data: locations},
		{Name: "Items", Data: items},
		{Name: "Images", Data: images},
	}

	return exporters.ExportMultiSheetXLSX(writer, sheets)
}

func (s *ExportService) getBrandResponses() ([]dto.BrandResponse, error) {
	if s.BrandRepository == nil {
		return []dto.BrandResponse{}, nil
	}
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

func (s *ExportService) getCategoryResponses() ([]dto.CategoryResponse, error) {
	if s.CategoryRepository == nil {
		return []dto.CategoryResponse{}, nil
	}
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

func (s *ExportService) getLocationResponses() ([]dto.LocationResponse, error) {
	if s.LocationRepository == nil {
		return []dto.LocationResponse{}, nil
	}
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

func (s *ExportService) getItemResponses() ([]dto.ItemResponse, error) {
	if s.ItemRepository == nil {
		return []dto.ItemResponse{}, nil
	}
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

func (s *ExportService) getImageResponses() ([]dto.ImageResponse, error) {
	if s.ImageRepository == nil {
		return []dto.ImageResponse{}, nil
	}
	images, err := s.ImageRepository.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.ImageResponse, len(images))
	for i := range images {
		responses[i] = *dto.ToImageResponse(&images[i])
	}
	return responses, nil
}