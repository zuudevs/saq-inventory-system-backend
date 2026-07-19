package services

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
	"github.com/zuudevs/saq-inventory-system-backend/internal/schema"
)

type ItemService struct {
	DB                          *sqlx.DB
	ItemRepository              *repositories.ItemRepository
	CategoryRepository          *repositories.CategoryRepository
	MetadataStructureRepository *repositories.MetadataStructureRepository
	MetadataRepository          *repositories.MetadataRepository
}

// Create menyimpan item baru. Jika kategorinya punya metadata structure
// terdaftar, payload req.Metadata divalidasi terhadap definisi field-nya
// (baca table_metadata_structure), lalu insert ke table_item dan insert ke
// table_<slug>_metadata dijalankan dalam satu SQL transaction yang sama —
// keduanya DML murni sehingga benar-benar atomic (beda dengan pembuatan
// tabel metadata di MetadataStructureService yang melibatkan DDL).
func (s *ItemService) Create(req dto.CreateItemRequest) (*dto.ItemResponse, error) {
	item := req.ToModel()

	if err := validateItem(item); err != nil {
		return nil, err
	}

	category, err := s.CategoryRepository.FindByID(item.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	structure, err := s.MetadataStructureRepository.FindByCategoryID(item.CategoryID)
	if err != nil {
		return nil, err
	}

	if structure == nil {
		if len(req.Metadata) > 0 {
			return nil, errors.New("this category has no metadata structure defined")
		}

		if err := s.ItemRepository.Create(item); err != nil {
			return nil, err
		}

		return dto.ToItemResponse(item), nil
	}

	fields, err := structure.DecodeFields()
	if err != nil {
		return nil, err
	}

	normalized, err := schema.ValidateMetadataPayload(fields, req.Metadata)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}

	if err := s.ItemRepository.CreateWithExecutor(tx, item); err != nil {
		tx.Rollback()
		return nil, err
	}

	tableName := schema.MetadataTableName(category.Slug)

	if err := s.MetadataRepository.InsertWithExecutor(tx, tableName, item.ID, normalized); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	response := dto.ToItemResponse(item)
	response.Metadata = normalized

	return response, nil
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

	response := dto.ToItemResponse(item)

	metadata, err := s.loadMetadata(item)
	if err != nil {
		return nil, err
	}
	response.Metadata = metadata

	return response, nil
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

	response := dto.ToItemResponse(item)

	metadata, err := s.loadMetadata(item)
	if err != nil {
		return nil, err
	}
	response.Metadata = metadata

	return response, nil
}

func (s *ItemService) Delete(id uint64) error {
	return s.ItemRepository.Delete(id)
}

// loadMetadata mengambil baris metadata milik sebuah item dari
// table_<slug>_metadata, kalau kategori item tersebut memang punya metadata
// structure terdaftar. Dipanggil dari FindByID dan Update supaya kedua
// endpoint itu ikut mengembalikan metadata dari DB, bukan cuma dari payload
// request seperti yang sebelumnya terjadi di Create.
//
// Mengembalikan (nil, nil) bila kategori tidak punya metadata structure,
// atau belum ada baris metadata untuk item ini (mis. race antara create item
// dan create structure, atau data lama sebelum structure ditambahkan).
func (s *ItemService) loadMetadata(item *models.Item) (map[string]any, error) {
	category, err := s.CategoryRepository.FindByID(item.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	structure, err := s.MetadataStructureRepository.FindByCategoryID(item.CategoryID)
	if err != nil {
		return nil, err
	}
	if structure == nil {
		return nil, nil
	}

	tableName := schema.MetadataTableName(category.Slug)

	row, err := s.MetadataRepository.FindByItemID(tableName, item.ID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	// `id`, `item_id`, `created_at`, `updated_at` adalah kolom housekeeping
	// tabel metadata, bukan bagian dari payload field yang didefinisikan
	// user di metadata structure — dibuang supaya bentuk response konsisten
	// dengan yang dikirim balik oleh Create (field structure saja).
	delete(row, "id")
	delete(row, "item_id")
	delete(row, "created_at")
	delete(row, "updated_at")

	return row, nil
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
