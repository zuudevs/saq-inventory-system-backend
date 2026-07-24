package services

import (
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/importers"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

type ImportService struct {
	DB                 *sqlx.DB
	BrandRepository    *repositories.BrandRepository
	CategoryRepository *repositories.CategoryRepository
	ItemRepository     *repositories.ItemRepository
	LocationRepository *repositories.LocationRepository
	ImageRepository    *repositories.ImageRepository
}

func (s *ImportService) ImportXLSX(reader io.Reader) (*dto.ImportSummary, error) {
	data, err := importers.ParseAndValidateXLSX(reader)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	summary := &dto.ImportSummary{}

	brandIDMap := make(map[uint64]uint64)
	catIDMap := make(map[uint64]uint64)
	locIDMap := make(map[uint64]uint64)
	itemIDMap := make(map[uint64]uint64)

	// 1. Insert Brands
	for _, b := range data.Brands {
		oldID := b.ID
		res, err := tx.Exec("INSERT INTO table_brand (name, slug) VALUES (?, ?)", b.Name, b.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to insert brand %q: %w", b.Name, err)
		}
		newID, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get brand last insert id: %w", err)
		}
		if oldID > 0 {
			brandIDMap[oldID] = uint64(newID)
		}
		summary.BrandsImported++
	}

	// 2. Insert Categories
	for _, c := range data.Categories {
		oldID := c.ID
		res, err := tx.Exec("INSERT INTO table_category (name, slug, description) VALUES (?, ?, ?)", c.Name, c.Slug, c.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to insert category %q: %w", c.Name, err)
		}
		newID, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get category last insert id: %w", err)
		}
		if oldID > 0 {
			catIDMap[oldID] = uint64(newID)
		}
		summary.CategoriesImported++
	}

	// 3. Insert Locations
	for _, l := range data.Locations {
		oldID := l.ID
		res, err := tx.Exec("INSERT INTO table_location (name, slug, room_code, description) VALUES (?, ?, ?, ?)", l.Name, l.Slug, l.RoomCode, l.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to insert location %q: %w", l.Name, err)
		}
		newID, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get location last insert id: %w", err)
		}
		if oldID > 0 {
			locIDMap[oldID] = uint64(newID)
		}
		summary.LocationsImported++
	}

	// 4. Insert Items
	for _, item := range data.Items {
		oldID := item.ID

		var brandID *uint64
		if item.BrandID != nil {
			if mapped, ok := brandIDMap[*item.BrandID]; ok {
				brandID = &mapped
			} else {
				brandID = item.BrandID
			}
			var exists int
			if err := tx.Get(&exists, "SELECT COUNT(1) FROM table_brand WHERE id = ?", *brandID); err != nil || exists == 0 {
				return nil, fmt.Errorf("item %q references non-existent brand_id %d", item.Name, *brandID)
			}
		}

		var catID uint64
		if mapped, ok := catIDMap[item.CategoryID]; ok {
			catID = mapped
		} else {
			catID = item.CategoryID
		}
		var catExists int
		if err := tx.Get(&catExists, "SELECT COUNT(1) FROM table_category WHERE id = ?", catID); err != nil || catExists == 0 {
			return nil, fmt.Errorf("item %q references non-existent category_id %d", item.Name, catID)
		}

		var locID *uint64
		if item.LocationID != nil {
			if mapped, ok := locIDMap[*item.LocationID]; ok {
				locID = &mapped
			} else {
				locID = item.LocationID
			}
			var exists int
			if err := tx.Get(&exists, "SELECT COUNT(1) FROM table_location WHERE id = ?", *locID); err != nil || exists == 0 {
				return nil, fmt.Errorf("item %q references non-existent location_id %d", item.Name, *locID)
			}
		}

		res, err := tx.Exec(`
			INSERT INTO table_item (brand_id, category_id, location_id, asset_code, name, slug, item_condition, item_status, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, brandID, catID, locID, item.AssetCode, item.Name, item.Slug, item.ItemCondition, item.ItemStatus, item.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to insert item %q: %w", item.Name, err)
		}
		newID, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get item last insert id: %w", err)
		}
		if oldID > 0 {
			itemIDMap[oldID] = uint64(newID)
		}
		summary.ItemsImported++
	}

	// 5. Insert Images
	for _, img := range data.Images {
		var locID *uint64
		if img.LocationID != nil {
			if mapped, ok := locIDMap[*img.LocationID]; ok {
				locID = &mapped
			} else {
				locID = img.LocationID
			}
			var exists int
			if err := tx.Get(&exists, "SELECT COUNT(1) FROM table_location WHERE id = ?", *locID); err != nil || exists == 0 {
				return nil, fmt.Errorf("image references non-existent location_id %d", *locID)
			}
		}

		var itemID *uint64
		if img.ItemID != nil {
			if mapped, ok := itemIDMap[*img.ItemID]; ok {
				itemID = &mapped
			} else {
				itemID = img.ItemID
			}
			var exists int
			if err := tx.Get(&exists, "SELECT COUNT(1) FROM table_item WHERE id = ?", *itemID); err != nil || exists == 0 {
				return nil, fmt.Errorf("image references non-existent item_id %d", *itemID)
			}
		}

		_, err := tx.Exec(`
			INSERT INTO table_image (location_id, item_id, image_path, is_primary)
			VALUES (?, ?, ?, ?)
		`, locID, itemID, img.ImagePath, img.IsPrimary)
		if err != nil {
			return nil, fmt.Errorf("failed to insert image: %w", err)
		}
		summary.ImagesImported++
	}

	summary.TotalImported = summary.BrandsImported + summary.CategoriesImported + summary.LocationsImported + summary.ItemsImported + summary.ImagesImported

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit import transaction: %w", err)
	}
	tx = nil

	return summary, nil
}
