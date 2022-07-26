
-- +migrate Up
ALTER TABLE `avatars` ADD COLUMN (
    `name` varchar(100) DEFAULT ''
);

-- +migrate Down
ALTER TABLE `avatars` DROP COLUMN `name`;