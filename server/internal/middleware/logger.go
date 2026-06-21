package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"haircut-server/pkg/logger"
	"haircut-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
// 记录每个HTTP请求的详细信息（方法、路径、状态码、耗时、IP等）
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 获取请求追踪ID
		requestID, _ := c.Get("request_id")

		// 记录请求基本信息
		logger.Info("【请求开始】 %s %s | IP: %s | RequestID: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			requestID,
		)

		// 包装responseWriter以捕获响应状态码和大小
		writer := &responseWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = writer

		// 执行请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 根据状态码决定日志级别
		statusCode := writer.statusCode
		switch {
		case statusCode >= 500:
			logger.Error("【请求完成】 %s %s | 状态: %d | 耗时: %v | 大小: %dB | IP: %s | RequestID: %s",
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				duration,
				writer.body.Len(),
				c.ClientIP(),
				requestID,
			)
		case statusCode >= 400:
			logger.Warn("【请求完成】 %s %s | 状态: %d | 耗时: %v | 大小: %dB | IP: %s | RequestID: %s",
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				duration,
				writer.body.Len(),
				c.ClientIP(),
				requestID,
			)
		default:
			logger.Info("【请求完成】 %s %s | 状态: %d | 耗时: %v | 大小: %dB | IP: %s | RequestID: %s",
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				duration,
				writer.body.Len(),
				c.ClientIP(),
				requestID,
			)
		}
	}
}

// responseWriter 用于捕获响应状态码和大小的包装器
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

// bodyLogger 请求体日志记录（可选，用于调试）
func bodyLogger(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}

	// 读取并缓存请求体（因为只能读取一次）
	var bodyBytes []byte
	if c.Request.GetHeader("Content-Type") == "application/json" {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新设置，供后续读取

		logger.Debug("【请求体】 %s", string(bodyBytes))
	} else if c.Request.Method == "GET" {
		logger.Debug("【Query参数】 %s", c.Request.URL.RawQuery)
	}
}
