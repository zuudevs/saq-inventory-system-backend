-- +goose Up
CREATE TABLE IF NOT EXISTS `table_category` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,

    `name` TEXT UNIQUE NOT NULL,
    `slug` TEXT UNIQUE NOT NULL,
    `description` TEXT NULL,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS `table_category`;