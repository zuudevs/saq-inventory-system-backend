-- +goose Up

-- +goose StatementBegin

CREATE TRIGGER `trg_table_category_updated_at`
AFTER UPDATE ON `table_category`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_category`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TRIGGER `trg_table_brand_updated_at`
AFTER UPDATE ON `table_brand`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_brand`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TRIGGER `trg_table_location_updated_at`
AFTER UPDATE ON `table_location`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_location`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TRIGGER `trg_table_item_updated_at`
AFTER UPDATE ON `table_item`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_item`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TRIGGER `trg_table_image_updated_at`
AFTER UPDATE ON `table_image`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_image`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;

-- +goose StatementEnd

-- +goose StatementBegin

CREATE TRIGGER `trg_table_metadata_structure_updated_at`
AFTER UPDATE ON `table_metadata_structure`
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE `table_metadata_structure`
    SET `updated_at` = CURRENT_TIMESTAMP
    WHERE `id` = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down

DROP TRIGGER IF EXISTS `trg_table_category_updated_at`;
DROP TRIGGER IF EXISTS `trg_table_brand_updated_at`;
DROP TRIGGER IF EXISTS `trg_table_location_updated_at`;
DROP TRIGGER IF EXISTS `trg_table_item_updated_at`;
DROP TRIGGER IF EXISTS `trg_table_image_updated_at`;
DROP TRIGGER IF EXISTS `trg_table_metadata_structure_updated_at`;