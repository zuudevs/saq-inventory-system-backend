package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	IMAGE_TABLE_NAME  = `table_image`
	IMAGE_FIND_FIELDS = `
		id,
		location_id,
		item_id,
		image_path,
		is_primary,
		created_at,
		updated_at
	`
	IMAGE_CREATE_FIELDS = `
		location_id,
		item_id,
		image_path,
		is_primary
	`
	IMAGE_UPDATE_FIELDS = `
		location_id = ?,
		item_id = ?,
		image_path = ?,
		is_primary = ?
	`
	IMAGE_PLACEHOLDER = `(?, ?, ?, ?)`
)

type ImageRepository struct {
	db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

func (r *ImageRepository) FindAll() ([]models.Image, error) {
	var images []models.Image

	query := `
		SELECT ` + IMAGE_FIND_FIELDS + `
		FROM ` + IMAGE_TABLE_NAME + `
		ORDER BY
			item_id IS NULL,
			COALESCE(item_id, location_id),
			is_primary DESC,
			id ASC
	`

	err := r.db.Select(&images, query)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepository) FindByItemID(itemID uint64) ([]models.Image, error) {
	var images []models.Image

	query := `
		SELECT ` + IMAGE_FIND_FIELDS + `
		FROM ` + IMAGE_TABLE_NAME + `
		WHERE item_id = ?
		ORDER BY is_primary DESC, id ASC
	`

	err := r.db.Select(&images, query, itemID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepository) FindByLocationID(locationID uint64) ([]models.Image, error) {
	var images []models.Image

	query := `
		SELECT ` + IMAGE_FIND_FIELDS + `
		FROM ` + IMAGE_TABLE_NAME + `
		WHERE location_id = ?
		ORDER BY is_primary DESC, id ASC
	`

	err := r.db.Select(&images, query, locationID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepository) FindByID(id uint64) (*models.Image, error) {
	var image models.Image

	query := `
		SELECT ` + IMAGE_FIND_FIELDS + `
		FROM ` + IMAGE_TABLE_NAME + `
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.Get(&image, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &image, nil
}

func (r *ImageRepository) Create(image *models.Image) error {
	return r.CreateWithExecutor(r.db, image)
}

// CreateWithExecutor sama seperti Create, tapi menerima sqlExecutor
// eksplisit (bisa *sqlx.DB atau *sqlx.Tx) supaya pemanggil bisa
// menyertakan operasi ini di dalam transaction yang lebih besar, mis.
// bersamaan dengan insert metadata dinamis di table_<slug>_metadata.
func (r *ImageRepository) CreateWithExecutor(exec sqlExecutor, image *models.Image) error {
	query := `
		INSERT INTO ` + IMAGE_TABLE_NAME + ` 
		(` + IMAGE_CREATE_FIELDS + `)
		VALUES ` + IMAGE_PLACEHOLDER + `
	`

	result, err := exec.Exec(
		query,
		image.LocationID,
		image.ItemID,
		image.ImagePath,
		image.IsPrimary,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	image.ID = uint64(id)

	return exec.Get(
		image,
		`
		SELECT `+IMAGE_FIND_FIELDS+`
		FROM `+IMAGE_TABLE_NAME+`
		WHERE id = ?
		`,
		image.ID,
	)
}

func (r *ImageRepository) Update(image *models.Image) error {
	return r.UpdateWithExecutor(r.db, image)
}

// UpdateWithExecutor sama seperti Update, tapi menerima sqlExecutor eksplisit
// supaya pemanggil bisa menyertakan operasi ini dalam transaction yang sama
// dengan UnsetPrimaryByItemIDWithExecutor/UnsetPrimaryByLocationIDWithExecutor
// — perlu atomic karena keduanya menyentuh unique partial index is_primary.
func (r *ImageRepository) UpdateWithExecutor(exec sqlExecutor, image *models.Image) error {
	query := `
		UPDATE ` + IMAGE_TABLE_NAME + `
		SET ` + IMAGE_UPDATE_FIELDS + `
		WHERE id = ?
	`

	_, err := exec.Exec(
		query,
		image.LocationID,
		image.ItemID,
		image.ImagePath,
		image.IsPrimary,
		image.ID,
	)

	if err != nil {
		return err
	}

	return exec.Get(
		image,
		`
		SELECT `+IMAGE_FIND_FIELDS+`
		FROM `+IMAGE_TABLE_NAME+`
		WHERE id = ?
		`,
		image.ID,
	)
}

// UnsetPrimaryByItemIDWithExecutor menghapus flag is_primary dari image lain
// milik item yang sama (excludeID dikecualikan supaya baris yang sedang
// diupdate tidak ikut ter-unset). Dipanggil dalam transaction sebelum
// insert/update sebuah image jadi is_primary = true, karena
// idx_image_item_primary adalah unique partial index — tanpa ini, insert
// atau update kedua akan gagal dengan UNIQUE constraint violation.
func (r *ImageRepository) UnsetPrimaryByItemIDWithExecutor(exec sqlExecutor, itemID uint64, excludeID uint64) error {
	query := `
		UPDATE ` + IMAGE_TABLE_NAME + `
		SET is_primary = 0
		WHERE item_id = ? AND is_primary = 1 AND id != ?
	`

	_, err := exec.Exec(query, itemID, excludeID)

	return err
}

// UnsetPrimaryByLocationIDWithExecutor adalah versi UnsetPrimaryByItemIDWithExecutor
// untuk idx_image_location_primary.
func (r *ImageRepository) UnsetPrimaryByLocationIDWithExecutor(exec sqlExecutor, locationID uint64, excludeID uint64) error {
	query := `
		UPDATE ` + IMAGE_TABLE_NAME + `
		SET is_primary = 0
		WHERE location_id = ? AND is_primary = 1 AND id != ?
	`

	_, err := exec.Exec(query, locationID, excludeID)

	return err
}

func (r *ImageRepository) Delete(id uint64) error {
	query := `
		DELETE FROM ` + IMAGE_TABLE_NAME + `
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
