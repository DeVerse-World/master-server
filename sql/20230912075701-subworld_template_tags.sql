
-- +migrate Up
CREATE TABLE `subworld_template_tags` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `tag_name` varchar(300) DEFAULT '',
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `subworld_template_id` int unsigned NOT NULL,
    CONSTRAINT `fk_subworld_template_entity` FOREIGN KEY (`subworld_template_id`) REFERENCES `subworld_templates` (`id`) ON DELETE CASCADE,
    UNIQUE KEY `subworld_template_tags_unique` (`subworld_template_id`, `tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `subworld_template_tags`;
