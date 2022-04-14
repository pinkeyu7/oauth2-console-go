-- +migrate Up
CREATE TABLE `sys_account` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `account` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0:啟用 1:禁用',
    `verify_at` datetime NULL,
    `forgot_pass_token` varchar(64) COLLATE utf8mb4_unicode_ci NULL,
    `forgot_pass_token_expired_at` datetime NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate Down
DROP TABLE `sys_account`;
