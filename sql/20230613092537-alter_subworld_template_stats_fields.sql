
-- +migrate Up
ALTER TABLE `subworld_templates` ADD COLUMN (
    `num_plays` int default 0,
    `num_views` int default 0,
    `num_clicks` int default 0
);


-- +migrate Down
ALTER TABLE `subworld_templates` DROP COLUMN (
    `num_plays`, `num_views`, `num_clicks`
);
