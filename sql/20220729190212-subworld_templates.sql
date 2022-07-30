-- +migrate Up
CREATE TABLE `subworld_templates` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `file_name` varchar(100) DEFAULT '',
    `display_name` varchar(100) DEFAULT '',
    `level_ipfs_uri` varchar(100) DEFAULT '',
    `level_centralized_uri` varchar(100) DEFAULT '',
    `thumbnail_centralized_uri` varchar(100) DEFAULT '',
    `derivative_uri` varchar(100) DEFAULT '',
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `file_name` (`file_name`),
    `parent_subworld_template_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_parent_id` FOREIGN KEY (`parent_subworld_template_id`) REFERENCES `subworld_templates` (`id`),
    `creator_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_template_creator` FOREIGN KEY (`creator_id`) REFERENCES `wallets` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `subworld_templates`;
