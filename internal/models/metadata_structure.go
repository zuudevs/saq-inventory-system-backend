package models

import (
	"encoding/json"
	"time"
)

type MetadataFieldType string

const (
	MetadataFieldTypeString   MetadataFieldType = `string`
	MetadataFieldTypeText     MetadataFieldType = `text`
	MetadataFieldTypeInt      MetadataFieldType = `int`
	MetadataFieldTypeFloat    MetadataFieldType = `float`
	MetadataFieldTypeBool     MetadataFieldType = `bool`
	MetadataFieldTypeDate     MetadataFieldType = `date`
	MetadataFieldTypeDatetime MetadataFieldType = `datetime`
	MetadataFieldTypeEnum     MetadataFieldType = `enum`
)

type MetadataField struct {
	Name      string            `json:"name"`
	Label     string            `json:"label"`
	Type      MetadataFieldType `json:"type"`
	Length    *int              `json:"length,omitempty"`
	Precision *int              `json:"precision,omitempty"`
	Scale     *int              `json:"scale,omitempty"`
	Options   []string          `json:"options,omitempty"`
	Nullable  bool              `json:"nullable"`
	Default   *string           `json:"default,omitempty"`
	Unique    bool              `json:"unique"`
}

type MetadataStructure struct {
	ID         uint64          `db:"id" json:"id"`
	CategoryID uint64          `db:"category_id" json:"category_id"`
	Fields     json.RawMessage `db:"fields" json:"fields"`
	Version    uint            `db:"version" json:"version"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

func (m *MetadataStructure) DecodeFields() ([]MetadataField, error) {
	var fields []MetadataField

	if err := json.Unmarshal(m.Fields, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func EncodeMetadataFields(fields []MetadataField) (json.RawMessage, error) {
	raw, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(raw), nil
}
