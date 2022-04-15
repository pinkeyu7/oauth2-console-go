-- +migrate Up
CREATE TABLE `oauth_scope` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `scope` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `path` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `method` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0:啟用 1:禁用',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate Down
DROP TABLE `oauth_scope`;
