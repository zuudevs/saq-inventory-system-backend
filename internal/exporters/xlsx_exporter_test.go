package exporters

import (
	"bytes"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestExportXLSX_StructSlice(t *testing.T) {
	optVal := "Hello Opt"
	now := time.Date(2026, 7, 24, 18, 0, 0, 0, time.UTC)
	data := []DummyStruct{
		{
			ID:        1,
			Name:      "Item One",
			IgnoreMe:  "Secret",
			Optional:  &optVal,
			ZeroValue: 42,
			CreatedAt: now,
		},
		{
			ID:        2,
			Name:      "Item Two",
			IgnoreMe:  "Secret 2",
			Optional:  nil,
			ZeroValue: 0,
			CreatedAt: time.Time{},
		},
	}

	var buf bytes.Buffer
	err := ExportXLSX(&buf, data)
	if err != nil {
		t.Fatalf("ExportXLSX failed: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("Failed to parse output as excel file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		t.Fatalf("Failed to get rows from Sheet1: %v", err)
	}

	if len(rows) != 3 {
		t.Fatalf("Expected 3 rows (1 header + 2 data rows), got %d", len(rows))
	}

	expectedHeaders := []string{"ID", "Name", "Optional Value", "Zero", "Created At"}
	for i, h := range expectedHeaders {
		if i >= len(rows[0]) || rows[0][i] != h {
			t.Errorf("Expected header %d to be %q, got %q", i, h, rows[0][i])
		}
	}

	expectedRow1 := []string{"1", "Item One", "Hello Opt", "42", "2026-07-24 18:00:00"}
	for i, val := range expectedRow1 {
		if i >= len(rows[1]) || rows[1][i] != val {
			t.Errorf("Expected row 1 col %d to be %q, got %q", i, val, rows[1][i])
		}
	}

	expectedRow2 := []string{"2", "Item Two", "", "0", ""}
	for i, val := range expectedRow2 {
		var actualVal string
		if i < len(rows[2]) {
			actualVal = rows[2][i]
		}
		if actualVal != val {
			t.Errorf("Expected row 2 col %d to be %q, got %q", i, val, actualVal)
		}
	}
}

func TestExportXLSX_PointerSlice(t *testing.T) {
	optVal := "Pointer Item"
	data := []*DummyStruct{
		{
			ID:        10,
			Name:      "Item Pointer",
			Optional:  &optVal,
			ZeroValue: 100,
		},
	}

	var buf bytes.Buffer
	err := ExportXLSX(&buf, data)
	if err != nil {
		t.Fatalf("ExportXLSX failed for pointer slice: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("Failed to parse XLSX: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		t.Fatalf("Failed to get rows: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rows))
	}
}

func TestExportXLSX_EmptySlice(t *testing.T) {
	var data []DummyStruct
	var buf bytes.Buffer

	err := ExportXLSX(&buf, data)
	if err != nil {
		t.Fatalf("ExportXLSX failed for empty slice: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("Failed to parse XLSX: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		t.Fatalf("Failed to get rows: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("Expected 1 header row, got %d", len(rows))
	}
}

func TestExportXLSX_InvalidInput(t *testing.T) {
	var buf bytes.Buffer

	err := ExportXLSX(&buf, "not a slice")
	if err == nil || err.Error() != "data must be a slice" {
		t.Errorf("Expected 'data must be a slice' error, got %v", err)
	}

	err = ExportXLSX(&buf, []int{1, 2, 3})
	if err == nil || err.Error() != "data must be a slice of structs" {
		t.Errorf("Expected 'data must be a slice of structs' error, got %v", err)
	}
}

func TestExportMultiSheetXLSX(t *testing.T) {
	sheets := []SheetData{
		{
			Name: "Brands",
			Data: []DummyStruct{{ID: 1, Name: "Brand One"}},
		},
		{
			Name: "Items",
			Data: []DummyStruct{{ID: 2, Name: "Item Two"}},
		},
	}

	var buf bytes.Buffer
	err := ExportMultiSheetXLSX(&buf, sheets)
	if err != nil {
		t.Fatalf("ExportMultiSheetXLSX failed: %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("Failed to parse output as excel file: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) != 2 || sheetList[0] != "Brands" || sheetList[1] != "Items" {
		t.Errorf("Unexpected sheet list: %v", sheetList)
	}

	rowsBrands, err := f.GetRows("Brands")
	if err != nil || len(rowsBrands) != 2 || rowsBrands[1][1] != "Brand One" {
		t.Errorf("Unexpected rows in Brands sheet: %v, err: %v", rowsBrands, err)
	}

	rowsItems, err := f.GetRows("Items")
	if err != nil || len(rowsItems) != 2 || rowsItems[1][1] != "Item Two" {
		t.Errorf("Unexpected rows in Items sheet: %v, err: %v", rowsItems, err)
	}
}

