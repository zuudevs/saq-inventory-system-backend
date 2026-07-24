package services

import (
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
}

// ExportItemsToCSV fetches all items and exports them in CSV format to the given writer.
func (s *ExportService) ExportItemsToCSV(writer io.Writer) error {
	items, err := s.ItemRepository.FindAll()
	if err != nil {
		return err
	}

	itemResponses := make([]dto.ItemResponse, len(items))
	for i := range items {
		itemResponses[i] = *dto.ToItemResponse(&items[i])
	}

	return exporters.ExportCSV(writer, itemResponses)
}

// ExportItemsToXLSX fetches all items and exports them in XLSX format to the given writer.
func (s *ExportService) ExportItemsToXLSX(writer io.Writer) error {
	items, err := s.ItemRepository.FindAll()
	if err != nil {
		return err
	}

	itemResponses := make([]dto.ItemResponse, len(items))
	for i := range items {
		itemResponses[i] = *dto.ToItemResponse(&items[i])
	}

	return exporters.ExportXLSX(writer, itemResponses)
}