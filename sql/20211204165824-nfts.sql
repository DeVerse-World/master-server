-- +migrate Up
CREATE TABLE `nfts` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `image_url` varchar(100) DEFAULT '',
    `name` varchar(100) DEFAULT '',
    `description` varchar(100) DEFAULT '',
    `require_fetch` boolean DEFAULT false,
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `wallet_id` int unsigned,
    CONSTRAINT `fk_nfts_wallets` FOREIGN KEY (`wallet_id`) REFERENCES `wallets` (`id`),
    `collection_id` int unsigned,
    CONSTRAINT `fk_nfts_collections` FOREIGN KEY (`collection_id`) REFERENCES `nft_collections` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `nfts`;
