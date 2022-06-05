-- +migrate Up
CREATE TABLE `minted_nfts` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `token_address` varchar(100) DEFAULT '',
    `token_id` varchar(100) DEFAULT '',
    `name` varchar(100) DEFAULT '',
    `description` varchar(100) DEFAULT '',
    `supply` int DEFAULT 1,
    `asset_type` varchar(100) NOT NULL CHECK (`asset_type` IN('2D Image', 'Character Race','Character Skin','New Gameplay mode','New Bot Logic')),
    `file_asset_name` varchar(100),
    `file_asset_uri` varchar(100),
    `file_asset_uri_from_centralized` varchar(100),
    `file_2d_uri` varchar(100),
    `file_3d_uri` varchar(100),
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `minted_nft_uniq` (`token_address`, `token_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `minted_nfts`;
