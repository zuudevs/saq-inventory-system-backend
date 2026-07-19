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

// CreateMetadataTable membangun dan mengeksekusi CREATE TABLE untuk
// kategori tertentu. MySQL memicu implicit commit untuk setiap DDL,
// sehingga operasi ini TIDAK ikut serta dalam transaction SQL apa pun —
// pemanggil (services.MetadataStructureService) bertanggung jawab
// melakukan compensating action (DropMetadataTable) bila langkah
// berikutnya setelah ini gagal.
func (s *Service) CreateMetadataTable(categorySlug string, fields []models.MetadataField) error {
	ddl, err := BuildCreateTableSQL(categorySlug, fields)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ddl)
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
