-- +migrate Up
CREATE TABLE `avatars` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `preprocess_url` varchar(100) DEFAULT '',
     `postprocess_url` varchar(100) DEFAULT '',
     `updated_at` timestamp NULL DEFAULT NULL,
     `created_at` timestamp NULL DEFAULT NULL,
     PRIMARY KEY (`id`),
     `user_id` int unsigned DEFAULT NULL,
     CONSTRAINT `fk_user_avatar` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `avatars`;
