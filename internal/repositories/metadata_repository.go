package repositories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/zuudevs/saq-inventory-system-backend/internal/schema"
)

// MetadataRepository membaca/menulis tabel metadata dinamis
// (table_<slug>_metadata). Berbeda dengan repository lain yang punya nama
// tabel & kolom tetap, di sini nama tabel dan daftar kolom baru diketahui
// saat runtime — karena itu setiap identifier divalidasi ulang lewat
// package schema sebelum disisipkan ke SQL text, meskipun sumbernya sudah
// tervalidasi sebelumnya saat metadata structure dibuat (defense in depth).
type MetadataRepository struct {
	db *sqlx.DB
}

func NewMetadataRepository(db *sqlx.DB) *MetadataRepository {
	return &MetadataRepository{
		db: db,
	}
}

// InsertWithExecutor menyimpan satu baris metadata untuk sebuah item.
// exec bisa berupa *sqlx.DB atau *sqlx.Tx — dipakai ItemService supaya
// insert ini ikut dalam transaction yang sama dengan insert ke table_item.
// values HARUS sudah melalui schema.ValidateMetadataPayload sebelum sampai
// di sini; keys pada values dipakai sebagai nama kolom.
func (r *MetadataRepository) InsertWithExecutor(
	exec sqlExecutor,
	tableName string,
	itemID uint64,
	values map[string]any,
) error {
	if err := schema.ValidateTableName(tableName); err != nil {
		return err
	}

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys) // urutan kolom deterministik, memudahkan debugging & test

	columns := make([]string, 0, len(keys)+1)
	placeholders := make([]string, 0, len(keys)+1)
	args := make([]any, 0, len(keys)+1)

	columns = append(columns, schema.QuoteIdentifier("item_id"))
	placeholders = append(placeholders, "?")
	args = append(args, itemID)

	for _, k := range keys {
		if err := schema.ValidateIdentifier(k); err != nil {
			return err
		}

		columns = append(columns, schema.QuoteIdentifier(k))
		placeholders = append(placeholders, "?")
		args = append(args, values[k])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		schema.QuoteIdentifier(tableName),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := exec.Exec(query, args...)
	return err
}

// FindByItemID mengambil satu baris metadata milik sebuah item. Hasilnya
// map generik karena struktur kolom berbeda-beda per kategori dan tidak
// diketahui saat compile time. Mengembalikan nil (tanpa error) bila belum
// ada baris metadata untuk item tersebut.
func (r *MetadataRepository) FindByItemID(tableName string, itemID uint64) (map[string]any, error) {
	if err := schema.ValidateTableName(tableName); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE `item_id` = ? LIMIT 1",
		schema.QuoteIdentifier(tableName),
	)

	rows, err := r.db.Queryx(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	row := make(map[string]any)
	if err := rows.MapScan(row); err != nil {
		return nil, err
	}

	return row, nil
}
