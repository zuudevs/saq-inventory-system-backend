package importers

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func createTestWorkbook(t *testing.T, modifyFunc func(f *excelize.File)) *bytes.Buffer {
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

	// Add sample valid rows
	_ = f.SetCellValue("Brands", "A2", "1")
	_ = f.SetCellValue("Brands", "B2", "Apple")

	_ = f.SetCellValue("Categories", "A2", "1")
	_ = f.SetCellValue("Categories", "B2", "Electronics")

	_ = f.SetCellValue("Locations", "A2", "1")
	_ = f.SetCellValue("Locations", "B2", "HQ Main Office")

	_ = f.SetCellValue("Items", "A2", "1")
	_ = f.SetCellValue("Items", "B2", "10")
	_ = f.SetCellValue("Items", "C2", "1")
	_ = f.SetCellValue("Items", "D2", "1")
	_ = f.SetCellValue("Items", "E2", "AST-001")
	_ = f.SetCellValue("Items", "F2", "MacBook Pro")
	_ = f.SetCellValue("Items", "G2", "good")
	_ = f.SetCellValue("Items", "H2", "active")

	_ = f.SetCellValue("Images", "A2", "1")
	_ = f.SetCellValue("Images", "C2", "10") // Item ID = 10
	_ = f.SetCellValue("Images", "D2", "images/macbook.png")
	_ = f.SetCellValue("Images", "E2", "true")

	if modifyFunc != nil {
		modifyFunc(f)
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		t.Fatalf("Failed to write test workbook: %v", err)
	}
	return &buf
}

func TestParseAndValidateXLSX_Valid(t *testing.T) {
	buf := createTestWorkbook(t, nil)
	data, err := ParseAndValidateXLSX(buf)
	if err != nil {
		t.Fatalf("Expected valid parse, got error: %v", err)
	}

	if len(data.Brands) != 1 || data.Brands[0].Name != "Apple" {
		t.Errorf("Unexpected Brands data: %v", data.Brands)
	}
	if len(data.Categories) != 1 || data.Categories[0].Name != "Electronics" {
		t.Errorf("Unexpected Categories data: %v", data.Categories)
	}
	if len(data.Items) != 1 || data.Items[0].AssetCode != "AST-001" {
		t.Errorf("Unexpected Items data: %v", data.Items)
	}
	if len(data.Images) != 1 || data.Images[0].ImagePath != "images/macbook.png" {
		t.Errorf("Unexpected Images data: %v", data.Images)
	}
}

func TestParseAndValidateXLSX_MissingSheet(t *testing.T) {
	buf := createTestWorkbook(t, func(f *excelize.File) {
		_ = f.DeleteSheet("Images")
	})

	_, err := ParseAndValidateXLSX(buf)
	if err == nil || !strings.Contains(err.Error(), "missing required sheet: Images") {
		t.Errorf("Expected missing sheet error for Images, got %v", err)
	}
}

func TestParseAndValidateXLSX_HeaderMismatch(t *testing.T) {
	buf := createTestWorkbook(t, func(f *excelize.File) {
		_ = f.SetCellValue("Brands", "B1", "WrongHeader")
	})

	_, err := ParseAndValidateXLSX(buf)
	if err == nil || !strings.Contains(err.Error(), "column header mismatch") {
		t.Errorf("Expected header mismatch error, got %v", err)
	}
}

func TestParseAndValidateXLSX_InvalidDataType(t *testing.T) {
	// Invalid Uint64 in Items.CategoryID
	buf := createTestWorkbook(t, func(f *excelize.File) {
		_ = f.SetCellValue("Items", "C2", "not_a_number")
	})

	_, err := ParseAndValidateXLSX(buf)
	if err == nil || !strings.Contains(err.Error(), "must be a valid positive integer") {
		t.Errorf("Expected integer validation error, got %v", err)
	}
}

func TestParseAndValidateXLSX_InvalidImageOwner(t *testing.T) {
	// Both Location ID and Item ID set
	buf := createTestWorkbook(t, func(f *excelize.File) {
		_ = f.SetCellValue("Images", "B2", "1") // Location ID = 1
		_ = f.SetCellValue("Images", "C2", "1") // Item ID = 1
	})

	_, err := ParseAndValidateXLSX(buf)
	if err == nil || !strings.Contains(err.Error(), "must have exactly one of Location ID or Item ID") {
		t.Errorf("Expected image owner validation error, got %v", err)
	}
}
