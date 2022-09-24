
-- +migrate Up
ALTER TABLE `events` ADD COLUMN (
    `category` varchar(100) DEFAULT ''
);

-- +migrate Down
ALTER TABLE `events` DROP COLUMN `category`;