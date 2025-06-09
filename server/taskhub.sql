-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `taskhub_user` DEFAULT CHARACTER
SET
    utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 使用数据库
USE `taskhub_user`;

-- 创建 users 表
CREATE TABLE
    IF NOT EXISTS `users` (
        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
        `username` varchar(50) NOT NULL,
        `email` varchar(100) NOT NULL,
        `password` varchar(255) NOT NULL,
        `nickname` varchar(50) DEFAULT NULL,
        `avatar` varchar(255) DEFAULT NULL,
        `status` int DEFAULT '1' COMMENT '1:正常 0:禁用',
        `role` varchar(20) DEFAULT 'user' COMMENT 'admin, user',
        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_username` (`username`),
        UNIQUE KEY `idx_email` (`email`),
        KEY `idx_deleted_at` (`deleted_at`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;