package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"haircut-server/pkg/logger"
	"haircut-server/pkg/response"

	"github.com/gin-gonic/gin"
	goerrors "github.com/pkg/errors"
)

// Recovery 异常恢复中间件
// 捕获所有未处理的panic，返回统一错误响应，防止服务崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录完整的错误堆栈（用于排查）
				stack := debug.Stack()
				logger.Error("【PANIC】请求异常: %v\n堆栈信息:\n%s", err, string(stack))

				// 根据环境决定返回信息的详细程度
				var errMsg string
				var errDetail interface{}

				// 判断是否为业务可预期错误
				var businessErr *response.BusinessError
				if errors.As(err.(error), &businessErr) {
					errMsg = businessErr.Message
					errDetail = businessErr.Detail
				} else {
					// 系统内部错误：生产环境隐藏细节，开发环境暴露堆栈
					if gin.Mode() == gin.ReleaseMode || gin.Mode() == gin.TestMode {
						errMsg = "服务器内部错误，请联系管理员"
						errDetail = nil
					} else {
						errMsg = "Internal Server Error"
						errDetail = map[string]interface{}{
							"error": goerrors.Wrap(err.(error), "").Error(),
							"stack": string(stack),
						}
					}
				}

				// 返回500响应
				response.Error(c, http.StatusInternalServerError, errMsg, errDetail)
				
				// 终止请求链
				c.Abort()
			}
		}()

		c.Next()
	}
}
