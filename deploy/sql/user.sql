CREATE TABLE `users` (
    `id` varchar(24) COMMENT '用户ID',
    `avatar` varchar(191) NOT NULL DEFAULT '' COMMENT '用户头像URL',
    `nickname` varchar(24) NOT NULL COMMENT '用户昵称',
    `phone` varchar(20) NOT NULL COMMENT '手机号码',
    `password` varchar(191) NOT NULL COMMENT '密码(加密存储)',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '用户状态(0-禁用,1-正常,2-锁定)',
    `sex` tinyint NOT NULL DEFAULT 0 COMMENT '性别(0-未知,1-男,2-女)',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone` (`phone`) COMMENT '手机号唯一索引',
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';