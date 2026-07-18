package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	ITEM_TABLE_NAME = `table_item`
	ITEM_FIND_FIELDS = `
		id,
		brand_id,
		category_id,
		location_id,
		asset_code,
		name,
		slug,
		item_condition,
		item_status,
		notes,
		created_at,
		updated_at
	`
	ITEM_CREATE_FIELDS = `
		brand_id,
		category_id,
		location_id,
		asset_code,
		name,
		slug,
		item_condition,
		item_status,
		notes
	`
	ITEM_UPDATE_FIELDS = `
		brand_id = ?,
		category_id = ?,
		location_id = ?,
		asset_code = ?,
		name = ?,
		slug = ?,
		item_condition = ?,
		item_status = ?,
		notes = ?
	`
	ITEM_PLACEHOLDER = `(?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) FindAll() ([]models.Item, error) {
	var items []models.Item

	query := `
		SELECT ` + ITEM_FIND_FIELDS + `
		FROM ` + ITEM_TABLE_NAME + `
		ORDER BY name ASC
	`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) FindByID(id uint64) (*models.Item, error) {
	var item models.Item

	query := `
		SELECT ` + ITEM_FIND_FIELDS + `
		FROM ` + ITEM_TABLE_NAME + `
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.Get(&item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) Create(item *models.Item) error {
	query := `
		INSERT INTO ` + ITEM_TABLE_NAME + ` 
		(` + ITEM_CREATE_FIELDS + `)
		VALUES ` + ITEM_PLACEHOLDER + `
	`

	result, err := r.db.Exec(
		query,
		item.BrandID,
		item.CategoryID,
		item.LocationID,
		item.AssetCode,
		item.Name,
		item.Slug,
		item.ItemCondition,
		item.ItemStatus,
		item.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = uint64(id)

	return nil
}

func (r *ItemRepository) Update(item *models.Item) error {
	query := `
		UPDATE ` + ITEM_TABLE_NAME + `
		SET ` + ITEM_UPDATE_FIELDS + `
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		item.BrandID,
		item.CategoryID,
		item.LocationID,
		item.AssetCode,
		item.Name,
		item.Slug,
		item.ItemCondition,
		item.ItemStatus,
		item.Notes,
		item.ID,
	)

	return err
}

func (r *ItemRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + ITEM_TABLE_NAME + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
