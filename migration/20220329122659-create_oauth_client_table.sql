-- +migrate Up
CREATE TABLE `oauth_client` (
    `id` VARCHAR(255) NOT NULL,
    `sys_account_id` int(11) NOT NULL COMMENT 'ref:sys_account.id',
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `secret` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `domain` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `scope` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `icon_path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `data` TEXT COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate Down
DROP TABLE `oauth_client`;
