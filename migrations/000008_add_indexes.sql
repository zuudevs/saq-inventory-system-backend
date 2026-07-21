-- +goose Up

CREATE INDEX `idx_image_location_id`
    ON `table_image`(`location_id`);

CREATE INDEX `idx_image_item_id`
    ON `table_image`(`item_id`);

CREATE UNIQUE INDEX `idx_image_location_primary`
	ON `table_image`(`location_id`)
	WHERE `is_primary` = 1;

CREATE UNIQUE INDEX `idx_image_item_primary`
	ON `table_image`(`item_id`)
	WHERE `is_primary` = 1;

-- +goose Down

DROP INDEX IF EXISTS `idx_image_location_id`;
DROP INDEX IF EXISTS `idx_image_item_id`;
DROP INDEX IF EXISTS `idx_image_location_primary`;
DROP INDEX IF EXISTS `idx_image_item_primary`;