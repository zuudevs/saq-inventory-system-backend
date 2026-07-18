package dto

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type CreateBrandRequest struct {
	Name string `json:"name"`
}

type UpdateBrandRequest struct {
	Name *string `json:"name,omitempty"`
}

type BrandResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r CreateBrandRequest) ToModel() *models.Brand {
	name := strings.TrimSpace(r.Name)

	brand := &models.Brand{
		Name: name,
		Slug: slug.Make(name),
	}

	return brand
}

func (r UpdateBrandRequest) Apply(brand *models.Brand) {
	if r.Name != nil {
		name := strings.TrimSpace(*r.Name)
		if name != "" {
			brand.Name = name
			brand.Slug = slug.Make(name)
		}
	}
}

func ToBrandResponse(brand *models.Brand) *BrandResponse {
	if brand == nil {
		return nil
	}

	return &BrandResponse{
		ID:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}
