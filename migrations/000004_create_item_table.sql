-- +goose Up
CREATE TABLE IF NOT EXISTS `table_item` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `brand_id` BIGINT UNSIGNED NULL,
    `category_id` BIGINT UNSIGNED NOT NULL,
    `location_id` BIGINT UNSIGNED NULL,

    `asset_code` VARCHAR(64) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(64) UNIQUE NOT NULL,

    `item_condition` ENUM('good', 'minor_damage', 'major_damage', 'lost') NOT NULL DEFAULT 'good',
    `item_status` ENUM('active', 'inactive', 'maintenance', 'borrowed') NOT NULL DEFAULT 'active',

    `notes` TEXT NULL,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT `fk_item_brand`
        FOREIGN KEY (`brand_id`) REFERENCES `table_brand`(`id`)
        ON UPDATE CASCADE
        ON DELETE SET NULL,

    CONSTRAINT `fk_item_category`
        FOREIGN KEY (`category_id`) REFERENCES `table_category`(`id`)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT `fk_item_location`
        FOREIGN KEY (`location_id`) REFERENCES `table_location`(`id`)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS `table_item`;