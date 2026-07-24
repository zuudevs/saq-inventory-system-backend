package services

import (
	"bytes"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"github.com/zuudevs/saq-inventory-system-backend/internal/config"
)

func setupFullImportTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := config.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	tables := []string{
		`CREATE TABLE table_brand (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE table_category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE table_location (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			room_code TEXT,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE table_item (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			brand_id INTEGER,
			category_id INTEGER NOT NULL,
			location_id INTEGER,
			asset_code TEXT NOT NULL,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			item_condition TEXT NOT NULL,
			item_status TEXT NOT NULL,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE table_image (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			location_id INTEGER,
			item_id INTEGER,
			image_path TEXT NOT NULL,
			is_primary INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, stmt := range tables {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("Failed to create test schema table: %v", err)
		}
	}

	return db
}

func createImportXLSXBuffer(t *testing.T) *bytes.Buffer {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := "Sheet1"

	sheets := map[string][]string{
		"Brands":     {"ID", "Name", "Created At", "Updated At"},
		"Categories": {"ID", "Name", "Description", "Created At", "Updated At"},
		"Locations":  {"ID", "Name", "Room Code", "Description", "Created At", "Updated At"},
		"Items":      {"Brand ID", "ID", "Category ID", "Location ID", "Asset Code", "Name", "Item Condition", "Item Status", "Notes", "Created At", "Updated At"},
		"Images":     {"ID", "Location ID", "Item ID", "Image Path", "Is Primary", "Created At", "Updated At"},
	}

	idx := 0
	for sheetName, headers := range sheets {
		if idx == 0 {
			_ = f.SetSheetName(defaultSheet, sheetName)
		} else {
			_, _ = f.NewSheet(sheetName)
		}
		for colIdx, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			_ = f.SetCellValue(sheetName, cell, h)
		}
		idx++
	}

	_ = f.SetCellValue("Brands", "A2", "1")
	_ = f.SetCellValue("Brands", "B2", "Dell")

	_ = f.SetCellValue("Categories", "A2", "1")
	_ = f.SetCellValue("Categories", "B2", "Laptops")

	_ = f.SetCellValue("Locations", "A2", "1")
	_ = f.SetCellValue("Locations", "B2", "Room 101")

	_ = f.SetCellValue("Items", "A2", "1")
	_ = f.SetCellValue("Items", "B2", "1")
	_ = f.SetCellValue("Items", "C2", "1")
	_ = f.SetCellValue("Items", "D2", "1")
	_ = f.SetCellValue("Items", "E2", "LAP-DELL-01")
	_ = f.SetCellValue("Items", "F2", "Dell XPS 15")
	_ = f.SetCellValue("Items", "G2", "good")
	_ = f.SetCellValue("Items", "H2", "active")

	_ = f.SetCellValue("Images", "A2", "1")
	_ = f.SetCellValue("Images", "C2", "1")
	_ = f.SetCellValue("Images", "D2", "images/dell.jpg")
	_ = f.SetCellValue("Images", "E2", "1")

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		t.Fatalf("Failed to create test xlsx buffer: %v", err)
	}
	return &buf
}

func TestImportService_ImportXLSX(t *testing.T) {
	db := setupFullImportTestDB(t)
	defer db.Close()

	svc := &ImportService{
		DB: db,
	}

	buf := createImportXLSXBuffer(t)

	summary, err := svc.ImportXLSX(buf)
	if err != nil {
		t.Fatalf("ImportXLSX failed: %v", err)
	}

	if summary.BrandsImported != 1 || summary.CategoriesImported != 1 || summary.LocationsImported != 1 || summary.ItemsImported != 1 || summary.ImagesImported != 1 {
		t.Errorf("Unexpected summary counts: %+v", summary)
	}

	var brandName string
	if err := db.Get(&brandName, "SELECT name FROM table_brand WHERE id = 1"); err != nil || brandName != "Dell" {
		t.Errorf("Expected imported brand 'Dell', got %q, err: %v", brandName, err)
	}

	var itemName string
	if err := db.Get(&itemName, "SELECT name FROM table_item WHERE id = 1"); err != nil || itemName != "Dell XPS 15" {
		t.Errorf("Expected imported item 'Dell XPS 15', got %q, err: %v", itemName, err)
	}
}

