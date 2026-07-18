package dto

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type CreateLocationRequest struct {
	Name        string  `json:"name"`
	RoomCode    *string `json:"room_code,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateLocationRequest struct {
	Name        *string `json:"name,omitempty"`
	RoomCode    *string `json:"room_code,omitempty"`
	Description *string `json:"description,omitempty"`
}

type LocationResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	RoomCode    *string   `json:"room_code,omitempty"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r CreateLocationRequest) ToModel() *models.Location {
	name := strings.TrimSpace(r.Name)

	location := &models.Location{
		Name: name,
		Slug: slug.Make(name),
	}

	if r.RoomCode != nil {
		roomCode := strings.TrimSpace(*r.RoomCode)
		location.RoomCode = sql.NullString{
			String: roomCode,
			Valid:  roomCode != "",
		}
	}

	if r.Description != nil {
		description := strings.TrimSpace(*r.Description)
		location.Description = sql.NullString{
			String: description,
			Valid:  description != "",
		}
	}

	return location
}

func (r UpdateLocationRequest) Apply(location *models.Location) {
	if r.Name != nil {
		name := strings.TrimSpace(*r.Name)
		if name != "" {
			location.Name = name
			location.Slug = slug.Make(name)
		}
	}

	if r.RoomCode != nil {
		roomCode := strings.TrimSpace(*r.RoomCode)
		location.RoomCode = sql.NullString{
			String: roomCode,
			Valid:  roomCode != "",
		}
	}

	if r.Description != nil {
		description := strings.TrimSpace(*r.Description)
		location.Description = sql.NullString{
			String: description,
			Valid:  description != "",
		}
	}
}

func ToLocationResponse(location *models.Location) *LocationResponse {
	if location == nil {
		return nil
	}

	var roomCode *string
	if location.RoomCode.Valid {
		value := location.RoomCode.String
		roomCode = &value
	}

	var description *string
	if location.Description.Valid {
		value := location.Description.String
		description = &value
	}

	return &LocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		Slug:        location.Slug,
		RoomCode:    roomCode,
		Description: description,
		CreatedAt:   location.CreatedAt,
		UpdatedAt:   location.UpdatedAt,
	}
}
