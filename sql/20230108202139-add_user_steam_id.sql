
-- +migrate Up
ALTER TABLE `users` ADD COLUMN (
    `steam_id` varchar(100) DEFAULT ''
);
ALTER TABLE `users` DROP INDEX `user_identifier`, ADD UNIQUE KEY `user_identifier` (`wallet_address`, `social_email`, `custom_email`, `steam_id`);

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `steam_id`;
ALTER TABLE `users` DROP INDEX `user_identifier`, ADD UNIQUE KEY `user_identifier` (`wallet_address`, `social_email`, `custom_email`);
