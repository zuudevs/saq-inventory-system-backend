package dto

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CategoryResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r CreateCategoryRequest) ToModel() *models.Category {
	name := strings.TrimSpace(r.Name)

	category := &models.Category{
		Name: name,
		Slug: slug.Make(name),
	}

	if r.Description != nil {
		desc := strings.TrimSpace(*r.Description)
		category.Description = sql.NullString{
			String: desc,
			Valid:  desc != "",
		}
	}

	return category
}

func (r UpdateCategoryRequest) Apply(category *models.Category) {
	if r.Name != nil {
		name := strings.TrimSpace(*r.Name)
		if name != "" {
			category.Name = name
			category.Slug = slug.Make(name)
		}
	}

	if r.Description != nil {
		desc := strings.TrimSpace(*r.Description)
		if desc == "" {
			category.Description = sql.NullString{}
			return
		}

		category.Description = sql.NullString{
			String: desc,
			Valid:  true,
		}
	}
}

func ToCategoryResponse(category *models.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	var description *string
	if category.Description.Valid {
		desc := category.Description.String
		description = &desc
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
