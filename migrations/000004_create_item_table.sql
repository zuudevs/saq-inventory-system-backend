-- +goose Up
CREATE TABLE IF NOT EXISTS `table_item` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,

	`brand_id` INTEGER NULL,
	`category_id` INTEGER NOT NULL,
	`location_id` INTEGER NULL,

	`asset_code` TEXT NOT NULL UNIQUE,
	`name` TEXT NOT NULL,
	`slug` TEXT NOT NULL UNIQUE,

	`item_condition` TEXT NOT NULL
        CHECK(`item_condition` IN ('good', 'minor_damage', 'major_damage', 'lost'))
        DEFAULT 'good',

    `item_status` TEXT NOT NULL
        CHECK(`item_status` IN ('active', 'inactive', 'maintenance', 'borrowed'))
        DEFAULT 'active',

    `notes` TEXT NULL,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT `fk_item_brand`
        FOREIGN KEY (`brand_id`)
        REFERENCES `table_brand`(`id`)
        ON UPDATE CASCADE
        ON DELETE SET NULL,

    CONSTRAINT `fk_item_category`
        FOREIGN KEY (`category_id`)
        REFERENCES `table_category`(`id`)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT `fk_item_location`
        FOREIGN KEY (`location_id`)
        REFERENCES `table_location`(`id`)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS table_item;