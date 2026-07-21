package dto

import (
	"strings"
	"time"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

// CreateImageRequest menerima persis salah satu dari LocationID atau ItemID
// (divalidasi di ImageService), sesuai CHECK constraint table_image yang
// mewajibkan image adalah milik location ATAU item, tidak boleh keduanya
// atau tidak keduanya.
type CreateImageRequest struct {
	LocationID *uint64 `json:"location_id,omitempty"`
	ItemID     *uint64 `json:"item_id,omitempty"`
	ImagePath  string  `json:"image_path"`
	IsPrimary  bool    `json:"is_primary,omitempty"`
}

// UpdateImageRequest sengaja tidak mengizinkan pemindahan kepemilikan image
// (location_id/item_id) — reparenting image ke owner lain di luar scope
// endpoint update biasa dan lebih aman dilakukan lewat delete + create ulang.
type UpdateImageRequest struct {
	ImagePath *string `json:"image_path,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
}

type ImageResponse struct {
	ID         uint64    `json:"id"`
	LocationID *uint64   `json:"location_id,omitempty"`
	ItemID     *uint64   `json:"item_id,omitempty"`
	ImagePath  string    `json:"image_path"`
	IsPrimary  bool      `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r CreateImageRequest) ToModel() *models.Image {
	return &models.Image{
		LocationID: r.LocationID,
		ItemID:     r.ItemID,
		ImagePath:  strings.TrimSpace(r.ImagePath),
		IsPrimary:  r.IsPrimary,
	}
}

func (r UpdateImageRequest) Apply(image *models.Image) {
	if r.ImagePath != nil {
		path := strings.TrimSpace(*r.ImagePath)
		if path != "" {
			image.ImagePath = path
		}
	}

	if r.IsPrimary != nil {
		image.IsPrimary = *r.IsPrimary
	}
}

func ToImageResponse(image *models.Image) *ImageResponse {
	if image == nil {
		return nil
	}

	return &ImageResponse{
		ID:         image.ID,
		LocationID: image.LocationID,
		ItemID:     image.ItemID,
		ImagePath:  image.ImagePath,
		IsPrimary:  image.IsPrimary,
		CreatedAt:  image.CreatedAt,
		UpdatedAt:  image.UpdatedAt,
	}
}
