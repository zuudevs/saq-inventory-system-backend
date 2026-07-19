package schema

import (
	"fmt"
	"strings"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

const MaxFieldsPerCategory = 50

// MetadataTableName menghasilkan nama tabel metadata dari slug kategori,
// mis. "elektronik" -> "table_elektronik_metadata". Slug kategori (dibuat
// oleh gosimple/slug) memakai '-' sebagai separator, sedangkan identifier
// SQL kita hanya mengizinkan '_', jadi dinormalisasi dulu di sini —
// satu-satunya tempat konversi ini terjadi agar konsisten di seluruh layer.
func MetadataTableName(categorySlug string) string {
	normalized := strings.ReplaceAll(categorySlug, "-", "_")
	return "table_" + normalized + "_metadata"
}

// BuildCreateTableSQL menyusun DDL CREATE TABLE untuk tabel metadata sebuah
// kategori. Setiap identifier (nama tabel & nama kolom) divalidasi lewat
// ValidateIdentifier sebelum disisipkan ke string SQL; setiap default value
// divalidasi & di-escape lewat defaultClause. Fungsi ini murni (tidak
// menyentuh koneksi DB) sehingga mudah di-unit test.
func BuildCreateTableSQL(categorySlug string, fields []models.MetadataField) (string, error) {
	if len(fields) == 0 {
		return "", fmt.Errorf("minimal harus ada satu field metadata")
	}

	if len(fields) > MaxFieldsPerCategory {
		return "", fmt.Errorf("jumlah field metadata melebihi batas maksimum (%d)", MaxFieldsPerCategory)
	}

	tableName := MetadataTableName(categorySlug)
	if err := ValidateTableName(tableName); err != nil {
		return "", err
	}

	seenNames := make(map[string]struct{}, len(fields))
	columnDefs := make([]string, 0, len(fields))

	for _, field := range fields {
		if err := ValidateIdentifier(field.Name); err != nil {
			return "", err
		}

		if _, dup := seenNames[field.Name]; dup {
			return "", fmt.Errorf("nama field '%s' duplikat", field.Name)
		}
		seenNames[field.Name] = struct{}{}

		colType, err := columnType(field)
		if err != nil {
			return "", err
		}

		nullability := "NOT NULL"
		if field.Nullable {
			nullability = "NULL"
		}

		def, err := defaultClause(field)
		if err != nil {
			return "", err
		}

		parts := []string{QuoteIdentifier(field.Name), colType, nullability}
		if def != "" {
			parts = append(parts, def)
		}
		if field.Unique {
			parts = append(parts, "UNIQUE")
		}

		columnDefs = append(columnDefs, strings.Join(parts, " "))
	}

	var b strings.Builder

	fmt.Fprintf(&b, "CREATE TABLE %s (\n", QuoteIdentifier(tableName))
	b.WriteString("    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,\n")
	b.WriteString("    `item_id` BIGINT UNSIGNED NOT NULL,\n")

	for _, col := range columnDefs {
		fmt.Fprintf(&b, "    %s,\n", col)
	}

	b.WriteString("    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n")
	b.WriteString("    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n")

	normalizedSlug := strings.ReplaceAll(categorySlug, "-", "_")
	fkName := "fk_" + normalizedSlug + "_metadata_item"
	if len(fkName) > MaxIdentifierLength {
		fkName = fkName[:MaxIdentifierLength]
	}
	fmt.Fprintf(
		&b,
		"    CONSTRAINT %s FOREIGN KEY (`item_id`) REFERENCES `table_item`(`id`) ON DELETE CASCADE\n",
		QuoteIdentifier(fkName),
	)
	b.WriteString(")")

	return b.String(), nil
}

// BuildDropTableSQL menyusun DDL DROP TABLE untuk kompensasi ketika
// penyimpanan definisi field ke table_metadata_structure gagal setelah
// tabel metadata sudah terlanjur dibuat (lihat catatan transaction & DDL
// implicit commit pada Service).
func BuildDropTableSQL(categorySlug string) (string, error) {
	tableName := MetadataTableName(categorySlug)
	if err := ValidateTableName(tableName); err != nil {
		return "", err
	}

	return fmt.Sprintf("DROP TABLE IF EXISTS %s", QuoteIdentifier(tableName)), nil
}
