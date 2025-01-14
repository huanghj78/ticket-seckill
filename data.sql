SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for goods
-- ----------------------------
DROP TABLE IF EXISTS `goods`;
CREATE TABLE `goods` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '商品id',
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `img` varchar(255) NOT NULL COMMENT '商品图片',
  `origin_price` bigint(20) NOT NULL COMMENT '商品价格',
  `price` bigint(20) NOT NULL COMMENT '秒杀价格',
  `stock` int(11) unsigned NOT NULL COMMENT '库存',
  `start_time` datetime NOT NULL COMMENT '秒杀开始时间',
  `end_time` datetime NOT NULL COMMENT '秒杀结束时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- ----------------------------
-- Records of goods
-- ----------------------------
INSERT INTO `goods` (`name`, `img`, `origin_price`, `price`, `stock`, `start_time`, `end_time`) VALUES
('魔力不粘锅', '', 1200, 600, 150, '2024-01-20 10:00:00', '2026-01-20 12:00:00'),
('智能食物称重器', '', 800, 450, 200, '2024-01-21 11:00:00', '2026-01-21 13:00:00'),
('空气炸锅', '', 1800, 1000, 80, '2024-01-22 09:00:00', '2026-01-22 11:00:00'),
('智能调味瓶套装', '', 650, 350, 300, '2024-01-23 14:00:00', '2026-01-23 16:00:00'),
('多功能料理机', '', 1500, 850, 120, '2024-01-24 15:00:00', '2026-01-24 17:00:00'),
('高端刀具套装', '', 2500, 1300, 100, '2024-01-25 08:00:00', '2026-01-25 10:00:00'),
('自清洁微波炉', '', 3200, 1600, 60, '2024-01-26 13:00:00', '2026-01-26 15:00:00'),
('便捷料理切割板', '', 500, 250, 400, '2024-01-27 12:00:00', '2026-01-27 14:00:00'),
('自动搅拌锅', '', 2200, 1100, 50, '2024-01-28 17:00:00', '2026-01-28 19:00:00'),
('厨房智能垃圾桶', '', 1500, 800, 30, '2024-01-29 18:00:00', '2026-01-29 20:00:00');


-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_id` varchar(32) NOT NULL COMMENT '订单id',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `goods_id` bigint(20) NOT NULL COMMENT '商品id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_gid` (`user_id`,`goods_id`)
) ENGINE=InnoDB AUTO_INCREMENT=313 DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- ----------------------------
-- Records of order
-- ----------------------------

-- ----------------------------
-- Table structure for order_info
-- ----------------------------
DROP TABLE IF EXISTS `order_infos`;
CREATE TABLE `order_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `order_id` varchar(32) NOT NULL COMMENT '订单号',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `goods_id` bigint(20) unsigned NOT NULL COMMENT '商品id',
  `goods_name` varchar(128) NOT NULL COMMENT '商品名称',
  `goods_img` varchar(128) NOT NULL COMMENT '商品图片',
  `goods_price` bigint(20) unsigned NOT NULL COMMENT '商品价格',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '订单状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_id` (`order_id`),
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=474 DEFAULT CHARSET=utf8mb4 COMMENT='订单信息表';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `mobile` varchar(16) NOT NULL COMMENT '手机号码',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_mobile` (`mobile`) USING BTREE COMMENT '手机号码索引'
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
