package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
	"github.com/zuudevs/saq-inventory-system-backend/internal/schema"
)

const MaxMetadataFields = 50

type MetadataStructureService struct {
	MetadataStructureRepository *repositories.MetadataStructureRepository
	CategoryRepository          *repositories.CategoryRepository
	SchemaService               *schema.Service
}

// Create menjalankan alur lengkap pendefinisian metadata untuk sebuah
// kategori:
//  1. Validasi kategori ada & belum punya metadata structure.
//  2. Validasi field-field yang diminta user (nama, tipe, default, dsb).
//  3. Generate & eksekusi CREATE TABLE lewat SchemaService.
//  4. Simpan definisi field ke table_metadata_structure.
//
// Langkah 3 dieksekusi di luar transaction (lihat Service.CreateMetadataTable),
// sehingga tidak benar-benar atomic dengan langkah 4. Bila langkah 4 gagal
// setelah tabel berhasil dibuat, service melakukan compensating action
// (DROP TABLE) supaya table_metadata_structure dan tabel fisik tetap
// konsisten satu sama lain.
func (s *MetadataStructureService) Create(
	categoryID uint64,
	req dto.CreateMetadataStructureRequest,
) (*dto.MetadataStructureResponse, error) {
	category, err := s.CategoryRepository.FindByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	existing, err := s.MetadataStructureRepository.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("metadata structure already exists for this category")
	}

	if err := validateFields(req.Fields); err != nil {
		return nil, err
	}

	if err := s.SchemaService.CreateMetadataTable(category.Slug, req.Fields); err != nil {
		return nil, fmt.Errorf("failed to create metadata table: %w", err)
	}

	fieldsJSON, err := models.EncodeMetadataFields(req.Fields)
	if err != nil {
		_ = s.SchemaService.DropMetadataTable(category.Slug)
		return nil, err
	}

	structure := &models.MetadataStructure{
		CategoryID: categoryID,
		Fields:     fieldsJSON,
		Version:    1,
	}

	if err := s.MetadataStructureRepository.Create(structure); err != nil {
		_ = s.SchemaService.DropMetadataTable(category.Slug)
		return nil, fmt.Errorf("failed to save metadata structure: %w", err)
	}

	return dto.ToMetadataStructureResponse(structure)
}

func (s *MetadataStructureService) Update(
	categoryID uint64, 
	req dto.UpdateMetadataStructureRequest,
) (*dto.MetadataStructureResponse, error) {
	category, err := s.CategoryRepository.FindByID(categoryID)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	existing, err := s.MetadataStructureRepository.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("metadata structure is not already exists for this category")
	}

	if err := validateFields(req.Fields); err != nil {
		return nil, err
	}

	oldFields, err := existing.DecodeFields()
	if err != nil {
		return nil, fmt.Errorf("failed to decode existing metadata fields: %w", err)
	}

	if err := s.SchemaService.UpdateMetadataTable(category.Slug, oldFields, req.Fields); err != nil {
		return nil, fmt.Errorf("failed to update metadata table schema: %w", err)
	}

	fieldsJSON, err := models.EncodeMetadataFields(req.Fields)
	if err != nil {
		// Compensating action: revert schema changes
		_ = s.SchemaService.UpdateMetadataTable(category.Slug, req.Fields, oldFields)
		return nil, err
	}

	structure := &models.MetadataStructure{
		ID:         existing.ID,
		CategoryID: categoryID,
		Fields:     fieldsJSON,
		Version:    existing.Version + 1,
	}

	if err := s.MetadataStructureRepository.Update(structure); err != nil {
		// Compensating action: revert schema changes
		_ = s.SchemaService.UpdateMetadataTable(category.Slug, req.Fields, oldFields)
		return nil, fmt.Errorf("failed to save metadata structure: %w", err)
	}

	return dto.ToMetadataStructureResponse(structure)
}

func (s *MetadataStructureService) FindByCategoryID(categoryID uint64) (*dto.MetadataStructureResponse, error) {
	structure, err := s.MetadataStructureRepository.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}
	if structure == nil {
		return nil, nil
	}

	return dto.ToMetadataStructureResponse(structure)
}

// validateFields melakukan validasi request-level (bukan validasi
// identifier/DDL — itu tanggung jawab schema.BuildCreateTableSQL). Fokus
// di sini adalah memastikan payload yang dikirim user lengkap dan masuk
// akal secara bisnis sebelum kita repot-repot menyusun DDL.
func validateFields(fields []models.MetadataField) error {
	if len(fields) == 0 {
		return errors.New("at least one metadata field is required")
	}

	if len(fields) > MaxMetadataFields {
		return fmt.Errorf("metadata fields exceed maximum limit of %d", MaxMetadataFields)
	}

	seen := make(map[string]struct{}, len(fields))

	for _, field := range fields {
		name := strings.TrimSpace(field.Name)
		if name == "" {
			return errors.New("field name is required")
		}

		if _, dup := seen[name]; dup {
			return fmt.Errorf("duplicate field name: %s", name)
		}
		seen[name] = struct{}{}

		if strings.TrimSpace(field.Label) == "" {
			return fmt.Errorf("label is required for field: %s", name)
		}

		switch field.Type {
		case models.MetadataFieldTypeString,
			models.MetadataFieldTypeText,
			models.MetadataFieldTypeInt,
			models.MetadataFieldTypeFloat,
			models.MetadataFieldTypeBool,
			models.MetadataFieldTypeDate,
			models.MetadataFieldTypeDatetime,
			models.MetadataFieldTypeEnum:
			// valid
		default:
			return fmt.Errorf("unsupported field type '%s' for field: %s", field.Type, name)
		}

		if field.Type == models.MetadataFieldTypeEnum && len(field.Options) == 0 {
			return fmt.Errorf("enum field '%s' requires at least one option", name)
		}
	}

	return nil
}