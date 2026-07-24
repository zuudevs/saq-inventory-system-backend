package exporters

import (
	"errors"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

type SheetData struct {
	Name string
	Data interface{}
}

// ExportMultiSheetXLSX writes multiple datasets into separate worksheets of a single XLSX workbook.
func ExportMultiSheetXLSX(writer io.Writer, sheets []SheetData) error {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	defaultSheet := "Sheet1"

	for idx, sheetData := range sheets {
		sheetName := sheetData.Name
		if sheetName == "" {
			sheetName = "Sheet"
		}

		if idx == 0 {
			if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
				return err
			}
		} else {
			if _, err := f.NewSheet(sheetName); err != nil {
				return err
			}
		}

		if err := writeSheetData(f, sheetName, sheetData.Data); err != nil {
			return err
		}
	}

	_, err := f.WriteTo(writer)
	return err
}

// ExportXLSX writes a single slice of structs to a writer in XLSX format on Sheet1.
func ExportXLSX(writer io.Writer, data interface{}) error {
	return ExportMultiSheetXLSX(writer, []SheetData{{Name: "Sheet1", Data: data}})
}

func writeSheetData(f *excelize.File, sheet string, data interface{}) error {
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

	return nil
}

