package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	CATEGORY_TABLE_NAME = `table_category`
	CATEGORY_FIND_FIELDS = `
		id,
		name,
		slug,
		description,
		created_at,
		updated_at
	`
	CATEGORY_CREATE_FIELDS = `
		name,
		slug,
		description
	`
	CATEGORY_UPDATE_FIELDS = `
		name = ?,
		slug = ?,
		description = ?
	`
	CATEGORY_PLACEHOLDER = `(?, ?, ?)`
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category

	query := `
		SELECT ` + CATEGORY_FIND_FIELDS + `
		FROM ` + CATEGORY_TABLE_NAME + `
		ORDER BY name ASC
	`

	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) FindByID(id uint64) (*models.Category, error) {
	var category models.Category

	query := `
		SELECT ` + CATEGORY_FIND_FIELDS + `
		FROM ` + CATEGORY_TABLE_NAME + `
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.Get(&category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO ` + CATEGORY_TABLE_NAME + ` 
		(` + CATEGORY_CREATE_FIELDS + `)
		VALUES ` + CATEGORY_PLACEHOLDER + `
	`

	result, err := r.db.Exec(
		query,
		category.Name,
		category.Slug,
		category.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	category.ID = uint64(id)

	return nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE ` + CATEGORY_TABLE_NAME + `
		SET ` + CATEGORY_UPDATE_FIELDS + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ID,
	)

	return err
}

func (r *CategoryRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + CATEGORY_TABLE_NAME + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
