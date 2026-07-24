package exporters

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// ExportCSV writes a slice of structs to a writer in CSV format.
// It uses reflection to inspect the "export" tags on the struct fields.
func ExportCSV(writer io.Writer, data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return errors.New("data must be a slice")
	}

	if val.Len() == 0 {
		return nil
	}

	elemType := val.Index(0).Type()
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

	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		itemVal := val.Index(i)
		if itemVal.Kind() == reflect.Ptr {
			itemVal = itemVal.Elem()
		}

		var row []string
		for _, fieldIdx := range fieldIndices {
			fieldVal := itemVal.Field(fieldIdx)
			row = append(row, formatValue(fieldVal))
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}

func formatValue(val reflect.Value) string {
	if !val.IsValid() {
		return ""
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		if t, ok := val.Interface().(interface{ Format(string) string }); ok {
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", val.Interface())
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}