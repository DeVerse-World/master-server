-- +migrate Up
CREATE TABLE `nft_collections` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `token_address` VARCHAR(100),
   `amount` int,
   `block_number` int,
   `minted_block_num` int,
   `contract_type` VARCHAR(100),
   `token_url` VARCHAR(100),
   `token_id` VARCHAR(100),
   `name` VARCHAR(100),
   `symbol` VARCHAR(100),
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `nft_collections`;
