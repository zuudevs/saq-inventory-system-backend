package repositories

import "database/sql"

// sqlExecutor adalah subset method yang dimiliki bersama oleh *sqlx.DB dan
// *sqlx.Tx. Repository yang menerima interface ini bisa dijalankan baik
// langsung terhadap koneksi maupun di dalam transaction yang sama dengan
// repository lain — dipakai supaya INSERT ke table_item dan INSERT ke
// table_<slug>_metadata bisa benar-benar atomic (keduanya DML, beda dengan
// DDL yang implicit-commit).
type sqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest any, query string, args ...any) error
}
