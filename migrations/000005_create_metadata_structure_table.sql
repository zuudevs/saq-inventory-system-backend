-- +goose Up
CREATE TABLE IF NOT EXISTS `table_metadata_structure` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `category_id` INTEGER NOT NULL UNIQUE,

    `fields` TEXT NOT NULL,
    `version` INTEGER NOT NULL DEFAULT 1,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT `fk_metadata_structure_category`
        FOREIGN KEY (`category_id`) REFERENCES `table_category`(`id`)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS `table_metadata_structure`;
