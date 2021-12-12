
-- +migrate Up
CREATE TABLE `wallets` (
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `address` varchar(100) DEFAULT '',
                           `nonce` varchar(100) DEFAULT '',
                           `updated_at` timestamp NULL DEFAULT NULL,
                           `created_at` timestamp NULL DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `address` (`address`),
                           `user_id` int unsigned DEFAULT NULL,
                           CONSTRAINT `fk_wallets_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `wallets`;
