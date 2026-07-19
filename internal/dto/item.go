package dto

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type CreateItemRequest struct {
	BrandID       *uint64        `json:"brand_id,omitempty"`
	CategoryID    uint64         `json:"category_id"`
	LocationID    *uint64        `json:"location_id,omitempty"`
	AssetCode     string         `json:"asset_code"`
	Name          string         `json:"name"`
	ItemCondition string         `json:"item_condition"`
	ItemStatus    string         `json:"item_status"`
	Notes         *string        `json:"notes,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

type UpdateItemRequest struct {
	BrandID       *uint64 `json:"brand_id,omitempty"`
	CategoryID    *uint64 `json:"category_id,omitempty"`
	LocationID    *uint64 `json:"location_id,omitempty"`
	AssetCode     *string `json:"asset_code,omitempty"`
	Name          *string `json:"name,omitempty"`
	ItemCondition *string `json:"item_condition,omitempty"`
	ItemStatus    *string `json:"item_status,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

type ItemResponse struct {
	ID            uint64         `json:"id"`
	BrandID       *uint64        `json:"brand_id,omitempty"`
	CategoryID    uint64         `json:"category_id"`
	LocationID    *uint64        `json:"location_id,omitempty"`
	AssetCode     string         `json:"asset_code"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	ItemCondition string         `json:"item_condition"`
	ItemStatus    string         `json:"item_status"`
	Notes         *string        `json:"notes,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (r CreateItemRequest) ToModel() *models.Item {
	name := strings.TrimSpace(r.Name)

	item := &models.Item{
		BrandID:       r.BrandID,
		CategoryID:    r.CategoryID,
		LocationID:    r.LocationID,
		AssetCode:     strings.TrimSpace(r.AssetCode),
		Name:          name,
		Slug:          slug.Make(name),
		ItemCondition: models.ItemCondition(r.ItemCondition),
		ItemStatus:    models.ItemStatus(r.ItemStatus),
	}

	if r.Notes != nil {
		notes := strings.TrimSpace(*r.Notes)
		item.Notes = sql.NullString{
			String: notes,
			Valid:  notes != "",
		}
	}

	return item
}

func (r UpdateItemRequest) Apply(item *models.Item) {
	if r.BrandID != nil {
		item.BrandID = r.BrandID
	}

	if r.CategoryID != nil {
		item.CategoryID = *r.CategoryID
	}

	if r.LocationID != nil {
		item.LocationID = r.LocationID
	}

	if r.AssetCode != nil {
		item.AssetCode = strings.TrimSpace(*r.AssetCode)
	}

	if r.Name != nil {
		name := strings.TrimSpace(*r.Name)
		if name != "" {
			item.Name = name
			item.Slug = slug.Make(name)
		}
	}

	if r.ItemCondition != nil {
		item.ItemCondition = models.ItemCondition(*r.ItemCondition)
	}

	if r.ItemStatus != nil {
		item.ItemStatus = models.ItemStatus(*r.ItemStatus)
	}

	if r.Notes != nil {
		notes := strings.TrimSpace(*r.Notes)
		item.Notes = sql.NullString{
			String: notes,
			Valid:  notes != "",
		}
	}
}

func ToItemResponse(item *models.Item) *ItemResponse {
	if item == nil {
		return nil
	}

	var notes *string
	if item.Notes.Valid {
		value := item.Notes.String
		notes = &value
	}

	return &ItemResponse{
		ID:            item.ID,
		BrandID:       item.BrandID,
		CategoryID:    item.CategoryID,
		LocationID:    item.LocationID,
		AssetCode:     item.AssetCode,
		Name:          item.Name,
		Slug:          item.Slug,
		ItemCondition: string(item.ItemCondition),
		ItemStatus:    string(item.ItemStatus),
		Notes:         notes,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}
