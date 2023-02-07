-- +migrate Up
CREATE TABLE `nft_galleries` (
      `id` int unsigned NOT NULL AUTO_INCREMENT,
      `display_name` varchar(1000) DEFAULT '',
      `owner_address` varchar(1000) DEFAULT '',
      `category` varchar(100) DEFAULT '',
      `collection_address` varchar(1000) DEFAULT '',
      `chain` varchar(100) DEFAULT '',
      `is_auto_fetch` int DEFAULT 1,
      `updated_at` timestamp NULL DEFAULT NULL,
      `created_at` timestamp NULL DEFAULT NULL,
      PRIMARY KEY (`id`),
      `subworld_template_id` int unsigned DEFAULT NULL,
      CONSTRAINT `fk_gallery_template_id` FOREIGN KEY (`subworld_template_id`) REFERENCES `subworld_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE `nft_galleries`;
