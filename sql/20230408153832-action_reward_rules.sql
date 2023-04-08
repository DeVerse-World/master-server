-- +migrate Up
CREATE TABLE `action_reward_rules` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `action_name` varchar(100) DEFAULT '',
   `amount` int unsigned DEFAULT 0,
   `limit` int unsigned DEFAULT 0,
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`),
   `entity_balance_id` int unsigned DEFAULT NULL,
   CONSTRAINT `fk_action_reward_entity` FOREIGN KEY (`entity_balance_id`) REFERENCES `entity_balances` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `action_reward_rules`;
