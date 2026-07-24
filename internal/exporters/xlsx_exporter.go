package exporters

import (
	"errors"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// ExportXLSX writes a slice of structs to a writer in XLSX format.
// It uses reflection to inspect the "export" tags on the struct fields.
func ExportXLSX(writer io.Writer, data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return errors.New("data must be a slice")
	}

	elemType := val.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return errors.New("data must be a slice of structs")
	}

	var headers []string
	var fieldIndices []int

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("export")
		if tag != "" && tag != "-" {
			headers = append(headers, tag)
			fieldIndices = append(fieldIndices, i)
		}
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheet := "Sheet1"

	// Write headers
	for colIdx, header := range headers {
		cell, err := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return err
		}
	}

	// Write rows
	for i := 0; i < val.Len(); i++ {
		itemVal := val.Index(i)
		if itemVal.Kind() == reflect.Ptr {
			if itemVal.IsNil() {
				continue
			}
			itemVal = itemVal.Elem()
		}

		rowNum := i + 2
		for colIdx, fieldIdx := range fieldIndices {
			fieldVal := itemVal.Field(fieldIdx)
			cell, err := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err != nil {
				return err
			}
			valStr := formatValue(fieldVal)
			if err := f.SetCellValue(sheet, cell, valStr); err != nil {
				return err
			}
		}
	}

	_, err := f.WriteTo(writer)
	return err
}
