-- +migrate Up
CREATE TABLE `event_participants` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `score` FLOAT,
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `event_participant_uniq` (`user_id`, `event_id`),
    `user_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_participants_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    `event_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_participants_event` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `event_participants`;
