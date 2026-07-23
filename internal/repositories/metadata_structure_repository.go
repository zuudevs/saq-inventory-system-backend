package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	kMetadataStructureTableName  = `table_metadata_structure`
	kMetadataStructureFindFields = `
		id,
		category_id,
		fields,
		version,
		created_at,
		updated_at
	`
	kMetadataStructureCreateFields = `
		category_id,
		fields,
		version
	`
	kMetadataStructureUpdateFields = `
		category_id = ?,
		fields = ?,
		version = ?
	`
	kMetadataStructurePlaceholder = `(?, ?, ?)`
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
		SELECT ` + kMetadataStructureFindFields + `
		FROM ` + kMetadataStructureTableName + `
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
		INSERT INTO ` + kMetadataStructureTableName + ` 
		(` + kMetadataStructureCreateFields + `)
		VALUES ` + kMetadataStructurePlaceholder + `
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
		SELECT `+kMetadataStructureFindFields+`
		FROM `+kMetadataStructureTableName+`
		WHERE id = ?
		`,
		structure.ID,
	)
}

func (r *MetadataStructureRepository) Update(structure *models.MetadataStructure) error {
	query := `
		UPDATE ` + kMetadataStructureTableName + `
		SET ` + kMetadataStructureUpdateFields + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		structure.CategoryID,
		structure.Fields,
		structure.Version,
		structure.ID,
	)

	if err != nil {
		return err
	}

	return r.db.Get(
		structure,
		`
		SELECT `+kMetadataStructureFindFields+`
		FROM `+kMetadataStructureTableName+`
		WHERE id = ?
		`,
		structure.ID,
	)
}

func (r *MetadataStructureRepository) Delete(categoryID uint64) error {
	query := `
		DELETE FROM ` + kMetadataStructureTableName + `
		WHERE category_id = ?
	`

	_, err := r.db.Exec(query, categoryID)

	return err
}
