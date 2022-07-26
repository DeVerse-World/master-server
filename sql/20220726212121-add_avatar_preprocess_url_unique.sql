
-- +migrate Up
ALTER TABLE `avatars` ADD CONSTRAINT `avatar_preprocess_url_unique` UNIQUE (`preprocess_url`);

-- +migrate Down
ALTER TABLE `avatars` DROP CONSTRAINT `avatar_preprocess_url_unique`;