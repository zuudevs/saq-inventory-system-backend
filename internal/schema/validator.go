package schema

import (
	"fmt"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

// ValidateMetadataPayload memvalidasi payload metadata (hasil form dinamis
// yang diisi user) terhadap definisi field dari table_metadata_structure,
// lalu mengembalikan map ternormalisasi (key = nama kolom, value = tipe Go
// yang siap dikirim sebagai parameter query).
//
// Hanya field yang benar-benar terdaftar di schema yang boleh muncul di
// payload maupun di hasil akhir — mencegah payload menyisipkan key di luar
// yang didefinisikan saat kategori dibuat. Field yang tidak dikirim (atau
// bernilai null) dan bersifat nullable atau punya default cukup dilewati
// saja di sini, biar MySQL yang menerapkan NULL/DEFAULT kolomnya.
func ValidateMetadataPayload(fields []models.MetadataField, payload map[string]any) (map[string]any, error) {
	allowed := make(map[string]models.MetadataField, len(fields))
	for _, field := range fields {
		allowed[field.Name] = field
	}

	for key := range payload {
		if _, ok := allowed[key]; !ok {
			return nil, fmt.Errorf("unknown metadata field: %s", key)
		}
	}

	result := make(map[string]any, len(fields))

	for _, field := range fields {
		raw, exists := payload[field.Name]

		if !exists || raw == nil {
			if field.Nullable || field.Default != nil {
				continue
			}
			return nil, fmt.Errorf("metadata field '%s' is required", field.Name)
		}

		value, err := coerceFieldValue(field, raw)
		if err != nil {
			return nil, err
		}

		result[field.Name] = value
	}

	return result, nil
}

// coerceFieldValue memvalidasi tipe nilai JSON yang dikirim user (hasil
// decode ke map[string]any: string, float64, bool, dst.) sesuai tipe field
// yang didefinisikan, lalu mengonversinya ke tipe Go yang tepat untuk
// dikirim ke database driver.
func coerceFieldValue(field models.MetadataField, raw any) (any, error) {
	switch field.Type {
	case models.MetadataFieldTypeString, models.MetadataFieldTypeText:
		s, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("metadata field '%s' must be a string", field.Name)
		}

		if field.Type == models.MetadataFieldTypeString {
			length := defaultStringLength
			if field.Length != nil {
				length = *field.Length
			}
			if len(s) > length {
				return nil, fmt.Errorf("metadata field '%s' exceeds max length %d", field.Name, length)
			}
		}

		return s, nil

	case models.MetadataFieldTypeInt:
		n, ok := raw.(float64)
		if !ok || n != float64(int64(n)) {
			return nil, fmt.Errorf("metadata field '%s' must be an integer", field.Name)
		}
		return int64(n), nil

	case models.MetadataFieldTypeFloat:
		n, ok := raw.(float64)
		if !ok {
			return nil, fmt.Errorf("metadata field '%s' must be a number", field.Name)
		}
		return n, nil

	case models.MetadataFieldTypeBool:
		b, ok := raw.(bool)
		if !ok {
			return nil, fmt.Errorf("metadata field '%s' must be a boolean", field.Name)
		}
		if b {
			return 1, nil
		}
		return 0, nil

	case models.MetadataFieldTypeDate:
		s, ok := raw.(string)
		if !ok || !dateOnlyPattern.MatchString(s) {
			return nil, fmt.Errorf("metadata field '%s' must be a valid date (YYYY-MM-DD)", field.Name)
		}
		return s, nil

	case models.MetadataFieldTypeDatetime:
		s, ok := raw.(string)
		if !ok || !dateLiteralPattern.MatchString(s) {
			return nil, fmt.Errorf("metadata field '%s' must be a valid datetime", field.Name)
		}
		return s, nil

	case models.MetadataFieldTypeEnum:
		s, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("metadata field '%s' must be a string", field.Name)
		}
		for _, opt := range field.Options {
			if opt == s {
				return s, nil
			}
		}
		return nil, fmt.Errorf("metadata field '%s' must be one of %v", field.Name, field.Options)

	default:
		return nil, fmt.Errorf("unsupported metadata field type for field: %s", field.Name)
	}
}
