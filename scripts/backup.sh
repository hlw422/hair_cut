#!/bin/bash

# ===========================================
# HairCut 数据备份脚本
# ===========================================

set -e

# 配置
BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_ROOT_USER:-root}"
MYSQL_PASSWORD="${MYSQL_ROOT_PASSWORD:-haircut_root_2024}"
MYSQL_DB="${MYSQL_DB:-haircut}"

# 创建备份目录
mkdir -p $BACKUP_DIR/mysql
mkdir -p $BACKUP_DIR/redis
mkdir -p $BACKUP_DIR/minio

echo "📦 开始数据备份... ($DATE)"

# 1. MySQL 数据库备份
echo "   - 备份 MySQL 数据库..."
mysqldump -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p"$MYSQL_PASSWORD" \
    --single-transaction --routines --triggers \
    --databases $MYSQL_DB | gzip > "$BACKUP_DIR/mysql/haircut_$DATE.sql.gz"
echo "     ✅ MySQL 备份完成: backup_$DATE.sql.gz"

# 2. Redis 备份（通过 docker exec）
if docker ps | grep -q "haircut-redis"; then
    echo "   - 备份 Redis 数据..."
    docker exec haircut-redis redis-cli BGSAVE
    sleep 5
    docker cp haircut-redis:/data/dump.rdb "$BACKUP_DIR/redis/dump_$DATE.rdb"
    echo "     ✅ Redis 备份完成"
fi

# 3. MinIO 文件备份（如果需要）
if docker ps | grep -q "haircut-minio"; then
    echo "   - 提示: MinIO 对象存储建议开启版本控制或使用 mc 工具备份"
    echo "     命令: mc mirror minio/haircut-uploads ./backups/minio/files_$DATE/"
fi

# 清理7天前的备份
echo "🧹 清理旧备份文件（保留最近7天）..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
find $BACKUP_DIR -name "*.rdb" -mtime +7 -delete

echo ""
echo "✅ 数据备份完成！备份位置: $BACKUP_DIR"
echo ""
echo "💡 恢复命令示例:"
echo "   MySQL: gunzip < $BACKUP_DIR/mysql/haircut_$DATE.sql.gz | mysql -u$MYSQL_USER -p$MYSQL_PASSWORD"
echo "   Redis: docker cp $BACKUP_DIR/redis/dump_$date.rdb haircut-redis:/data/dump.rdb && docker restart haircut-redis"
