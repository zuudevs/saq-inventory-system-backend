package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	kItemTableName  = `table_item`
	kItemFindFields = `
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
	kItemCreateFields = `
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
	kItemUpdateFields = `
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
	kItemPlaceholder = `(?, ?, ?, ?, ?, ?, ?, ?, ?)`
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
		SELECT ` + kItemFindFields + `
		FROM ` + kItemTableName + `
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
		SELECT ` + kItemFindFields + `
		FROM ` + kItemTableName + `
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
	return r.CreateWithExecutor(r.db, item)
}

// CreateWithExecutor sama seperti Create, tapi menerima sqlExecutor
// eksplisit (bisa *sqlx.DB atau *sqlx.Tx) supaya pemanggil bisa
// menyertakan operasi ini di dalam transaction yang lebih besar, mis.
// bersamaan dengan insert metadata dinamis di table_<slug>_metadata.
func (r *ItemRepository) CreateWithExecutor(exec sqlExecutor, item *models.Item) error {
	query := `
		INSERT INTO ` + kItemTableName + ` 
		(` + kItemCreateFields + `)
		VALUES ` + kItemPlaceholder + `
	`

	result, err := exec.Exec(
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

	return exec.Get(
		item,
		`
		SELECT `+kItemFindFields+`
		FROM `+kItemTableName+`
		WHERE id = ?
		`,
		item.ID,
	)
}

func (r *ItemRepository) Update(item *models.Item) error {
	query := `
		UPDATE ` + kItemTableName + `
		SET ` + kItemUpdateFields + `
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

	if err != nil {
		return err
	}

	return r.db.Get(
		item,
		`
		SELECT `+kItemFindFields+`
		FROM `+kItemTableName+`
		WHERE id = ?
		`,
		item.ID,
	)
}

func (r *ItemRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + kItemTableName + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
