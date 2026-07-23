package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	kCategoryTableName  = `table_category`
	kCategoryFindFields = `
		id,
		name,
		slug,
		description,
		created_at,
		updated_at
	`
	kCategoryCreateFields = `
		name,
		slug,
		description
	`
	kCategoryUpdateFields = `
		name = ?,
		slug = ?,
		description = ?
	`
	kCategoryPlaceholder = `(?, ?, ?)`
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
		SELECT ` + kCategoryFindFields + `
		FROM ` + kCategoryTableName + `
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
		SELECT ` + kCategoryFindFields + `
		FROM ` + kCategoryTableName + `
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
		INSERT INTO ` + kCategoryTableName + ` 
		(` + kCategoryCreateFields + `)
		VALUES ` + kCategoryPlaceholder + `
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

	return r.db.Get(
		category,
		`
		SELECT `+kCategoryFindFields+`
		FROM `+kCategoryTableName+`
		WHERE id = ?
		`,
		category.ID,
	)
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE ` + kCategoryTableName + `
		SET ` + kCategoryUpdateFields + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ID,
	)

	if err != nil {
		return err
	}

	return r.db.Get(
		category,
		`
		SELECT `+kCategoryFindFields+`
		FROM `+kCategoryTableName+`
		WHERE id = ?
		`,
		category.ID,
	)
}

func (r *CategoryRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + kCategoryTableName + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
