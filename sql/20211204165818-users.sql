
-- +migrate Up
CREATE TABLE `users` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `wallet_address` varchar(100) DEFAULT '',
   `wallet_nonce` varchar(100) DEFAULT '',
   `social_email` varchar(100) DEFAULT '',
   `custom_email` varchar(100) DEFAULT '',
   `custom_password_hash` varchar(100) DEFAULT '',
   `custom_salt` varchar(100) DEFAULT '',
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`),
   UNIQUE KEY `user_identifier` (`wallet_address`, `social_email`, `custom_email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `users`;
