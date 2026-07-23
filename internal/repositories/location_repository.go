package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	kLocationTableName  = `table_location`
	kLocationFindFields = `
		id,
		name,
		slug,
		room_code,
		description,
		created_at,
		updated_at
	`
	kLocationCreateFields = `
		name,
		slug,
		room_code,
		description
	`
	kLocationUpdateFields = `
		name = ?,
		slug = ?,
		room_code = ?,
		description = ?
	`
	kLocationPlaceholder = `(?, ?, ?, ?)`
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
		SELECT ` + kLocationFindFields + `
		FROM ` + kLocationTableName + `
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
		SELECT ` + kLocationFindFields + `
		FROM ` + kLocationTableName + `
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
		INSERT INTO ` + kLocationTableName + ` 
		(` + kLocationCreateFields + `)
		VALUES ` + kLocationPlaceholder + `
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

	if err != nil {
		return err
	}

	return r.db.Get(
		location,
		`
		SELECT `+kLocationFindFields+`
		FROM `+kLocationTableName+`
		WHERE id = ?
		`,
		location.ID,
	)
}

func (r *LocationRepository) Update(location *models.Location) error {
	query := `
		UPDATE ` + kLocationTableName + `
		SET ` + kLocationUpdateFields + `
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

	if err != nil {
		return err
	}

	return r.db.Get(
		location,
		`
		SELECT `+kLocationFindFields+`
		FROM `+kLocationTableName+`
		WHERE id = ?
		`,
		location.ID,
	)
}

func (r *LocationRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + kLocationTableName + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
