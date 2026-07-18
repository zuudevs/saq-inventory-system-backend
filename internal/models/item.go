package models

import (
	"database/sql"
	"time"
)

type ItemCondition string
type ItemStatus string

const (
	ItemConditionGood        ItemCondition = `good`
	ItemConditionMinorDamage ItemCondition = `minor_damage`
	ItemConditionMajorDamage ItemCondition = `major_damage`
	ItemConditionLost        ItemCondition = `lost`
	ItemStatusActive         ItemStatus    = `active`
	ItemStatusInactive       ItemStatus    = `inactive`
	ItemStatusMaintenance    ItemStatus    = `maintenance`
	ItemStatusBorrowed       ItemStatus    = `borrowed`
)

type Item struct {
	ID            uint64         `db:"id" json:"id"`
	BrandID       *uint64        `db:"brand_id" json:"brand_id"`
	CategoryID    uint64         `db:"category_id" json:"category_id"`
	LocationID    *uint64        `db:"location_id" json:"location_id"`
	AssetCode     string         `db:"asset_code" json:"asset_code"`
	Name          string         `db:"name" json:"name"`
	Slug          string         `db:"slug" json:"slug"`
	ItemCondition ItemCondition  `db:"item_condition" json:"item_condition"`
	ItemStatus    ItemStatus     `db:"item_status" json:"item_status"`
	Notes         sql.NullString `db:"notes" json:"notes"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}
