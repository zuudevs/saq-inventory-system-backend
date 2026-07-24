package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"github.com/zuudevs/saq-inventory-system-backend/internal/config"
	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
)

func setupImportHandlerTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := config.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	tables := []string{
		`CREATE TABLE table_brand (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, slug TEXT NOT NULL, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE table_category (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, slug TEXT NOT NULL, description TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE table_location (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, slug TEXT NOT NULL, room_code TEXT, description TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE table_item (id INTEGER PRIMARY KEY AUTOINCREMENT, brand_id INTEGER, category_id INTEGER NOT NULL, location_id INTEGER, asset_code TEXT NOT NULL, name TEXT NOT NULL, slug TEXT NOT NULL, item_condition TEXT NOT NULL, item_status TEXT NOT NULL, notes TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE table_image (id INTEGER PRIMARY KEY AUTOINCREMENT, location_id INTEGER, item_id INTEGER, image_path TEXT NOT NULL, is_primary INTEGER NOT NULL DEFAULT 0, created_at DATETIME, updated_at DATETIME);`,
	}

	for _, stmt := range tables {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("Failed to create test schema table: %v", err)
		}
	}

	return db
}

func createImportFormFile(t *testing.T) (*bytes.Buffer, string) {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()

	sheets := map[string][]string{
		"Brands":     {"ID", "Name", "Created At", "Updated At"},
		"Categories": {"ID", "Name", "Description", "Created At", "Updated At"},
		"Locations":  {"ID", "Name", "Room Code", "Description", "Created At", "Updated At"},
		"Items":      {"Brand ID", "ID", "Category ID", "Location ID", "Asset Code", "Name", "Item Condition", "Item Status", "Notes", "Created At", "Updated At"},
		"Images":     {"ID", "Location ID", "Item ID", "Image Path", "Is Primary", "Created At", "Updated At"},
	}

	defaultSheet := "Sheet1"
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
	_ = f.SetCellValue("Brands", "B2", "Sony")

	_ = f.SetCellValue("Categories", "A2", "1")
	_ = f.SetCellValue("Categories", "B2", "Audio")

	_ = f.SetCellValue("Locations", "A2", "1")
	_ = f.SetCellValue("Locations", "B2", "Studio A")

	_ = f.SetCellValue("Items", "A2", "1")
	_ = f.SetCellValue("Items", "B2", "1")
	_ = f.SetCellValue("Items", "C2", "1")
	_ = f.SetCellValue("Items", "D2", "1")
	_ = f.SetCellValue("Items", "E2", "MIC-01")
	_ = f.SetCellValue("Items", "F2", "Sony Mic")
	_ = f.SetCellValue("Items", "G2", "good")
	_ = f.SetCellValue("Items", "H2", "active")

	_ = f.SetCellValue("Images", "A2", "1")
	_ = f.SetCellValue("Images", "C2", "1")
	_ = f.SetCellValue("Images", "D2", "images/mic.png")
	_ = f.SetCellValue("Images", "E2", "1")

	var xlsxBuf bytes.Buffer
	if _, err := f.WriteTo(&xlsxBuf); err != nil {
		t.Fatalf("Failed to write xlsx: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "import.xlsx")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, _ = part.Write(xlsxBuf.Bytes())
	_ = writer.Close()

	return body, writer.FormDataContentType()
}

func TestImportXLSXHandler(t *testing.T) {
	db := setupImportHandlerTestDB(t)
	defer db.Close()

	importService := &services.ImportService{
		DB: db,
	}
	handler := NewImportHandler(importService)

	body, contentType := createImportFormFile(t)

	req := httptest.NewRequest(http.MethodPost, "/imports/xlsx", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	handler.ImportXLSX(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d, body: %s", rr.Code, rr.Body.String())
	}
}

