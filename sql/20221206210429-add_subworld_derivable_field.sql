
-- +migrate Up
ALTER TABLE `subworld_templates` ADD COLUMN (
    `derivable` BOOLEAN DEFAULT 1
);

-- +migrate Down
ALTER TABLE `subworld_templates` DROP COLUMN `derivable`;