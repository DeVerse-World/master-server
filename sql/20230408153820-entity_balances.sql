-- +migrate Up
CREATE TABLE `entity_balances` (
       `id` int unsigned NOT NULL AUTO_INCREMENT,
       `entity_id` int unsigned DEFAULT NULL,
       `entity_type` varchar(50) DEFAULT '',
       `balance_amount` int unsigned DEFAULT 0,
       `balance_type` varchar(50) DEFAULT '',
       `updated_at` timestamp NULL DEFAULT NULL,
       `created_at` timestamp NULL DEFAULT NULL,
       PRIMARY KEY (`id`),
       UNIQUE KEY `entity_balance_identifier` (`entity_type`, `entity_id`, `balance_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `entity_balances`;
