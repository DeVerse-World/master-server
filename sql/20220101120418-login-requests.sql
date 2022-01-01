-- +migrate Up
CREATE TABLE `login_requests` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `session_key` VARCHAR(100),
    UNIQUE KEY `session_key` (`session_key`),
    `updated_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    `wallet_id` int unsigned DEFAULT NULL,
    CONSTRAINT `fk_login_requests_wallets` FOREIGN KEY (`wallet_id`) REFERENCES `wallets` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `login_requests`;
