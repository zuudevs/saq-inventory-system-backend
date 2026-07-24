package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestExportItemsHandler(t *testing.T) {
	itemRepo := setupExportTestDB(t)

	exportService := &services.ExportService{
		ItemRepository: itemRepo,
	}
	handler := NewExportHandler(exportService)

	req := httptest.NewRequest(http.MethodGet, "/exports/items", nil)
	rr := httptest.NewRecorder()

	handler.ExportItems(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("Expected Content-Type 'text/csv', got %q", contentType)
	}

	contentDisp := rr.Header().Get("Content-Disposition")
	expectedDisp := "attachment; filename=items.csv"
	if contentDisp != expectedDisp {
		t.Errorf("Expected Content-Disposition %q, got %q", expectedDisp, contentDisp)
	}

	body := rr.Body.String()
	expectedHeader := "Brand ID,ID,Category ID,Location ID,Asset Code,Name,Item Condition,Item Status,Notes,Created At,Updated At"
	if !strings.Contains(body, expectedHeader) {
		t.Errorf("Expected body to contain header %q, got %q", expectedHeader, body)
	}

	expectedData := "1,1,2,3,LAP-001,MacBook Pro M3,good,active,Test laptop notes"
	if !strings.Contains(body, expectedData) {
		t.Errorf("Expected body to contain row data %q, got %q", expectedData, body)
	}
}
