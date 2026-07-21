package models

import "time"

// Image merepresentasikan baris table_image. LocationID dan ItemID sengaja
// berupa pointer (bukan uint64 biasa) karena kolomnya nullable di DB dan
// CHECK constraint table_image mewajibkan tepat satu di antara keduanya
// terisi, sisanya NULL — pola yang sama seperti Item.BrandID/Item.LocationID.
type Image struct {
	ID         uint64    `db:"id" json:"id"`
	LocationID *uint64   `db:"location_id" json:"location_id,omitempty"`
	ItemID     *uint64   `db:"item_id" json:"item_id,omitempty"`
	ImagePath  string    `db:"image_path" json:"image_path"`
	IsPrimary  bool      `db:"is_primary" json:"is_primary"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
