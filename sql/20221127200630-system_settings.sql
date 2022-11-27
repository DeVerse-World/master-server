-- +migrate Up
CREATE TABLE `system_settings` (
                          `id` int unsigned NOT NULL AUTO_INCREMENT,
                          `key_name` varchar(100) DEFAULT '',
                          `value` varchar(100) DEFAULT '',
                          `category` varchar(100) DEFAULT '',
                          `object_reference_id` varchar(100) DEFAULT '',
                          `updated_at` timestamp NULL DEFAULT NULL,
                          `created_at` timestamp NULL DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `setting_identifier` (`key_name`, `category`, `object_reference_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `system_settings`;
