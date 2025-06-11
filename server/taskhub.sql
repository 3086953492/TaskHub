-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `taskhub_user` DEFAULT CHARACTER
SET
    utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 使用数据库
USE `taskhub_user`;

-- 创建 users 表 
CREATE TABLE
    IF NOT EXISTS `users` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `username` VARCHAR(50) NOT NULL,
        `email` VARCHAR(100) NOT NULL, -- 移入核心表
        `password` VARCHAR(255) NOT NULL,
        `status` INT DEFAULT '1' COMMENT '1:正常 0:禁用',
        `role` VARCHAR(20) DEFAULT 'user' COMMENT 'admin, user',
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_username` (`username`),
        UNIQUE KEY `idx_email` (`email`), -- 邮箱唯一索引
        KEY `idx_status_deleted` (`status`, `deleted_at`) -- 组合索引优化查询
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表';

-- 创建用户信息表
CREATE TABLE
    IF NOT EXISTS `user_profiles` ( -- 统一复数命名
        `user_id` BIGINT UNSIGNED NOT NULL,
        `nickname` VARCHAR(50) DEFAULT NULL,
        `avatar` VARCHAR(255) DEFAULT NULL,
        PRIMARY KEY (`user_id`),
        CONSTRAINT `fk_user_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE -- 级联删除防止脏数据
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户信息表';