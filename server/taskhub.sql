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


CREATE DATABASE IF NOT EXISTS `taskhub_task` DEFAULT CHARACTER
SET
    utf8mb4 COLLATE utf8mb4_0900_ai_ci;
    
USE `taskhub_task`;

-- 创建任务表
CREATE TABLE
    IF NOT EXISTS `tasks` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `status` INT DEFAULT '1' COMMENT '状态 1:待处理 2:进行中 3:已完成 4:已取消',
        `assignee_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '被分配人ID (关联taskhub_user.users.id)',
        `creator_id` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID (关联taskhub_user.users.id)',
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_assignee_status` (`assignee_id`, `status`),
        KEY `idx_creator_status` (`creator_id`, `status`), 
        KEY `idx_status_created` (`status`, `created_at`), 
        KEY `idx_status_deleted` (`status`, `deleted_at`),
        KEY `idx_created_at` (`created_at`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务表';

-- 创建任务信息表
CREATE TABLE
    IF NOT EXISTS `task_info` (
        `task_id` BIGINT UNSIGNED NOT NULL,
        `title` VARCHAR(200) NOT NULL COMMENT '任务标题',
        `description` TEXT DEFAULT NULL COMMENT '任务描述',
        `priority` INT DEFAULT '3' COMMENT '优先级 1:高 2:中 3:低',
        `due_date` DATETIME DEFAULT NULL COMMENT '截止日期',
        `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
        UNIQUE KEY `idx_task_id` (`task_id`), 
        KEY `idx_priority_due` (`priority`, `due_date`), 
        KEY `idx_due_date` (`due_date`), 
        KEY `idx_title` (`title`(50)), 
        CONSTRAINT `fk_task_info_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务信息表';


-- 创建任务历史记录表
CREATE TABLE
    IF NOT EXISTS `task_histories` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `task_id` BIGINT UNSIGNED NOT NULL,
        `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
        `field_name` VARCHAR(50) DEFAULT NULL COMMENT '字段名称',
        `old_value` TEXT DEFAULT NULL COMMENT '旧值',
        `new_value` TEXT DEFAULT NULL COMMENT '新值',
        `operator_id` BIGINT UNSIGNED NOT NULL COMMENT '操作人ID (关联taskhub_user.users.id)',
        `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注',
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        KEY `idx_task_action` (`task_id`, `action`), 
        KEY `idx_task_created` (`task_id`, `created_at`), 
        KEY `idx_operator_created` (`operator_id`, `created_at`), 
        KEY `idx_action_created` (`action`, `created_at`), 
        CONSTRAINT `fk_task_history_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务历史记录表';

-- 创建任务图像表
CREATE TABLE
    IF NOT EXISTS `task_images` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `task_id` BIGINT UNSIGNED NOT NULL,
        `image_url` VARCHAR(500) NOT NULL COMMENT '图像URL路径',
        `sort_order` INT DEFAULT '0' COMMENT '排序顺序',
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `deleted_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_task_sort` (`task_id`, `sort_order`),
        KEY `idx_task_created` (`task_id`, `created_at`),
        CONSTRAINT `fk_task_image_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务图像表';

-- 创建任务历史记录图像表
CREATE TABLE
    IF NOT EXISTS `task_history_images` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `history_id` BIGINT UNSIGNED NOT NULL,
        `image_url` VARCHAR(500) NOT NULL COMMENT '图像URL路径',
        `sort_order` INT DEFAULT '0' COMMENT '排序顺序',
        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `deleted_at` DATETIME DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_history_sort` (`history_id`, `sort_order`),
        KEY `idx_history_created` (`history_id`, `created_at`),
        CONSTRAINT `fk_task_history_image_history` FOREIGN KEY (`history_id`) REFERENCES `task_histories` (`id`) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务历史记录图像表';