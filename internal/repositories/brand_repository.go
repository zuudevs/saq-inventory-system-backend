package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	kBrandTableName  = `table_brand`
	kBrandFindFields = `
		id,
		name,
		slug,
		created_at,
		updated_at
	`
	kBrandCreateFields = `
		name,
		slug
	`
	kBrandUpdateFields = `
		name = ?,
		slug = ?
	`
	kBrandPlaceholder = `(?, ?)`
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
		SELECT ` + kBrandFindFields + `
		FROM ` + kBrandTableName + `
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
		SELECT ` + kBrandFindFields + `
		FROM ` + kBrandTableName + `
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
		INSERT INTO ` + kBrandTableName + ` 
		(` + kBrandCreateFields + `)
		VALUES ` + kBrandPlaceholder + `
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
		SELECT `+kBrandFindFields+`
		FROM `+kBrandTableName+`
		WHERE id = ?
		`,
		brand.ID,
	)
}

func (r *BrandRepository) Update(brand *models.Brand) error {
	query := `
		UPDATE ` + kBrandTableName + `
		SET ` + kBrandUpdateFields + `
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
		SELECT `+kBrandFindFields+`
		FROM `+kBrandTableName+`
		WHERE id = ?
		`,
		brand.ID,
	)
}

func (r *BrandRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + kBrandTableName + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
