
-- +migrate Up
ALTER TABLE `users` ADD COLUMN (
    `name` varchar(100) DEFAULT ''
);

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `name`;