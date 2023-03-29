
-- +migrate Up
ALTER TABLE `subworld_templates` MODIFY COLUMN `derivative_uri` varchar(500) DEFAULT '';

-- +migrate Down
ALTER TABLE `subworld_templates` MODIFY COLUMN `derivative_uri` varchar(100) DEFAULT '';
