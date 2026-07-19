package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const (
	defaultStringLength     = 255
	maxStringLength         = 2000
	defaultDecimalPrecision = 10
	defaultDecimalScale     = 2
	maxEnumOptionLength     = 100
)

var dateLiteralPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}([ T]\d{2}:\d{2}(:\d{2})?)?$`)
var dateOnlyPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// columnType menghasilkan fragmen tipe kolom MySQL (mis. "VARCHAR(255)")
// dari definisi field abstrak. Semua angka yang dipakai (length, precision,
// scale) sudah melalui validasi range di sini, tidak pernah disisipkan mentah.
func columnType(field models.MetadataField) (string, error) {
	switch field.Type {
	case models.MetadataFieldTypeString:
		length := defaultStringLength
		if field.Length != nil {
			length = *field.Length
		}
		if length < 1 || length > maxStringLength {
			return "", fmt.Errorf("length untuk field '%s' harus antara 1-%d", field.Name, maxStringLength)
		}
		return fmt.Sprintf("VARCHAR(%d)", length), nil

	case models.MetadataFieldTypeText:
		return "TEXT", nil

	case models.MetadataFieldTypeInt:
		return "BIGINT", nil

	case models.MetadataFieldTypeFloat:
		precision := defaultDecimalPrecision
		scale := defaultDecimalScale
		if field.Precision != nil {
			precision = *field.Precision
		}
		if field.Scale != nil {
			scale = *field.Scale
		}
		if precision < 1 || precision > 65 {
			return "", fmt.Errorf("precision untuk field '%s' harus antara 1-65", field.Name)
		}
		if scale < 0 || scale > precision {
			return "", fmt.Errorf("scale untuk field '%s' tidak valid", field.Name)
		}
		return fmt.Sprintf("DECIMAL(%d,%d)", precision, scale), nil

	case models.MetadataFieldTypeBool:
		return "TINYINT(1)", nil

	case models.MetadataFieldTypeDate:
		return "DATE", nil

	case models.MetadataFieldTypeDatetime:
		return "DATETIME", nil

	case models.MetadataFieldTypeEnum:
		if len(field.Options) == 0 {
			return "", fmt.Errorf("field enum '%s' wajib memiliki minimal satu option", field.Name)
		}

		quoted := make([]string, len(field.Options))
		for i, opt := range field.Options {
			if opt == "" || len(opt) > maxEnumOptionLength {
				return "", fmt.Errorf("option enum pada field '%s' tidak valid", field.Name)
			}
			quoted[i] = "'" + escapeStringLiteral(opt) + "'"
		}

		return fmt.Sprintf("ENUM(%s)", strings.Join(quoted, ",")), nil

	default:
		return "", fmt.Errorf("tipe field '%s' tidak dikenali", field.Type)
	}
}

// escapeStringLiteral meng-escape backslash dan single quote sesuai aturan
// MySQL string literal. Dipakai hanya untuk nilai yang akan disisipkan
// sebagai literal string di dalam DDL (default value, opsi enum) — bukan
// pengganti parameterized query untuk DML.
func escapeStringLiteral(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `'`, `\'`)
	return value
}

// defaultClause memvalidasi nilai default sesuai tipe field, lalu
// menghasilkan fragmen "DEFAULT ..." yang aman untuk DDL.
func defaultClause(field models.MetadataField) (string, error) {
	if field.Default == nil {
		return "", nil
	}

	value := *field.Default

	switch field.Type {
	case models.MetadataFieldTypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return "", fmt.Errorf("default field '%s' bukan integer valid", field.Name)
		}
		return "DEFAULT " + value, nil

	case models.MetadataFieldTypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return "", fmt.Errorf("default field '%s' bukan angka desimal valid", field.Name)
		}
		return "DEFAULT " + value, nil

	case models.MetadataFieldTypeBool:
		if value != "0" && value != "1" {
			return "", fmt.Errorf("default field '%s' harus '0' atau '1'", field.Name)
		}
		return "DEFAULT " + value, nil

	case models.MetadataFieldTypeEnum:
		valid := false
		for _, opt := range field.Options {
			if opt == value {
				valid = true
				break
			}
		}
		if !valid {
			return "", fmt.Errorf("default field '%s' harus salah satu dari options", field.Name)
		}
		return "DEFAULT '" + escapeStringLiteral(value) + "'", nil

	case models.MetadataFieldTypeString, models.MetadataFieldTypeText:
		return "DEFAULT '" + escapeStringLiteral(value) + "'", nil

	case models.MetadataFieldTypeDate, models.MetadataFieldTypeDatetime:
		// Nilai default tanggal dibatasi hanya digit, dash, colon, dan spasi
		// (format ISO 8601), atau literal khusus MySQL CURRENT_TIMESTAMP.
		if strings.EqualFold(value, "CURRENT_TIMESTAMP") {
			return "DEFAULT CURRENT_TIMESTAMP", nil
		}
		if !dateLiteralPattern.MatchString(value) {
			return "", fmt.Errorf("default field '%s' bukan format tanggal yang valid", field.Name)
		}
		return "DEFAULT '" + value + "'", nil

	default:
		return "", fmt.Errorf("tipe field '%s' tidak dikenali", field.Type)
	}
}
