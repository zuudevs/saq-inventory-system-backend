-- +goose Up
CREATE TABLE IF NOT EXISTS `table_metadata_structure` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `category_id` BIGINT UNSIGNED NOT NULL UNIQUE,

    `fields` JSON NOT NULL,
    `version` INT UNSIGNED NOT NULL DEFAULT 1,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT `fk_metadata_structure_category`
        FOREIGN KEY (`category_id`) REFERENCES `table_category`(`id`)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS `table_metadata_structure`;
