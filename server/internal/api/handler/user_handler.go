package handler

import (
	"haircut-server/internal/model/mysql"
	"haircut-server/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户相关API处理器
type UserHandler struct {
	// TODO: 注入UserService依赖
}

// GetProfile 获取当前用户个人信息
// GET /api/v1/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// TODO: 调用Service获取用户详情
	user := &mysql.UserPublic{
		ID:       userID.(uint64),
		Nickname: "测试用户",
		AvatarURL: "",
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户个人信息
// PUT /api/v1/user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req struct {
		Nickname  string `json:"nickname" binding:"required,max=50"`
		AvatarURL string `json:"avatar_url"`
		Gender    int8   `json:"gender"`
		Birthday  string `json:"birthday"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// TODO: 调用Service更新用户信息

	response.SuccessWithMessage(c, "个人信息更新成功", nil)
}

// GetMemberInfo 获取会员信息（等级、积分、余额等）
// GET /api/v1/user/member
func (h *UserHandler) GetMemberInfo(c *gin.Context) {
	userID := c.GetUint64("user_id")

	// TODO: 调用MembershipService查询会员信息
	memberInfo := map[string]interface{}{
		"level":        2,
		"level_name":   "银卡会员",
		"points":       1580,
		"balance":      128.50,
		"total_spent":  3580.00,
		"order_count":  12,
		"next_level":   "金卡会员",
		"progress":     71.6, // 升级进度百分比
	}

	response.Success(c, memberInfo)
}
