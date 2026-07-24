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

	// 1. Insert/Upsert Brands
	for _, b := range data.Brands {
		oldID := b.ID
		var targetID uint64
		if oldID > 0 {
			_, err := tx.Exec(`
				INSERT INTO table_brand (id, name, slug) VALUES (?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET name=EXCLUDED.name, slug=EXCLUDED.slug
			`, oldID, b.Name, b.Slug)
			if err != nil {
				return nil, fmt.Errorf("failed to import brand %q: %w", b.Name, err)
			}
			targetID = oldID
		} else {
			res, err := tx.Exec(`
				INSERT INTO table_brand (name, slug) VALUES (?, ?)
				ON CONFLICT(slug) DO UPDATE SET name=EXCLUDED.name
			`, b.Name, b.Slug)
			if err != nil {
				return nil, fmt.Errorf("failed to import brand %q: %w", b.Name, err)
			}
			newID, _ := res.LastInsertId()
			if newID > 0 {
				targetID = uint64(newID)
			} else {
				_ = tx.Get(&targetID, "SELECT id FROM table_brand WHERE slug = ?", b.Slug)
			}
		}
		if oldID > 0 {
			brandIDMap[oldID] = targetID
		}
		summary.BrandsImported++
	}

	// 2. Insert/Upsert Categories
	for _, c := range data.Categories {
		oldID := c.ID
		var targetID uint64
		if oldID > 0 {
			_, err := tx.Exec(`
				INSERT INTO table_category (id, name, slug, description) VALUES (?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET name=EXCLUDED.name, slug=EXCLUDED.slug, description=EXCLUDED.description
			`, oldID, c.Name, c.Slug, c.Description)
			if err != nil {
				return nil, fmt.Errorf("failed to import category %q: %w", c.Name, err)
			}
			targetID = oldID
		} else {
			res, err := tx.Exec(`
				INSERT INTO table_category (name, slug, description) VALUES (?, ?, ?)
				ON CONFLICT(slug) DO UPDATE SET name=EXCLUDED.name, description=EXCLUDED.description
			`, c.Name, c.Slug, c.Description)
			if err != nil {
				return nil, fmt.Errorf("failed to import category %q: %w", c.Name, err)
			}
			newID, _ := res.LastInsertId()
			if newID > 0 {
				targetID = uint64(newID)
			} else {
				_ = tx.Get(&targetID, "SELECT id FROM table_category WHERE slug = ?", c.Slug)
			}
		}
		if oldID > 0 {
			catIDMap[oldID] = targetID
		}
		summary.CategoriesImported++
	}

	// 3. Insert/Upsert Locations
	for _, l := range data.Locations {
		oldID := l.ID
		var targetID uint64
		if oldID > 0 {
			_, err := tx.Exec(`
				INSERT INTO table_location (id, name, slug, room_code, description) VALUES (?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET name=EXCLUDED.name, slug=EXCLUDED.slug, room_code=EXCLUDED.room_code, description=EXCLUDED.description
			`, oldID, l.Name, l.Slug, l.RoomCode, l.Description)
			if err != nil {
				return nil, fmt.Errorf("failed to import location %q: %w", l.Name, err)
			}
			targetID = oldID
		} else {
			res, err := tx.Exec(`
				INSERT INTO table_location (name, slug, room_code, description) VALUES (?, ?, ?, ?)
				ON CONFLICT(slug) DO UPDATE SET name=EXCLUDED.name, room_code=EXCLUDED.room_code, description=EXCLUDED.description
			`, l.Name, l.Slug, l.RoomCode, l.Description)
			if err != nil {
				return nil, fmt.Errorf("failed to import location %q: %w", l.Name, err)
			}
			newID, _ := res.LastInsertId()
			if newID > 0 {
				targetID = uint64(newID)
			} else {
				_ = tx.Get(&targetID, "SELECT id FROM table_location WHERE slug = ?", l.Slug)
			}
		}
		if oldID > 0 {
			locIDMap[oldID] = targetID
		}
		summary.LocationsImported++
	}

	// 4. Insert/Upsert Items
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

		var targetID uint64
		if oldID > 0 {
			_, err := tx.Exec(`
				INSERT INTO table_item (id, brand_id, category_id, location_id, asset_code, name, slug, item_condition, item_status, notes)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET brand_id=EXCLUDED.brand_id, category_id=EXCLUDED.category_id, location_id=EXCLUDED.location_id, asset_code=EXCLUDED.asset_code, name=EXCLUDED.name, slug=EXCLUDED.slug, item_condition=EXCLUDED.item_condition, item_status=EXCLUDED.item_status, notes=EXCLUDED.notes
			`, oldID, brandID, catID, locID, item.AssetCode, item.Name, item.Slug, item.ItemCondition, item.ItemStatus, item.Notes)
			if err != nil {
				return nil, fmt.Errorf("failed to import item %q: %w", item.Name, err)
			}
			targetID = oldID
		} else {
			res, err := tx.Exec(`
				INSERT INTO table_item (brand_id, category_id, location_id, asset_code, name, slug, item_condition, item_status, notes)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(asset_code) DO UPDATE SET brand_id=EXCLUDED.brand_id, category_id=EXCLUDED.category_id, location_id=EXCLUDED.location_id, name=EXCLUDED.name, slug=EXCLUDED.slug, item_condition=EXCLUDED.item_condition, item_status=EXCLUDED.item_status, notes=EXCLUDED.notes
			`, brandID, catID, locID, item.AssetCode, item.Name, item.Slug, item.ItemCondition, item.ItemStatus, item.Notes)
			if err != nil {
				return nil, fmt.Errorf("failed to import item %q: %w", item.Name, err)
			}
			newID, _ := res.LastInsertId()
			if newID > 0 {
				targetID = uint64(newID)
			} else {
				_ = tx.Get(&targetID, "SELECT id FROM table_item WHERE asset_code = ?", item.AssetCode)
			}
		}
		if oldID > 0 {
			itemIDMap[oldID] = targetID
		}
		summary.ItemsImported++
	}

	// 5. Insert/Upsert Images
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

		if img.ID > 0 {
			_, err := tx.Exec(`
				INSERT INTO table_image (id, location_id, item_id, image_path, is_primary)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET location_id=EXCLUDED.location_id, item_id=EXCLUDED.item_id, image_path=EXCLUDED.image_path, is_primary=EXCLUDED.is_primary
			`, img.ID, locID, itemID, img.ImagePath, img.IsPrimary)
			if err != nil {
				return nil, fmt.Errorf("failed to import image: %w", err)
			}
		} else {
			_, err := tx.Exec(`
				INSERT INTO table_image (location_id, item_id, image_path, is_primary)
				VALUES (?, ?, ?, ?)
			`, locID, itemID, img.ImagePath, img.IsPrimary)
			if err != nil {
				return nil, fmt.Errorf("failed to import image: %w", err)
			}
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
