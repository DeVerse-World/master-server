-- +migrate Up
ALTER TABLE events ADD user_id int UNSIGNED DEFAULT NULL;
ALTER TABLE events ADD CONSTRAINT fk_event_user FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE;
-- +migrate Down
