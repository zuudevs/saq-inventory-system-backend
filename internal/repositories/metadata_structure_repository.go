package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	METADATA_STRUCTURE_TABLE_NAME  = `table_metadata_structure`
	METADATA_STRUCTURE_FIND_FIELDS = `
		id,
		category_id,
		fields,
		version,
		created_at,
		updated_at
	`
	METADATA_STRUCTURE_CREATE_FIELDS = `
		category_id,
		fields,
		version
	`
	METADATA_STRUCTURE_PLACEHOLDER = `(?, ?, ?)`
)

type MetadataStructureRepository struct {
	db *sqlx.DB
}

func NewMetadataStructureRepository(db *sqlx.DB) *MetadataStructureRepository {
	return &MetadataStructureRepository{
		db: db,
	}
}

func (r *MetadataStructureRepository) FindByCategoryID(categoryID uint64) (*models.MetadataStructure, error) {
	var structure models.MetadataStructure

	query := `
		SELECT ` + METADATA_STRUCTURE_FIND_FIELDS + `
		FROM ` + METADATA_STRUCTURE_TABLE_NAME + `
		WHERE category_id = ?
		LIMIT 1
	`

	err := r.db.Get(&structure, query, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &structure, nil
}

func (r *MetadataStructureRepository) Create(structure *models.MetadataStructure) error {
	query := `
		INSERT INTO ` + METADATA_STRUCTURE_TABLE_NAME + ` 
		(` + METADATA_STRUCTURE_CREATE_FIELDS + `)
		VALUES ` + METADATA_STRUCTURE_PLACEHOLDER + `
	`

	result, err := r.db.Exec(
		query,
		structure.CategoryID,
		structure.Fields,
		structure.Version,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	structure.ID = uint64(id)

	return r.db.Get(
		structure,
		`
		SELECT `+METADATA_STRUCTURE_FIND_FIELDS+`
		FROM `+METADATA_STRUCTURE_TABLE_NAME+`
		WHERE id = ?
		`,
		structure.ID,
	)
}

func (r *MetadataStructureRepository) Delete(categoryID uint64) error {
	query := `
		DELETE FROM ` + METADATA_STRUCTURE_TABLE_NAME + `
		WHERE category_id = ?
	`

	_, err := r.db.Exec(query, categoryID)

	return err
}
