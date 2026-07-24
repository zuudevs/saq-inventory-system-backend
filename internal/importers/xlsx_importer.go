package importers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/xuri/excelize/v2"
	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

type ImportData struct {
	Brands     []models.Brand
	Categories []models.Category
	Locations  []models.Location
	Items      []models.Item
	Images     []models.Image
}

var expectedHeaders = map[string][]string{
	"Brands":     {"ID", "Name", "Created At", "Updated At"},
	"Categories": {"ID", "Name", "Description", "Created At", "Updated At"},
	"Locations":  {"ID", "Name", "Room Code", "Description", "Created At", "Updated At"},
	"Items":      {"Brand ID", "ID", "Category ID", "Location ID", "Asset Code", "Name", "Item Condition", "Item Status", "Notes", "Created At", "Updated At"},
	"Images":     {"ID", "Location ID", "Item ID", "Image Path", "Is Primary", "Created At", "Updated At"},
}

// ParseAndValidateXLSX reads an XLSX file and performs sheet name, column header, and data type validations.
func ParseAndValidateXLSX(reader io.Reader) (*ImportData, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("invalid excel file: %w", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	existingSheets := make(map[string]bool)
	for _, name := range sheetList {
		existingSheets[name] = true
	}

	requiredSheets := []string{"Brands", "Categories", "Locations", "Items", "Images"}
	for _, reqSheet := range requiredSheets {
		if !existingSheets[reqSheet] {
			return nil, fmt.Errorf("missing required sheet: %s", reqSheet)
		}
	}

	data := &ImportData{}

	// 1. Validate and parse Brands
	brandsRows, err := f.GetRows("Brands")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet Brands: %w", err)
	}
	if err := validateHeaders("Brands", brandsRows, expectedHeaders["Brands"]); err != nil {
		return nil, err
	}
	brands, err := parseBrands(brandsRows)
	if err != nil {
		return nil, err
	}
	data.Brands = brands

	// 2. Validate and parse Categories
	catRows, err := f.GetRows("Categories")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet Categories: %w", err)
	}
	if err := validateHeaders("Categories", catRows, expectedHeaders["Categories"]); err != nil {
		return nil, err
	}
	categories, err := parseCategories(catRows)
	if err != nil {
		return nil, err
	}
	data.Categories = categories

	// 3. Validate and parse Locations
	locRows, err := f.GetRows("Locations")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet Locations: %w", err)
	}
	if err := validateHeaders("Locations", locRows, expectedHeaders["Locations"]); err != nil {
		return nil, err
	}
	locations, err := parseLocations(locRows)
	if err != nil {
		return nil, err
	}
	data.Locations = locations

	// 4. Validate and parse Items
	itemRows, err := f.GetRows("Items")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet Items: %w", err)
	}
	if err := validateHeaders("Items", itemRows, expectedHeaders["Items"]); err != nil {
		return nil, err
	}
	items, err := parseItems(itemRows)
	if err != nil {
		return nil, err
	}
	data.Items = items

	// 5. Validate and parse Images
	imgRows, err := f.GetRows("Images")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet Images: %w", err)
	}
	if err := validateHeaders("Images", imgRows, expectedHeaders["Images"]); err != nil {
		return nil, err
	}
	images, err := parseImages(imgRows)
	if err != nil {
		return nil, err
	}
	data.Images = images

	return data, nil
}

func validateHeaders(sheetName string, rows [][]string, expected []string) error {
	if len(rows) == 0 {
		return fmt.Errorf("sheet %q is empty", sheetName)
	}
	headerRow := rows[0]
	if len(headerRow) < len(expected) {
		return fmt.Errorf("sheet %q column header mismatch: expected headers %v, got %v", sheetName, expected, headerRow)
	}
	for i, exp := range expected {
		actual := strings.TrimSpace(headerRow[i])
		if actual != exp {
			return fmt.Errorf("sheet %q column header mismatch at index %d: expected %q, got %q", sheetName, i, exp, actual)
		}
	}
	return nil
}

