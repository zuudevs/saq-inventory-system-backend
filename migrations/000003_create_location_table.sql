-- +goose Up
CREATE TABLE IF NOT EXISTS `table_location` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,

    `name` TEXT NOT NULL,
    `slug` TEXT UNIQUE NOT NULL,
    `room_code` TEXT NULL,
    `description` TEXT NULL,

    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS `table_location`;