-- +migrate Up
CREATE TABLE `subworld_instances` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `host_name` varchar(100) DEFAULT '',
    `region` varchar(100) DEFAULT '',
    `max_players` int DEFAULT 0,
    `num_current_players` int DEFAULT 0,
    `instance_port` varchar(100) DEFAULT '',
    `beacon_port` varchar(100) DEFAULT '',
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `subworld_template_id` int unsigned DEFAULT NULL,
     CONSTRAINT `fk_template_id` FOREIGN KEY (`subworld_template_id`) REFERENCES `subworld_templates` (`id`),
    `host_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_instance_host` FOREIGN KEY (`host_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `subworld_instances`;
