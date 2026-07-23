package schema

import (
	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/models"
)

// Service bertanggung jawab menerjemahkan definisi field metadata menjadi
// DDL dan mengeksekusinya. Ia tidak tahu apa-apa soal business rules
// (validasi kategori, dsb.) — itu tanggung jawab layer Service di atasnya
// (services.MetadataStructureService).
type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

// TableName mengembalikan nama tabel metadata untuk sebuah slug kategori.
func (s *Service) TableName(categorySlug string) string {
	return MetadataTableName(categorySlug)
}

// CreateMetadataTable membangun dan mengeksekusi CREATE TABLE (plus trigger
// updated_at) untuk kategori tertentu. Operasi ini TIDAK ikut serta dalam
// transaction SQL apa pun — pemanggil (services.MetadataStructureService)
// bertanggung jawab melakukan compensating action (DropMetadataTable) bila
// langkah berikutnya setelah ini gagal.
func (s *Service) CreateMetadataTable(categorySlug string, fields []models.MetadataField) error {
	ddl, err := BuildCreateTableSQL(categorySlug, fields)
	if err != nil {
		return err
	}

	if _, err := s.db.Exec(ddl); err != nil {
		return err
	}

	triggerDDL, err := BuildUpdatedAtTriggerSQL(categorySlug)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(triggerDDL)
	return err
}

// DropMetadataTable adalah compensating action: menghapus tabel metadata
// yang sudah terlanjur dibuat ketika penyimpanan definisi field ke
// table_metadata_structure gagal setelahnya.
func (s *Service) DropMetadataTable(categorySlug string) error {
	ddl, err := BuildDropTableSQL(categorySlug)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ddl)
	return err
}

// UpdateMetadataTable membandingkan schema metadata lama dengan yang baru
// dan menerapkan perubahan kolom secara dinamis menggunakan ALTER TABLE.
func (s *Service) UpdateMetadataTable(categorySlug string, oldFields []models.MetadataField, newFields []models.MetadataField) error {
	oldMap := make(map[string]models.MetadataField)
	for _, f := range oldFields {
		oldMap[f.Name] = f
	}

	newMap := make(map[string]models.MetadataField)
	for _, f := range newFields {
		newMap[f.Name] = f
	}

	var sqls []string

	// 1. Identifikasi kolom yang perlu di-drop (dihapus atau dimodifikasi)
	for _, oldField := range oldFields {
		newField, exists := newMap[oldField.Name]
		if !exists {
			sql, err := BuildDropColumnSQL(categorySlug, oldField.Name)
			if err != nil {
				return err
			}
			sqls = append(sqls, sql)
		} else if isFieldModified(oldField, newField) {
			sql, err := BuildDropColumnSQL(categorySlug, oldField.Name)
			if err != nil {
				return err
			}
			sqls = append(sqls, sql)
		}
	}

	// 2. Identifikasi kolom yang perlu di-add (baru atau dimodifikasi)
	for _, newField := range newFields {
		oldField, exists := oldMap[newField.Name]
		if !exists {
			sql, err := BuildAddColumnSQL(categorySlug, newField)
			if err != nil {
				return err
			}
			sqls = append(sqls, sql)
		} else if isFieldModified(oldField, newField) {
			sql, err := BuildAddColumnSQL(categorySlug, newField)
			if err != nil {
				return err
			}
			sqls = append(sqls, sql)
		}
	}

	if len(sqls) == 0 {
		return nil
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range sqls {
		if _, err := tx.Exec(query); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func isFieldModified(f1, f2 models.MetadataField) bool {
	if f1.Type != f2.Type || f1.Nullable != f2.Nullable || f1.Unique != f2.Unique {
		return true
	}
	if (f1.Length == nil) != (f2.Length == nil) || (f1.Length != nil && *f1.Length != *f2.Length) {
		return true
	}
	if (f1.Precision == nil) != (f2.Precision == nil) || (f1.Precision != nil && *f1.Precision != *f2.Precision) {
		return true
	}
	if (f1.Scale == nil) != (f2.Scale == nil) || (f1.Scale != nil && *f1.Scale != *f2.Scale) {
		return true
	}
	if (f1.Default == nil) != (f2.Default == nil) || (f1.Default != nil && *f1.Default != *f2.Default) {
		return true
	}
	if len(f1.Options) != len(f2.Options) {
		return true
	}
	for i := range f1.Options {
		if f1.Options[i] != f2.Options[i] {
			return true
		}
	}
	return false
}