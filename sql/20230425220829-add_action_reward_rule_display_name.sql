
-- +migrate Up
ALTER TABLE `action_reward_rules` ADD COLUMN (
   `display_name` varchar(100) DEFAULT ''
);

-- +migrate Down
ALTER TABLE `action_reward_rules` DROP COLUMN `display_name`;