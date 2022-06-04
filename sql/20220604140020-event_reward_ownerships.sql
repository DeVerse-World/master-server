-- +migrate Up
CREATE TABLE `event_reward_ownerships` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `amount` int,
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `event_reward_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_rewards_ownerships_reward` FOREIGN KEY (`event_reward_id`) REFERENCES `event_rewards` (`id`),
    `participant_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_reward_ownerships_participant` FOREIGN KEY (`participant_id`) REFERENCES `event_participants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `event_reward_ownerships`;