func parseBrands(rows [][]string) ([]models.Brand, error) {
	var brands []models.Brand
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]
		if isEmptyRow(row) {
			continue
		}
		id, err := parseOptionalUint64(getCol(row, 0))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Brands", rIdx+1, "ID", err)
		}
		name := strings.TrimSpace(getCol(row, 1))
		if name == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Brands", rIdx+1, "Name")
		}

		var brandID uint64
		if id != nil {
			brandID = *id
		}
		createdAt := parseOptionalTime(getCol(row, 2))
		updatedAt := parseOptionalTime(getCol(row, 3))

		brands = append(brands, models.Brand{
			ID:        brandID,
			Name:      name,
			Slug:      slug.Make(name),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return brands, nil
}

func parseCategories(rows [][]string) ([]models.Category, error) {
	var categories []models.Category
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]
		if isEmptyRow(row) {
			continue
		}
		id, err := parseOptionalUint64(getCol(row, 0))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Categories", rIdx+1, "ID", err)
		}
		name := strings.TrimSpace(getCol(row, 1))
		if name == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Categories", rIdx+1, "Name")
		}

		desc := strings.TrimSpace(getCol(row, 2))
		var descNull sql.NullString
		if desc != "" {
			descNull = sql.NullString{String: desc, Valid: true}
		}

		var catID uint64
		if id != nil {
			catID = *id
		}
		createdAt := parseOptionalTime(getCol(row, 3))
		updatedAt := parseOptionalTime(getCol(row, 4))

		categories = append(categories, models.Category{
			ID:          catID,
			Name:        name,
			Slug:        slug.Make(name),
			Description: descNull,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}
	return categories, nil
}

func parseLocations(rows [][]string) ([]models.Location, error) {
	var locations []models.Location
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]
		if isEmptyRow(row) {
			continue
		}
		id, err := parseOptionalUint64(getCol(row, 0))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Locations", rIdx+1, "ID", err)
		}
		name := strings.TrimSpace(getCol(row, 1))
		if name == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Locations", rIdx+1, "Name")
		}

		roomCode := strings.TrimSpace(getCol(row, 2))
		var roomNull sql.NullString
		if roomCode != "" {
			roomNull = sql.NullString{String: roomCode, Valid: true}
		}

		desc := strings.TrimSpace(getCol(row, 3))
		var descNull sql.NullString
		if desc != "" {
			descNull = sql.NullString{String: desc, Valid: true}
		}

		var locID uint64
		if id != nil {
			locID = *id
		}
		createdAt := parseOptionalTime(getCol(row, 4))
		updatedAt := parseOptionalTime(getCol(row, 5))

		locations = append(locations, models.Location{
			ID:          locID,
			Name:        name,
			Slug:        slug.Make(name),
			RoomCode:    roomNull,
			Description: descNull,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}
	return locations, nil
}

