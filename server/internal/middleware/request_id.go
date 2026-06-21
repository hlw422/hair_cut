package middleware

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// RequestID 请求追踪ID中间件
// 为每个请求生成唯一标识，用于链路追踪、日志关联、问题排查
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取（前端或网关传入）
		requestID := c.GetHeader("X-Request-ID")
		
		// 如果未提供，生成新的UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 注入到上下文（后续Handler/日志可获取）
		c.Set("request_id", requestID)

		// 设置响应头（方便客户端关联日志）
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
