package services

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
	"github.com/zuudevs/saq-inventory-system-backend/internal/config"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
)

func setupServiceTestDB(t *testing.T) *repositories.ItemRepository {
	t.Helper()
	db, err := config.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	createStmt := `
	CREATE TABLE table_item (
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
	);`

	if _, err := db.Exec(createStmt); err != nil {
		t.Fatalf("Failed to create table_item: %v", err)
	}

	insertStmt := `
	INSERT INTO table_item (brand_id, category_id, location_id, asset_code, name, slug, item_condition, item_status, notes)
	VALUES (1, 2, 3, 'LAP-001', 'MacBook Pro M3', 'macbook-pro-m3', 'good', 'active', 'Test laptop notes');`

	if _, err := db.Exec(insertStmt); err != nil {
		t.Fatalf("Failed to insert mock item: %v", err)
	}

	return repositories.NewItemRepository(db)
}

func TestExportItemsToXLSX(t *testing.T) {
	itemRepo := setupServiceTestDB(t)
	exportService := &ExportService{
		ItemRepository: itemRepo,
	}

	var buf bytes.Buffer
	err := exportService.ExportItemsToXLSX(&buf)
	if err != nil {
		t.Fatalf("ExportItemsToXLSX failed: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("Failed to parse output as Excel file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		t.Fatalf("Failed to read rows from Sheet1: %v", err)
	}

	if len(rows) < 2 {
		t.Fatalf("Expected at least 2 rows (1 header + 1 data), got %d", len(rows))
	}

	if rows[1][4] != "LAP-001" {
		t.Errorf("Expected asset code 'LAP-001' at row 2 col 5, got %q", rows[1][4])
	}
}
