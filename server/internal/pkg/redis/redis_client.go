package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"haircut-server/internal/config"
	"haircut-server/pkg/logger"

	"github.com/go-redis/redis/v8"
)

// Client Redis客户端（全局单例）
var Client *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	logger.Info("✅ Redis连接成功")
	return nil
}

// Set 设置缓存（带过期时间）
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string, dest interface{}) error {
	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key不存在: %s", key)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists 检查Key是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := Client.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX 分布式锁（Set if Not eXists）
// 返回：是否获取锁成功、释放锁的函数
func SetNX(ctx context.Context, key string, expiration time.Duration) (bool, func()) {
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
	
	ok, err := Client.SetNX(ctx, key, lockValue, expiration).Result()
	if err != nil || !ok {
		return false, nil
	}

	unlock := func() {
		// Lua脚本确保原子性删除（只删除自己加的锁）
		script := `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end
		`
		Client.Eval(ctx, script, []string{key}, lockValue)
	}

	return true, unlock
}

// Incr 自增计数器
func Incr(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// ZAdd 有序集合添加（排行榜用）
func ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	return Client.ZAdd(ctx, key, members...).Err()
}

// ZRevRangeWithScores 倒序获取有序集合（排行榜TOP N）
func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return Client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// Close 关闭连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
