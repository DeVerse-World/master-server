
-- +migrate Up
ALTER TABLE `action_reward_records` RENAME COLUMN `amount` to `occur_count`;

-- +migrate Down
ALTER TABLE `action_reward_records` MODIFY COLUMN `occur_count` to `amount`;
