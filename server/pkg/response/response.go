package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务状态码（0=成功）
	Message string      `json:"message"`            // 响应消息
	Data    interface{} `json:"data,omitempty"`      // 数据（可选）
	Meta    *Meta       `json:"meta,omitempty"`     // 元数据（分页等，可选）
	Error   *ErrorInfo  `json:"error,omitempty"`    // 错误详情（仅错误时返回）
}

// Meta 分页元数据
type Meta struct {
	Page         int `json:"page"`          // 当前页码
	PerPage      int `json:"per_page"`      // 每页数量
	Total        int64 `json:"total"`       // 总记录数
	TotalPages   int  `json:"total_pages"`  // 总页数
}

// ErrorInfo 错误详情
type ErrorInfo {
	Field   string `json:"field,omitempty"`    // 字段名（表单校验错误时使用）
	Message string `json:"message"`             // 错误描述
	Detail  interface{} `json:"detail,omitempty"` // 详细信息（开发环境）
}

// BusinessError 业务异常（可被Recovery中间件捕获并友好处理）
type BusinessError struct {
	Code    int
	Message string
	Detail  interface{}
	Error   error
}

func (e *BusinessError) Error() string {
	return e.Message
}

// 预定义业务错误码
var (
	ErrSuccess           = 0     // 成功
	ErrBadRequest        = 400   // 请求参数错误
	ErrUnauthorized      = 401   // 未认证
	ErrForbidden         = 403   // 无权限
	ErrNotFound          = 404   // 资源不存在
	ErrMethodNotAllowed  = 405   // 方法不允许
	ErrTooManyRequests   = 429   // 请求过于频繁
	ErrInternalServer    = 500   // 内部服务器错误

	// 业务错误（10000+）
	ErrUserNotFound      = 10001 // 用户不存在
	ErrPasswordWrong     = 10002 // 密码错误
	ErrTokenInvalid      = 10003 // Token无效或过期
	ErrPhoneExists       = 10004 // 手机号已注册
	ErrStoreNotFound     = 20001 // 门店不存在
	ErrStylistNotFound   = 30001 // 理发师不存在
	ErrAppointmentConflict = 40001 // 时间段已被预约
	ErrOrderNotPaid      = 50001 // 订单未支付
	ErrInsufficientBalance = 60001 // 余额不足
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    ErrSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    ErrSuccess,
		Message: message,
		Data:    data,
	})
}

// PagedResponse 分页成功响应
func PagedResponse(c *gin.Context, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, Response{
		Code:    ErrSuccess,
		Message: "success",
		Data:    data,
		Meta:    &meta,
	})
}

// Error 错误响应（通用）
func Error(c *gin.Context, httpCode int, message string, detail interface{}) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    nil,
		Error: &ErrorInfo{
			Message: message,
			Detail:  detail,
		},
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message, nil)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// ValidationError 表单验证失败错误（带字段级错误详情）
func ValidationError(c *gin.Context, errors []ErrorInfo) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    ErrBadRequest,
		Message: "请求参数验证失败",
		Error: &ErrorInfo{
			Message: "字段验证未通过",
			Detail:  errors,
		},
	})
}
