
-- +migrate Up
ALTER TABLE `avatars` MODIFY COLUMN `preprocess_url` varchar(500) DEFAULT '';
ALTER TABLE `events` MODIFY COLUMN `event_config_uri` varchar(500) DEFAULT '';
ALTER TABLE `events` MODIFY COLUMN `description` varchar(500) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `level_ipfs_uri` varchar(500) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `level_centralized_uri` varchar(500) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `thumbnail_centralized_uri` varchar(500) DEFAULT '';

-- +migrate Down
ALTER TABLE `avatars` MODIFY COLUMN `preprocess_url` varchar(100) DEFAULT '';
ALTER TABLE `events` MODIFY COLUMN `event_config_uri` varchar(100);
ALTER TABLE `events` MODIFY COLUMN `description` varchar(100) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `level_ipfs_uri` varchar(100) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `level_centralized_uri` varchar(100) DEFAULT '';
ALTER TABLE `subworld_templates` MODIFY COLUMN `thumbnail_centralized_uri` varchar(100) DEFAULT '';
