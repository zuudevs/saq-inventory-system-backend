package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	BRAND_TABLE_NAME  = `table_brand`
	BRAND_FIND_FIELDS = `
		id,
		name,
		slug,
		created_at,
		updated_at
	`
	BRAND_CREATE_FIELDS = `
		name,
		slug
	`
	BRAND_UPDATE_FIELDS = `
		name = ?,
		slug = ?
	`
	BRAND_PLACEHOLDER = `(?, ?)`
)

type BrandRepository struct {
	db *sqlx.DB
}

func NewBrandRepository(db *sqlx.DB) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

func (r *BrandRepository) FindAll() ([]models.Brand, error) {
	var brands []models.Brand

	query := `
		SELECT ` + BRAND_FIND_FIELDS + `
		FROM ` + BRAND_TABLE_NAME + `
		ORDER BY name ASC
	`

	err := r.db.Select(&brands, query)
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *BrandRepository) FindByID(id uint64) (*models.Brand, error) {
	var brand models.Brand

	query := `
		SELECT ` + BRAND_FIND_FIELDS + `
		FROM ` + BRAND_TABLE_NAME + `
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.Get(&brand, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &brand, nil
}

func (r *BrandRepository) Create(brand *models.Brand) error {
	query := `
		INSERT INTO ` + BRAND_TABLE_NAME + ` 
		(` + BRAND_CREATE_FIELDS + `)
		VALUES ` + BRAND_PLACEHOLDER + `
	`

	result, err := r.db.Exec(
		query,
		brand.Name,
		brand.Slug,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	brand.ID = uint64(id)

	return r.db.Get(
		brand,
		`
		SELECT `+BRAND_FIND_FIELDS+`
		FROM `+BRAND_TABLE_NAME+`
		WHERE id = ?
		`,
		brand.ID,
	)
}

func (r *BrandRepository) Update(brand *models.Brand) error {
	query := `
		UPDATE ` + BRAND_TABLE_NAME + `
		SET ` + BRAND_UPDATE_FIELDS + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		brand.Name,
		brand.Slug,
		brand.ID,
	)

	if err != nil {
		return err
	}

	return r.db.Get(
		brand,
		`
		SELECT `+BRAND_FIND_FIELDS+`
		FROM `+BRAND_TABLE_NAME+`
		WHERE id = ?
		`,
		brand.ID,
	)
}

func (r *BrandRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + BRAND_TABLE_NAME + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
