package schema

import (
	"testing"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

func TestBuildAddColumnSQL(t *testing.T) {
	slug := "test-category"

	t.Run("valid string field", func(t *testing.T) {
		length := 100
		field := models.MetadataField{
			Name:     "description",
			Type:     models.MetadataFieldTypeString,
			Length:   &length,
			Nullable: true,
		}

		sql, err := BuildAddColumnSQL(slug, field)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "ALTER TABLE `table_test_category_metadata` ADD COLUMN `description` VARCHAR(100) NULL"
		if sql != expected {
			t.Errorf("expected: %q, got: %q", expected, sql)
		}
	})

	t.Run("valid enum field", func(t *testing.T) {
		defaultValue := "red"
		field := models.MetadataField{
			Name:     "color",
			Type:     models.MetadataFieldTypeEnum,
			Options:  []string{"red", "green", "blue"},
			Nullable: false,
			Default:  &defaultValue,
		}

		sql, err := BuildAddColumnSQL(slug, field)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "ALTER TABLE `table_test_category_metadata` ADD COLUMN `color` TEXT NOT NULL DEFAULT 'red' CHECK (`color` IN ('red','green','blue'))"
		if sql != expected {
			t.Errorf("expected: %q, got: %q", expected, sql)
		}
	})
}

func TestBuildDropColumnSQL(t *testing.T) {
	slug := "test-category"

	t.Run("valid column name", func(t *testing.T) {
		sql, err := BuildDropColumnSQL(slug, "old_column")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "ALTER TABLE `table_test_category_metadata` DROP COLUMN `old_column`"
		if sql != expected {
			t.Errorf("expected: %q, got: %q", expected, sql)
		}
	})

	t.Run("invalid column name", func(t *testing.T) {
		_, err := BuildDropColumnSQL(slug, "invalid name;")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
