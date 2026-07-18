package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	LOCATION_TABLE_NAME = `table_location`
	LOCATION_FIND_FIELDS = `
		id,
		name,
		slug,
		room_code,
		description,
		created_at,
		updated_at
	`
	LOCATION_CREATE_FIELDS = `
		name,
		slug,
		room_code,
		description
	`
	LOCATION_UPDATE_FIELDS = `
		name = ?,
		slug = ?,
		room_code = ?,
		description = ?
	`
	LOCATION_PLACEHOLDER = `(?, ?, ?, ?)`
)

type LocationRepository struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) FindAll() ([]models.Location, error) {
	var locations []models.Location

	query := `
		SELECT ` + LOCATION_FIND_FIELDS + `
		FROM ` + LOCATION_TABLE_NAME + `
		ORDER BY name ASC
	`

	err := r.db.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) FindByID(id uint64) (*models.Location, error) {
	var location models.Location

	query := `
		SELECT ` + LOCATION_FIND_FIELDS + `
		FROM ` + LOCATION_TABLE_NAME + `
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.Get(&location, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &location, nil
}

func (r *LocationRepository) Create(location *models.Location) error {
	query := `
		INSERT INTO ` + LOCATION_TABLE_NAME + ` 
		(` + LOCATION_CREATE_FIELDS + `)
		VALUES ` + LOCATION_PLACEHOLDER + `
	`

	result, err := r.db.Exec(
		query,
		location.Name,
		location.Slug,
		location.RoomCode,
		location.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	location.ID = uint64(id)

	return nil
}

func (r *LocationRepository) Update(location *models.Location) error {
	query := `
		UPDATE ` + LOCATION_TABLE_NAME + `
		SET ` + LOCATION_UPDATE_FIELDS + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		location.Name,
		location.Slug,
		location.RoomCode,
		location.Description,
		location.ID,
	)

	return err
}

func (r *LocationRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + LOCATION_TABLE_NAME + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
