package exporters

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

type DummyStruct struct {
	ID        uint64    `export:"ID"`
	Name      string    `export:"Name"`
	IgnoreMe  string    `json:"ignore"`
	Optional  *string   `export:"Optional Value"`
	ZeroValue int       `export:"Zero"`
	CreatedAt time.Time `export:"Created At"`
}

func TestExportCSV_StructSlice(t *testing.T) {
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
			CreatedAt: time.Time{}, // zero time
		},
	}

	var buf bytes.Buffer
	err := ExportCSV(&buf, data)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	result := buf.String()
	expectedHeaders := "ID,Name,Optional Value,Zero,Created At"
	if !strings.Contains(result, expectedHeaders) {
		t.Errorf("Expected headers %q, got %q", expectedHeaders, result)
	}

	expectedRow1 := "1,Item One,Hello Opt,42,2026-07-24 18:00:00"
	if !strings.Contains(result, expectedRow1) {
		t.Errorf("Expected row 1 %q, got %q", expectedRow1, result)
	}

	expectedRow2 := "2,Item Two,,0,"
	if !strings.Contains(result, expectedRow2) {
		t.Errorf("Expected row 2 %q, got %q", expectedRow2, result)
	}
}

func TestExportCSV_PointerSlice(t *testing.T) {
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
	err := ExportCSV(&buf, data)
	if err != nil {
		t.Fatalf("ExportCSV failed for pointer slice: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "10,Item Pointer,Pointer Item,100,") {
		t.Errorf("Unexpected result for pointer slice: %q", result)
	}
}

func TestExportCSV_EmptySlice(t *testing.T) {
	var data []DummyStruct
	var buf bytes.Buffer

	err := ExportCSV(&buf, data)
	if err != nil {
		t.Fatalf("ExportCSV failed for empty slice: %v", err)
	}

	result := strings.TrimSpace(buf.String())
	expectedHeaders := "ID,Name,Optional Value,Zero,Created At"
	if result != expectedHeaders {
		t.Errorf("Expected headers %q, got %q", expectedHeaders, result)
	}
}

func TestExportCSV_InvalidInput(t *testing.T) {
	var buf bytes.Buffer

	// Test non-slice
	err := ExportCSV(&buf, "not a slice")
	if err == nil || err.Error() != "data must be a slice" {
		t.Errorf("Expected 'data must be a slice' error, got %v", err)
	}

	// Test slice of non-struct
	err = ExportCSV(&buf, []int{1, 2, 3})
	if err == nil || err.Error() != "data must be a slice of structs" {
		t.Errorf("Expected 'data must be a slice of structs' error, got %v", err)
	}
}
