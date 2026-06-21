package middleware

import (
	"net/http"
	"strings"

	"haircut-server/internal/pkg/casbin"
	"haircut-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// RBAC Casbin权限校验中间件
// 基于用户角色+请求路径+方法进行访问控制
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取当前用户ID和角色列表（从JWT中间件注入的上下文）
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusForbidden, "无法获取用户身份", nil)
			c.Abort()
			return
		}

		rolesInterface, _ := c.Get("roles")
		var roles []string
		if rolesInterface != nil {
			if r, ok := rolesInterface.([]string); ok {
				roles = r
			}
		}

		// 2. 如果没有角色，拒绝访问（超级管理员在JWT中已包含角色）
		if len(roles) == 0 {
			response.Error(c, http.StatusForbidden, "无访问权限", nil)
			c.Abort()
			return
		}

		// 3. 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 4. 逐个检查角色是否有权限（任一角色有权限即可）
		hasPermission := false
		for _, role := range roles {
			ok, err := casbin.Enforcer.Enforce(role, path, method)
			if err != nil {
				continue // 出错时跳过，检查下一个角色
			}
			if ok {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, "无权执行此操作", map[string]interface{}{
				"path":   path,
				"method": method,
				"roles":  roles,
			})
			c.Abort()
			return
		}

		// 5. 数据权限过滤（可选：根据租户/门店/区域）
		dataScopeCheck(c)

		c.Next()
	}
}

// dataScopeCheck 数据权限范围检查
// 确保用户只能访问其权限范围内的数据（门店隔离、区域隔离等）
func dataScopeCheck(c *gin.Context) {
	userID := c.GetUint64("user_id")
	tenantID := c.GetUint64("tenant_id")

	// 将数据权限上下文注入，供后续Repository使用
	c.Set("data_scope", map[string]uint64{
		"user_id":  userID,
		"tenant_id": tenantID,
	})
}

// RequireRole 角色白名单中间件工厂函数
// 用于特定路由需要指定角色的场景
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesInterface, _ := c.Get("roles")
		var userRoles []string
		if rolesInterface != nil {
			if r, ok := rolesInterface.([]string); ok {
				userRoles = r
			}
		}

		// 检查用户是否拥有任一所需角色
		for _, required := range requiredRoles {
			for _, has := range userRoles {
				if strings.EqualFold(has, required) {
					c.Next() // 有权限，放行
					return
				}
			}
		}

		response.Error(c, http.StatusForbidden, "需要"+strings.Join(requiredRoles, "或")+"角色", nil)
		c.Abort()
	}
}
