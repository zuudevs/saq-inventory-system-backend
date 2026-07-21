-- +goose Up

CREATE TABLE IF NOT EXISTS `table_image` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,

    `location_id` INTEGER,
    `item_id` INTEGER,

    `image_path` TEXT NOT NULL,

    `is_primary` INTEGER NOT NULL DEFAULT 0
        CHECK (`is_primary` IN (0, 1)),

	`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        (`item_id` IS NOT NULL AND `location_id` IS NULL)
        OR
        (`item_id` IS NULL AND `location_id` IS NOT NULL)
    ),

    CONSTRAINT `fk_image_location`
        FOREIGN KEY (`location_id`)
        REFERENCES `table_location`(`id`)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT `fk_image_item`
        FOREIGN KEY (`item_id`)
        REFERENCES `table_item`(`id`)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS `table_image`;