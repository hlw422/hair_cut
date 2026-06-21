package middleware

import (
	"net/http"
	"sync"
	"time"

	"haircut-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// rateLimiter 令牌桶/滑动窗口限流器
type rateLimiter struct {
	// 使用sync.Map实现并发安全的请求计数
	requests sync.Map // key: clientIP, value: *rateLimitEntry
	mu        sync.Mutex
	window    time.Duration // 时间窗口（如1分钟）
	maxRequests int         // 窗口内最大请求数
}

type rateLimitEntry struct {
	count      int
	resetTime  time.Time
}

// 全局限流实例（生产环境应使用Redis实现分布式限流）
var globalLimiter = &rateLimiter{
	window:      time.Minute,
	maxRequests: 120, // 每IP每分钟最多120次请求（可根据接口类型调整）
}

// RateLimit 请求频率限制中间件
// 基于客户端IP的滑动窗口算法，防止恶意刷接口/DDOS
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查是否超限
		if !globalLimiter.allowRequest(clientIP) {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试", map[string]interface{}{
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		// 设置响应头告知客户端剩余配额
		entry, _ := globalLimiter.requests.Load(clientIP)
		if e, ok := entry.(*rateLimitEntry); ok {
			remaining := globalLimiter.maxRequests - e.count
			c.Header("X-RateLimit-Limit", string(globalLimiter.maxRequests))
			c.Header("X-RateLimit-Remaining", string(max(0, remaining)))
			c.Header("X-RateLimit-Reset", e.resetTime.Format(time.RFC3339))
		}

		c.Next()
	}
}

// allowRequest 判断是否允许该请求通过
func (rl *rateLimiter) allowRequest(clientIP string) bool {
	now := time.Now()

	// 原子加载或创建条目
	actual, _ := rl.requests.LoadOrStore(clientIP, &rateLimitEntry{
		count:     0,
		resetTime: now.Add(rl.window),
	})
	entry := actual.(*rateLimitEntry)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 检查是否需要重置时间窗口
	if now.After(entry.resetTime) {
		entry.count = 0
		entry.resetTime = now.Add(rl.window)
	}

	// 检查是否超过限制
	if entry.count >= rl.maxRequests {
		return false
	}

	// 增加计数
	entry.count++
	return true
}

// RateLimitByPath 按路径差异化限流（可选）
// 不同API可设置不同阈值（如登录接口更严格）
func RateLimitByPath(limits map[string]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 获取该路径的限流配置，默认使用全局配置
		maxReq, ok := limits[path]
		if !ok {
			maxReq = globalLimiter.maxRequests
		}

		// TODO: 实现基于路径的独立限流逻辑
		_ = maxReq

		c.Next()
	}
}
