package schema

import (
	"fmt"
	"regexp"
)

const MaxIdentifierLength = 64

var identifierPattern = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,62}$`)

var reservedIdentifiers = map[string]struct{}{
	`id`:         {},
	`item_id`:    {},
	`created_at`: {},
	`updated_at`: {},

	`select`: {}, `insert`: {}, `update`: {}, `delete`: {}, `drop`: {},
	`alter`: {}, `table`: {}, `database`: {}, `from`: {}, `where`: {},
	`join`: {}, `union`: {}, `into`: {}, `values`: {}, `set`: {},
	`grant`: {}, `revoke`: {}, `index`: {}, `primary`: {}, `foreign`: {},
	`key`: {}, `constraint`: {}, `null`: {}, `default`: {}, `order`: {},
	`group`: {}, `having`: {}, `limit`: {}, `and`: {}, `or`: {}, `not`: {},
}

// ValidateIdentifier memastikan sebuah nama (nama tabel atau nama kolom)
// aman untuk disisipkan langsung ke dalam DDL string. DDL statement tidak
// bisa diparameterisasi seperti DML, jadi validasi whitelist ini adalah
// lapisan pertahanan utama terhadap SQL injection lewat identifier.
func ValidateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("identifier tidak boleh kosong")
	}

	if len(name) > MaxIdentifierLength {
		return fmt.Errorf("identifier '%s' melebihi %d karakter", name, MaxIdentifierLength)
	}

	if !identifierPattern.MatchString(name) {
		return fmt.Errorf(
			"identifier '%s' tidak valid, hanya boleh huruf kecil, angka, dan underscore, diawali huruf/underscore",
			name,
		)
	}

	if _, reserved := reservedIdentifiers[name]; reserved {
		return fmt.Errorf("identifier '%s' adalah reserved word/nama kolom sistem", name)
	}

	return nil
}

// QuoteIdentifier membungkus identifier yang sudah tervalidasi dengan
// backtick. Jangan pernah memanggil ini terhadap identifier yang belum
// lolos ValidateIdentifier.
func QuoteIdentifier(name string) string {
	return "`" + name + "`"
}

// ValidateTableName memvalidasi nama tabel metadata final (termasuk
// prefix/suffix) tidak melebihi batas panjang identifier (64 karakter).
func ValidateTableName(tableName string) error {
	if len(tableName) > MaxIdentifierLength {
		return fmt.Errorf(
			"nama tabel '%s' melebihi %d karakter, perpendek slug kategori",
			tableName,
			MaxIdentifierLength,
		)
	}

	if !identifierPattern.MatchString(tableName) {
		return fmt.Errorf("nama tabel '%s' tidak valid", tableName)
	}

	return nil
}