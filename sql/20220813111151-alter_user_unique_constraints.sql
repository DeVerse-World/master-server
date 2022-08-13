--
-- -- +migrate Up
-- ALTER TABLE `users` DROP INDEX `user_identifier`;
-- ALTER TABLE `users` ADD CONSTRAINT UNIQUE KEY `user_identifier_1` (`wallet_address`) WHERE TRIM(`wallet_address`) <> '';
-- ALTER TABLE `users` ADD CONSTRAINT UNIQUE KEY `user_identifier_2` (`social_email`) WHERE TRIM(`social_email`) <> '';;
-- -- ALTER TABLE `users` ADD CONSTRAINT UNIQUE KEY `user_identifier_3` (`custom_email`);
--
-- -- +migrate Down
-- ALTER TABLE `users` DROP INDEX `user_identifier_1`;
-- ALTER TABLE `users` DROP INDEX `user_identifier_2`;
-- -- ALTER TABLE `users` DROP INDEX `user_identifier_3`;
-- ALTER TABLE `users` ADD CONSTRAINT UNIQUE KEY `user_identifier` (`wallet_address`, `social_email`, `custom_email`);
