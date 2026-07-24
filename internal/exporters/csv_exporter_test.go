package exporters

import (
	"bytes"
	"strings"
	"testing"
)

type DummyStruct struct {
	ID        uint64  `export:"ID"`
	Name      string  `export:"Name"`
	IgnoreMe  string  `json:"ignore"`
	Optional  *string `export:"Optional Value"`
	ZeroValue int     `export:"Zero"`
}

func TestExportCSV(t *testing.T) {
	optVal := "Hello Opt"
	data := []DummyStruct{
		{
			ID:        1,
			Name:      "Item One",
			IgnoreMe:  "Secret",
			Optional:  &optVal,
			ZeroValue: 42,
		},
		{
			ID:        2,
			Name:      "Item Two",
			IgnoreMe:  "Secret 2",
			Optional:  nil,
			ZeroValue: 0,
		},
	}

	var buf bytes.Buffer
	err := ExportCSV(&buf, data)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	result := buf.String()
	expectedHeaders := "ID,Name,Optional Value,Zero"
	if !strings.Contains(result, expectedHeaders) {
		t.Errorf("Expected headers %q, got %q", expectedHeaders, result)
	}

	expectedRow1 := "1,Item One,Hello Opt,42"
	if !strings.Contains(result, expectedRow1) {
		t.Errorf("Expected row 1 %q, got %q", expectedRow1, result)
	}

	expectedRow2 := "2,Item Two,,0"
	if !strings.Contains(result, expectedRow2) {
		t.Errorf("Expected row 2 %q, got %q", expectedRow2, result)
	}
}
