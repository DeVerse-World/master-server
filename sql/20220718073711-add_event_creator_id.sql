-- +migrate Up
ALTER TABLE events ADD wallet_id int UNSIGNED DEFAULT NULL;
ALTER TABLE events ADD CONSTRAINT fk_event_wallet FOREIGN KEY (`wallet_id`) REFERENCES `wallets`(`id`);
-- +migrate Down
