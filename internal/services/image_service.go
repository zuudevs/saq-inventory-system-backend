package services

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
	"github.com/zuudevs/saq-inventory-system-backend/internal/utils"
)

type ImageService struct {
	DB                 *sqlx.DB
	StoragePath        string
	ImageRepository    *repositories.ImageRepository
	ItemRepository     *repositories.ItemRepository
	LocationRepository *repositories.LocationRepository
}

// Create menyimpan image baru. Owner-nya (item atau location) divalidasi
// benar-benar ada, lalu kalau image ini di-set sebagai primary, unset
// primary lama milik owner yang sama dan insert dijalankan dalam satu
// transaction — wajib atomic karena idx_image_item_primary /
// idx_image_location_primary adalah unique partial index, jadi dua baris
// is_primary = 1 untuk owner yang sama akan ditolak oleh SQLite.
//
// ImagePath di sini diasumsikan sudah berupa hasil POST /images/upload
// (path relatif di dalam StoragePath) — Create tidak melakukan upload file,
// cuma menyimpan pointer ke file yang sudah ada di disk.
func (s *ImageService) Create(req dto.CreateImageRequest) (*dto.ImageResponse, error) {
	image := req.ToModel()

	if err := validateImageOwner(image); err != nil {
		return nil, err
	}

	if strings.TrimSpace(image.ImagePath) == "" {
		return nil, errors.New("image_path is required")
	}

	if err := s.ensureOwnerExists(image); err != nil {
		return nil, err
	}

	if !image.IsPrimary {
		if err := s.ImageRepository.Create(image); err != nil {
			return nil, err
		}

		return dto.ToImageResponse(image), nil
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}

	if err := s.unsetExistingPrimary(tx, image, 0); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.ImageRepository.CreateWithExecutor(tx, image); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dto.ToImageResponse(image), nil
}

func (s *ImageService) FindAll() ([]dto.ImageResponse, error) {
	images, err := s.ImageRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return toImageResponses(images), nil
}

func (s *ImageService) FindByItemID(itemID uint64) ([]dto.ImageResponse, error) {
	images, err := s.ImageRepository.FindByItemID(itemID)
	if err != nil {
		return nil, err
	}

	return toImageResponses(images), nil
}

func (s *ImageService) FindByLocationID(locationID uint64) ([]dto.ImageResponse, error) {
	images, err := s.ImageRepository.FindByLocationID(locationID)
	if err != nil {
		return nil, err
	}

	return toImageResponses(images), nil
}

func (s *ImageService) FindByID(id uint64) (*dto.ImageResponse, error) {
	image, err := s.ImageRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, nil
	}

	return dto.ToImageResponse(image), nil
}

// Update tidak mengizinkan perpindahan owner (lihat dto.UpdateImageRequest),
// jadi owner image tetap sama dengan sebelumnya — hanya image_path dan
// is_primary yang bisa berubah. Sama seperti Create, perubahan is_primary
// jadi true dijalankan dalam transaction bersama unset primary lama.
//
// Kalau image_path berubah (client upload file baru buat gantiin yang lama),
// file fisik lama dihapus best-effort setelah record berhasil diupdate —
// supaya tidak ada file yatim menumpuk di storage.
func (s *ImageService) Update(id uint64, req dto.UpdateImageRequest) (*dto.ImageResponse, error) {
	image, err := s.ImageRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, nil
	}

	oldImagePath := image.ImagePath
	wasPrimary := image.IsPrimary

	req.Apply(image)

	if strings.TrimSpace(image.ImagePath) == "" {
		return nil, errors.New("image_path is required")
	}

	if image.IsPrimary && !wasPrimary {
		tx, err := s.DB.Beginx()
		if err != nil {
			return nil, err
		}

		if err := s.unsetExistingPrimary(tx, image, image.ID); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := s.ImageRepository.UpdateWithExecutor(tx, image); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}
	} else {
		if err := s.ImageRepository.Update(image); err != nil {
			return nil, err
		}
	}

	if image.ImagePath != oldImagePath {
		_ = utils.DeleteImageFile(s.StoragePath, oldImagePath)
	}

	return dto.ToImageResponse(image), nil
}

// Delete menghapus record table_image, lalu best-effort menghapus file
// fisiknya juga. Kegagalan hapus file fisik sengaja tidak membatalkan
// (rollback) penghapusan record — kalau dibalik, file yang sudah kepencet
// manual duluan bakal bikin record tidak pernah bisa dihapus.
func (s *ImageService) Delete(id uint64) error {
	image, err := s.ImageRepository.FindByID(id)
	if err != nil {
		return err
	}
	if image == nil {
		return nil
	}

	if err := s.ImageRepository.Delete(id); err != nil {
		return err
	}

	_ = utils.DeleteImageFile(s.StoragePath, image.ImagePath)

	return nil
}

func (s *ImageService) ensureOwnerExists(image *models.Image) error {
	if image.ItemID != nil {
		item, err := s.ItemRepository.FindByID(*image.ItemID)
		if err != nil {
			return err
		}
		if item == nil {
			return errors.New("item not found")
		}
	}

	if image.LocationID != nil {
		location, err := s.LocationRepository.FindByID(*image.LocationID)
		if err != nil {
			return err
		}
		if location == nil {
			return errors.New("location not found")
		}
	}

	return nil
}

func (s *ImageService) unsetExistingPrimary(tx *sqlx.Tx, image *models.Image, excludeID uint64) error {
	if image.ItemID != nil {
		return s.ImageRepository.UnsetPrimaryByItemIDWithExecutor(tx, *image.ItemID, excludeID)
	}

	return s.ImageRepository.UnsetPrimaryByLocationIDWithExecutor(tx, *image.LocationID, excludeID)
}

func validateImageOwner(image *models.Image) error {
	if image.ItemID == nil && image.LocationID == nil {
		return errors.New("either item_id or location_id is required")
	}

	if image.ItemID != nil && image.LocationID != nil {
		return errors.New("image cannot belong to both an item and a location")
	}

	return nil
}

func toImageResponses(images []models.Image) []dto.ImageResponse {
	responses := make([]dto.ImageResponse, len(images))
	for i := range images {
		responses[i] = *dto.ToImageResponse(&images[i])
	}

	return responses
}
