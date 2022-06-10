-- +migrate Up
CREATE TABLE `event_rewards` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `supply` int,
    `min_eligible_rank` int,
    `max_eligible_rank` int,
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `minted_nft_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_rewards_minted_nft` FOREIGN KEY (`minted_nft_id`) REFERENCES `minted_nfts` (`id`),
    `event_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_event_rewards_event` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `event_rewards`;
