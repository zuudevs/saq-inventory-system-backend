package handlers

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
	"github.com/zuudevs/saq-inventory-system-backend/internal/config"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
)

func setupExportTestDB(t *testing.T) *repositories.ItemRepository {
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

func TestExportCSVHandler(t *testing.T) {
	itemRepo := setupExportTestDB(t)

	exportService := &services.ExportService{
		ItemRepository: itemRepo,
	}
	handler := NewExportHandler(exportService)

	req := httptest.NewRequest(http.MethodGet, "/exports/csv", nil)
	rr := httptest.NewRecorder()

	handler.ExportCSV(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/zip" {
		t.Errorf("Expected Content-Type 'application/zip', got %q", contentType)
	}

	contentDisp := rr.Header().Get("Content-Disposition")
	expectedDisp := "attachment; filename=exports.zip"
	if contentDisp != expectedDisp {
		t.Errorf("Expected Content-Disposition %q, got %q", expectedDisp, contentDisp)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(rr.Body.Bytes()), int64(rr.Body.Len()))
	if err != nil {
		t.Fatalf("Failed to parse response body as ZIP archive: %v", err)
	}

	expectedFiles := map[string]bool{
		"brands.csv":     false,
		"categories.csv": false,
		"locations.csv":  false,
		"items.csv":      false,
		"images.csv":     false,
	}

	for _, f := range zipReader.File {
		if _, exists := expectedFiles[f.Name]; exists {
			expectedFiles[f.Name] = true
		}
		if f.Name == "items.csv" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("Failed to open items.csv inside zip: %v", err)
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatalf("Failed to read items.csv: %v", err)
			}
			if !strings.Contains(string(content), "LAP-001") {
				t.Errorf("Expected items.csv to contain 'LAP-001', got %s", string(content))
			}
		}
	}

	for fileName, found := range expectedFiles {
		if !found {
			t.Errorf("Expected zip archive to contain file %s", fileName)
		}
	}
}

func TestExportXLSXHandler(t *testing.T) {
	itemRepo := setupExportTestDB(t)

	exportService := &services.ExportService{
		ItemRepository: itemRepo,
	}
	handler := NewExportHandler(exportService)

	req := httptest.NewRequest(http.MethodGet, "/exports/xlsx", nil)
	rr := httptest.NewRecorder()

	handler.ExportXLSX(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	expectedContentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type %q, got %q", expectedContentType, contentType)
	}

	contentDisp := rr.Header().Get("Content-Disposition")
	expectedDisp := "attachment; filename=exports.xlsx"
	if contentDisp != expectedDisp {
		t.Errorf("Expected Content-Disposition %q, got %q", expectedDisp, contentDisp)
	}

	f, err := excelize.OpenReader(rr.Body)
	if err != nil {
		t.Fatalf("Failed to parse XLSX body: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Items")
	if err != nil {
		t.Fatalf("Failed to read rows from Items sheet: %v", err)
	}

	if len(rows) < 2 {
		t.Fatalf("Expected at least 2 rows (1 header + 1 data), got %d", len(rows))
	}
}

