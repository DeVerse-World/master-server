
-- +migrate Up
ALTER TABLE `action_reward_rules` MODIFY COLUMN `amount` FLOAT default 0;

-- +migrate Down
ALTER TABLE `action_reward_rules` MODIFY COLUMN `amount` int unsigned DEFAULT 0;