-- +migrate Up
CREATE TABLE `events` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) DEFAULT '',
    `description` varchar(100) DEFAULT '',
    `event_config_uri` varchar(100),
    `max_num_participants` int,
    `allow_temporary_hold` int,
    `stage` varchar(100),
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `events`;
