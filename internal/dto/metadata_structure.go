package dto

import (
	"time"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type CreateMetadataStructureRequest struct {
	Fields []models.MetadataField `json:"fields"`
}

type UpdateMetadataStructureRequest struct {
	Fields []models.MetadataField `json:"fields,omitempty"`
}

type MetadataStructureResponse struct {
	ID         uint64                 `json:"id" export:"ID"`
	CategoryID uint64                 `json:"category_id" export:"Category ID"`
	Fields     []models.MetadataField `json:"fields" export:"Fields"`
	Version    uint                   `json:"version" export:"Version"`
	CreatedAt  time.Time              `json:"created_at" export:"Created At"`
	UpdatedAt  time.Time              `json:"updated_at" export:"Updated At"`
}

func ToMetadataStructureResponse(structure *models.MetadataStructure) (*MetadataStructureResponse, error) {
	if structure == nil {
		return nil, nil
	}

	fields, err := structure.DecodeFields()
	if err != nil {
		return nil, err
	}

	return &MetadataStructureResponse{
		ID:         structure.ID,
		CategoryID: structure.CategoryID,
		Fields:     fields,
		Version:    structure.Version,
		CreatedAt:  structure.CreatedAt,
		UpdatedAt:  structure.UpdatedAt,
	}, nil
}
