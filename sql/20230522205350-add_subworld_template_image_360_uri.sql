
-- +migrate Up
ALTER TABLE `subworld_templates` ADD COLUMN (
    `image_360_uri` varchar(500) DEFAULT ''
);

-- +migrate Down
ALTER TABLE `subworld_templates` DROP COLUMN `image_360_uri`;