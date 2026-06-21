package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
// 配置允许的来源、方法、头部等，支持开发环境通配符和生产环境精确控制
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Origin头
		origin := c.GetHeader("Origin")

		// 判断是否为允许的来源（生产环境应配置白名单）
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 允许的方法和头部（预检请求用）
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 预检请求缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400") // 24小时

		// 处理OPTIONS预检请求（直接返回200）
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查是否为允许的跨域来源
func isAllowedOrigin(origin string) bool {
	// 生产环境白名单列表（从配置文件或数据库读取）
	allowedOrigins := []string{
		"http://localhost:5173",      // 开发环境前端
		"http://localhost:3000",       // 开发环境官网
		"http://127.0.0.1:5173",
		"https://admin.your-haircut.com",   // 生产环境后台
		"https://www.your-haircut.com",     // 生产环境官网
		// 微信小程序不需要CORS（不涉及浏览器同源策略）
	}

	// 空来源（如Postman、微信小程序）直接放行
	if origin == "" {
		return true
	}

	for _, allowed := range allowedOrigins {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	// 开发模式：允许所有localhost来源
	if strings.HasPrefix(origin, "http://localhost:") ||
	   strings.HasPrefix(origin, "http://127.0.0.1:") {
		return true
	}

	return false
}