func parseItems(rows [][]string) ([]models.Item, error) {
	var items []models.Item
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]
		if isEmptyRow(row) {
			continue
		}
		brandID, err := parseOptionalUint64(getCol(row, 0))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Items", rIdx+1, "Brand ID", err)
		}

		id, err := parseOptionalUint64(getCol(row, 1))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Items", rIdx+1, "ID", err)
		}

		catID, err := parseRequiredUint64(getCol(row, 2), "Category ID")
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d: %w", "Items", rIdx+1, err)
		}

		locID, err := parseOptionalUint64(getCol(row, 3))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Items", rIdx+1, "Location ID", err)
		}

		assetCode := strings.TrimSpace(getCol(row, 4))
		if assetCode == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Items", rIdx+1, "Asset Code")
		}

		name := strings.TrimSpace(getCol(row, 5))
		if name == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Items", rIdx+1, "Name")
		}

		itemCondStr := strings.TrimSpace(getCol(row, 6))
		if itemCondStr == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Items", rIdx+1, "Item Condition")
		}

		itemStatusStr := strings.TrimSpace(getCol(row, 7))
		if itemStatusStr == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Items", rIdx+1, "Item Status")
		}

		notes := strings.TrimSpace(getCol(row, 8))
		var notesNull sql.NullString
		if notes != "" {
			notesNull = sql.NullString{String: notes, Valid: true}
		}

		var itemID uint64
		if id != nil {
			itemID = *id
		}
		createdAt := parseOptionalTime(getCol(row, 9))
		updatedAt := parseOptionalTime(getCol(row, 10))

		items = append(items, models.Item{
			ID:            itemID,
			BrandID:       brandID,
			CategoryID:    catID,
			LocationID:    locID,
			AssetCode:     assetCode,
			Name:          name,
			Slug:          slug.Make(name),
			ItemCondition: models.ItemCondition(itemCondStr),
			ItemStatus:    models.ItemStatus(itemStatusStr),
			Notes:         notesNull,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}
	return items, nil
}

func parseImages(rows [][]string) ([]models.Image, error) {
	var images []models.Image
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]
		if isEmptyRow(row) {
			continue
		}
		id, err := parseOptionalUint64(getCol(row, 0))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Images", rIdx+1, "ID", err)
		}
		locID, err := parseOptionalUint64(getCol(row, 1))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Images", rIdx+1, "Location ID", err)
		}
		itemID, err := parseOptionalUint64(getCol(row, 2))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid uint64: %w", "Images", rIdx+1, "Item ID", err)
		}

		if (locID == nil && itemID == nil) || (locID != nil && itemID != nil) {
			return nil, fmt.Errorf("sheet %q row %d: image must have exactly one of Location ID or Item ID", "Images", rIdx+1)
		}

		imagePath := strings.TrimSpace(getCol(row, 3))
		if imagePath == "" {
			return nil, fmt.Errorf("sheet %q row %d col %q: required field is empty", "Images", rIdx+1, "Image Path")
		}

		isPrimary, err := parseBool(getCol(row, 4))
		if err != nil {
			return nil, fmt.Errorf("sheet %q row %d col %q: invalid boolean value: %w", "Images", rIdx+1, "Is Primary", err)
		}

		var imgID uint64
		if id != nil {
			imgID = *id
		}
		createdAt := parseOptionalTime(getCol(row, 5))
		updatedAt := parseOptionalTime(getCol(row, 6))

		images = append(images, models.Image{
			ID:         imgID,
			LocationID: locID,
			ItemID:     itemID,
			ImagePath:  imagePath,
			IsPrimary:  isPrimary,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}
	return images, nil
}

func getCol(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func isEmptyRow(row []string) bool {
	for _, col := range row {
		if strings.TrimSpace(col) != "" {
			return false
		}
	}
	return true
}

func parseOptionalUint64(val string) (*uint64, error) {
	if val == "" {
		return nil, nil
	}
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func parseRequiredUint64(val string, colName string) (uint64, error) {
	if val == "" {
		return 0, fmt.Errorf("required field %q is empty", colName)
	}
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("field %q must be a valid positive integer: %w", colName, err)
	}
	return n, nil
}

func parseBool(val string) (bool, error) {
	v := strings.ToLower(strings.TrimSpace(val))
	if v == "true" || v == "1" {
		return true, nil
	}
	if v == "false" || v == "0" || v == "" {
		return false, nil
	}
	return false, errors.New("must be true, false, 1, or 0")
}

func parseOptionalTime(val string) time.Time {
	val = strings.TrimSpace(val)
	if val == "" {
		return time.Time{}
	}
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02",
	}
	for _, fmtStr := range formats {
		if t, err := time.Parse(fmtStr, val); err == nil {
			return t
		}
	}
	return time.Time{}
}
