#!/bin/bash

# ===========================================
# HairCut 数据库初始化脚本
# ===========================================

set -e

echo "🚀 开始初始化 HairCut 数据库..."

# 配置（可通过环境变量覆盖）
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_ROOT_USER="${MYSQL_ROOT_USER:-root}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-haircut_root_2024}"
MYSQL_DB="${MYSQL_DB:-haircut}"
MYSQL_USER="${MYSQL_USER:-haircut_user}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-haircut_dev_2024}"

# 等待 MySQL 就绪
echo "⏳ 等待 MySQL 服务就绪..."
until mysqladmin ping -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" --silent; do
    sleep 1
done
echo "✅ MySQL 已就绪"

# 创建数据库和用户
echo "📦 创建数据库和用户..."
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" <<EOF
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS \`${MYSQL_DB}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS '${MYSQL_USER}'@'%' IDENTIFIED BY '${MYSQL_PASSWORD}';

-- 授权
GRANT ALL PRIVILEGES ON \`${MYSQL_DB}\`.* TO '${MYSQL_USER}'@'%';

-- 刷新权限
FLUSH PRIVILEGES;

-- 选择数据库
USE \`${MYSQL_DB}\`;

-- ===========================================
-- 以下是核心表结构初始化 SQL
-- （完整表结构将由 GORM AutoMigrate 自动生成）
-- ===========================================

-- 用户表
CREATE TABLE IF NOT EXISTS \`users\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`tenant_id\` BIGINT UNSIGNED DEFAULT NULL,
    \`openid\` VARCHAR(100) DEFAULT NULL COMMENT '微信OpenID',
    \`union_id\` VARCHAR(100) DEFAULT NULL COMMENT '微信UnionID',
    \`phone\` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    \`nickname\` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    \`avatar_url\` TEXT DEFAULT NULL COMMENT '头像URL',
    \`gender\` TINYINT DEFAULT 0 COMMENT '性别: 0未知 1男 2女',
    \`birthday\` DATE DEFAULT NULL COMMENT '生日',
    \`city_code\` VARCHAR(10) DEFAULT NULL COMMENT '城市代码',
    \`status\` TINYINT DEFAULT 1 COMMENT '状态: 0禁用 1正常',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    \`deleted_at\` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (\`id\`),
    UNIQUE KEY \`idx_openid\` (\`openid\`),
    UNIQUE KEY \`idx_phone\` (\`phone\`),
    KEY \`idx_tenant\` (\`tenant_id\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 门店表
CREATE TABLE IF NOT EXISTS \`stores\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`tenant_id\` BIGINT UNSIGNED DEFAULT NULL,
    \`org_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '组织ID',
    \`name\` VARCHAR(100) NOT NULL COMMENT '门店名称',
    \`logo_url\` TEXT DEFAULT NULL COMMENT 'Logo图片URL',
    \`cover_images\` JSON DEFAULT NULL COMMENT '封面图片列表',
    \`province\` VARCHAR(50) DEFAULT NULL COMMENT '省',
    \`city\` VARCHAR(50) DEFAULT NULL COMMENT '市',
    \`district\` VARCHAR(50) DEFAULT NULL COMMENT '区',
    \`address\` VARCHAR(200) NOT NULL COMMENT '详细地址',
    \`latitude\` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
    \`longitude\` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
    \`phone\` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    \`open_time\` VARCHAR(20) DEFAULT '09:00' COMMENT '营业开始时间',
    \`close_time\` VARCHAR(20) DEFAULT '21:00' COMMENT '营业结束时间',
    \`description\` TEXT DEFAULT NULL COMMENT '门店描述',
    \`avg_price\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '人均消费',
    \`rating\` DECIMAL(2,1) DEFAULT 5.0 COMMENT '评分(1-5)',
    \`review_count\` INT DEFAULT 0 COMMENT '评价数',
    \`star_level\` TINYINT DEFAULT 1 COMMENT '星级(1-5)',
    \`status\` TINYINT DEFAULT 1 COMMENT '状态: 0停业 1营业 2筹备中',
    \`is_featured\` TINYINT DEFAULT 0 COMMENT '是否推荐',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    \`deleted_at\` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (\`id\`),
    KEY \`idx_city\` (\`city\`),
    KEY \`idx_location\` (\`latitude\`, \`longitude\`),
    KEY \`idx_status\` (\`status\`),
    KEY \`idx_tenant\` (\`tenant_id\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='门店表';

-- 理发师表
CREATE TABLE IF NOT EXISTS \`stylists\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`tenant_id\` BIGINT UNSIGNED DEFAULT NULL,
    \`store_id\` BIGINT UNSIGNED NOT NULL COMMENT '所属门店ID',
    \`user_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联用户ID',
    \`name\` VARCHAR(50) NOT NULL COMMENT '姓名',
    \`avatar_url\` TEXT DEFAULT NULL COMMENT '头像URL',
    \`gender\` TINYINT DEFAULT 0 COMMENT '性别',
    \`phone\` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    \`title\` VARCHAR(50) DEFAULT NULL COMMENT '职称(首席/高级/资深等)',
    \`experience_years\` INT DEFAULT 0 COMMENT '从业年限',
    \`specialties\` JSON DEFAULT NULL COMMENT '擅长风格标签',
    \`introduction\` TEXT DEFAULT NULL COMMENT '个人简介',
    \`portfolio_count\` INT DEFAULT 0 COMMENT '作品数量',
    \`fan_count\` INT DEFAULT 0 COMMENT '粉丝数',
    \`appointment_count\` INT DEFAULT 0 COMMENT '预约次数',
    \`rating\` DECIMAL(2,1) DEFAULT 5.0 COMMENT '评分',
    \`level\` TINYINT DEFAULT 1 COMMENT '等级(1-5)',
    \`status\` TINYINT DEFAULT 1 COMMENT '状态: 0离职 1在职 2休假',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    \`deleted_at\` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (\`id\`),
    KEY \`idx_store\` (\`store_id\`),
    KEY \`idx_user\` (\`user_id\`),
    KEY \`idx_tenant\` (\`tenant_id\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='理发师表';

-- 服务项目表
CREATE TABLE IF NOT EXISTS \`service_items\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`store_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '门店ID(NULL表示通用项目)',
    \`category_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '分类ID',
    \`name\` VARCHAR(100) NOT NULL COMMENT '服务名称',
    \`description\` TEXT DEFAULT NULL COMMENT '服务描述',
    \`image_url\` TEXT DEFAULT NULL COMMENT '服务图片',
    \`original_price\` DECIMAL(10,2) NOT NULL COMMENT '原价',
    \`price\` DECIMAL(10,2) NOT NULL COMMENT '现价',
    \`duration\` INT DEFAULT 30 COMMENT '服务时长(分钟)',
    \`target_gender\` TINYINT DEFAULT 0 COMMENT '适用性别: 0通用 1男 2女',
    \`sort_order\` INT DEFAULT 0 COMMENT '排序权重',
    \`is_hot\` TINYINT DEFAULT 0 COMMENT '是否热门',
    \`status\` TINYINT DEFAULT 1 COMMENT '状态: 0下架 1上架',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    \`deleted_at\` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (\`id\`),
    KEY \`idx_store\` (\`store_id\`),
    KEY \`idx_category\` (\`category_id\`),
    KEY \`idx_hot\` (\`is_hot\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务项目表';

-- 预约记录表
CREATE TABLE IF NOT EXISTS \`appointments\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`tenant_id\` BIGINT UNSIGNED DEFAULT NULL,
    \`order_no\` VARCHAR(64) NOT NULL COMMENT '预约单号',
    \`user_id\` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    \`store_id\` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
    \`stylist_id\` BIGINT UNSIGNED NOT NULL COMMENT '理发师ID',
    \`appointment_date\` DATE NOT NULL COMMENT '预约日期',
    \`appointment_time\` VARCHAR(10) NOT NULL COMMENT '预约时间段(如14:00-15:00)',
    \`service_ids\` JSON DEFAULT NULL COMMENT '服务项目ID列表',
    \`total_amount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '总金额',
    \`status\` TINYINT DEFAULT 0 COMMENT '状态: 0待确认 1已确认 2进行中 3已完成 4已取消',
    \`remark\` VARCHAR(500) DEFAULT NULL COMMENT '备注',
    \`source\` TINYINT DEFAULT 1 COMMENT '来源: 1小程序 2电话 3到店',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (\`id\`),
    UNIQUE KEY \`idx_order_no\` (\`order_no\`),
    KEY \`idx_user\` (\`user_id\`),
    KEY \`idx_store\` (\`store_id\`),
    KEY \`idx_stylist\` (\`stylist_id\`),
    KEY \`idx_datetime\` (\`appointment_date\`, \`appointment_time\`),
    KEY \`idx_status\` (\`status\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约记录表';

-- 订单表
CREATE TABLE IF NOT EXISTS \`orders\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`tenant_id\` BIGINT UNSIGNED DEFAULT NULL,
    \`order_no\` VARCHAR(64) NOT NULL COMMENT '订单号',
    \`user_id\` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    \`store_id\` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
    \`stylist_id\` BIGINT UNSIGNED NOT NULL COMMENT '理发师ID',
    \`appointment_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联预约ID',
    \`total_amount\` DECIMAL(10,2) NOT NULL COMMENT '订单总额',
    \`discount_amount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '优惠减免',
    \`coupon_amount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '优惠券抵扣',
    \`points_amount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '积分抵扣',
    \`pay_amount\` DECIMAL(10,2) NOT NULL COMMENT '实付金额',
    \`pay_method\` TINYINT DEFAULT 1 COMMENT '支付方式: 1微信 2余额 3积分 4混合',
    \`pay_time\` TIMESTAMP NULL DEFAULT NULL COMMENT '支付时间',
    \`transaction_id\` VARCHAR(64) DEFAULT NULL COMMENT '微信交易号',
    \`status\` TINYINT DEFAULT 0 COMMENT '状态: 0待支付 1已支付 2服务中 3已完成 4已退款 5取消',
    \`refund_amount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '退款金额',
    \`remark\` VARCHAR(500) DEFAULT NULL COMMENT '备注',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    \`deleted_at\` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (\`id\`),
    UNIQUE KEY \`idx_order_no\` (\`order_no\`),
    KEY \`idx_user\` (\`user_id\`),
    KEY \`idx_store\` (\`store_id\`),
    KEY \`idx_status\` (\`status\`),
    KEY \`idx_pay_time\` (\`pay_time\`),
    KEY \`idx_tenant\` (\`tenant_id\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 会员信息表
CREATE TABLE IF NOT EXISTS \`members\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`user_id\` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    \`level\` TINYINT DEFAULT 1 COMMENT '等级: 1普通 2银卡 3金卡 4黑金 5钻石',
    \`points\` INT DEFAULT 0 COMMENT '积分余额',
    \`balance\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '余额',
    \`total_spent\` DECIMAL(12,2) DEFAULT 0.00 COMMENT '累计消费',
    \`order_count\` INT DEFAULT 0 COMMENT '订单数',
    \`expire_date\` DATE DEFAULT NULL COMMENT '有效期至',
    \`upgrade_threshold\` DECIMAL(12,2) DEFAULT 0.00 COMMENT '升级门槛',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (\`id\`),
    UNIQUE KEY \`idx_user\` (\`user_id\`),
    KEY \`idx_level\` (\`level\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员信息表';

-- 优惠券模板表
CREATE TABLE IF NOT EXISTS \`coupon_templates\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`name\` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
    \`type\` TINYINT NOT NULL COMMENT '类型: 1满减 2折扣 3新人券 4节日券 5门店券 6通用券',
    \`value\` DECIMAL(10,2) NOT NULL COMMENT '面值/折扣率',
    \`min_spend\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '最低消费',
    \`max_discount\` DECIMAL(10,2) DEFAULT 0.00 COMMENT '最大优惠金额(折扣券使用)',
    \`valid_days\` INT DEFAULT 30 COMMENT '有效天数',
    \`total_count\` INT DEFAULT 0 COMMENT '发放总量(-1不限)',
    \`per_user_limit\` INT DEFAULT 1 COMMENT '每人限领数量',
    \`scope\` TINYINT DEFAULT 1 COMMENT '适用范围: 1全国 2指定门店 3指定服务',
    \`scope_config\` JSON DEFAULT NULL COMMENT '范围配置',
    \`status\` TINYINT DEFAULT 1 COMMENT '状态: 0停用 1启用',
    \`start_time\` TIMESTAMP NULL DEFAULT NULL COMMENT '活动开始时间',
    \`end_time\` TIMESTAMP NULL DEFAULT NULL COMMENT '活动结束时间',
    \`created_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    \`updated_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (\`id\`),
    KEY \`idx_type\` (\`type\`),
    KEY \`idx_status\` (\`status\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券模板表';

-- 用户优惠券表
CREATE TABLE IF NOT EXISTS \`coupons\` (
    \`id\` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    \`template_id\` BIGINT UNSIGNED NOT NULL COMMENT '模板ID',
    \`user_id\` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    \`code\` VARCHAR(32) NOT NULL COMMENT '优惠券码',
    \`status\` TINYINT DEFAULT 0 COMMENT '状态: 0未使用 1已使用 2已过期',
    \`obtained_at\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
    \`used_at\` TIMESTAMP NULL DEFAULT NULL COMMENT '使用时间',
    \`used_order_id\` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用订单ID',
    \`expire_at\` TIMESTAMP NOT NULL COMMENT '过期时间',
    PRIMARY KEY (\`id\`),
    UNIQUE KEY \`idx_code\` (\`code\`),
    KEY \`idx_user\` (\`user_id\`),
    KEY \`idx_template\` (\`template_id\`),
    KEY \`idx_status\` (\`status\`),
    KEY \`idx_expire\` (\`expire_at\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券表';

EOF

echo "✅ 数据库初始化完成！"
echo ""
echo "📊 已创建以下核心数据表:"
echo "   - users (用户表)"
echo "   - stores (门店表)"
echo "   - stylists (理发师表)"
echo "   - service_items (服务项目表)"
echo "   - appointments (预约记录表)"
echo "   - orders (订单表)"
echo "   - members (会员信息表)"
echo "   - coupon_templates (优惠券模板表)"
echo "   - coupons (用户优惠券表)"
echo ""
echo "💡 提示: 完整表结构将通过 GORM AutoMigrate 自动生成剩余30+张业务表"
