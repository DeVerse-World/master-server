-- +migrate Up
CREATE TABLE `action_reward_records` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `amount` int unsigned NOT NULL DEFAULT 0,
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`),
   `action_reward_rule_id` int unsigned DEFAULT NULL,
   CONSTRAINT `fk_action_reward_rule` FOREIGN KEY (`action_reward_rule_id`) REFERENCES `action_reward_rules` (`id`) ON DELETE CASCADE,
   `user_id` int unsigned DEFAULT NULL,
   CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `action_reward_records`;
