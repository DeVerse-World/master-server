
-- +migrate Up
ALTER TABLE `subworld_templates` ADD COLUMN (
    `rating` int unsigned NOT NULL DEFAULT 5
);

-- +migrate Down
ALTER TABLE `subworld_templates` DROP COLUMN `rating`;