package middleware

import (
	"net/http"
	"strings"
	"time"

	"haircut-server/internal/pkg/jwt"
	"haircut-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
// 从请求头提取Bearer Token，解析并验证有效性，将用户信息注入上下文
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未提供认证令牌", nil)
			c.Abort()
			return
		}

		// 2. 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "认证格式错误", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 解析并验证Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "令牌无效或已过期", nil)
			c.Abort()
			return
		}

		// 4. 检查是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			response.Error(c, http.StatusUnauthorized, "令牌已过期，请重新登录", nil)
			c.Abort()
			return
		}

		// 5. 将用户信息注入上下文（后续Handler可获取）
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("tenant_id", claims.TenantID)

		// 继续执行下一个中间件/Handler
		c.Next()
	}
}

// OptionalJWT 可选的JWT认证（部分公开接口可能携带Token以获得更多信息）
func OptionalJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				if claims, err := jwt.ParseToken(parts[1]); err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
				}
			}
		}
		c.Next()
	}
}
