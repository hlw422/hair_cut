package handler

import (
	"haircut-server/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AppointmentHandler 预约相关API处理器
type AppointmentHandler struct{}

// CreateAppointment 创建预约
// POST /api/v1/appointments
// 请求体: { store_id, stylist_id, service_ids: [], appointment_date, appointment_time }
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req struct {
		StoreID        uint64   `json:"store_id" binding:"required"`
		StylistID      uint64   `json:"stylist_id" binding:"required"`
		ServiceIDs     []uint64 `json:"service_ids" binding:"required,min=1"`
		AppointmentDate string   `json:"appointment_date" binding:"required"` // 格式: 2024-01-20
		AppointmentTime string   `json:"appointment_time" binding:"required"` // 格式: "14:00-15:00"
		Remark         string   `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		response.BadRequest(c, "日期格式错误，应为 YYYY-MM-DD")
		return
	}

	// 检查是否为过去的日期
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		response.BadRequest(c, "不能预约过去的日期")
		return
	}

	// TODO: 核心业务逻辑（事务）
	// 1. 检查门店和理发师是否存在且状态正常
	// 2. 检查服务项目是否有效并计算总金额
	// 3. 使用Redis分布式锁检查时间段冲突
	// 4. 创建预约记录
	// 5. 发送确认消息（WebSocket + 微信模板消息）

	mockOrderNo := generateOrderNo("APT")
	
	result := map[string]interface{}{
		"order_no":       mockOrderNo,
		"store_id":       req.StoreID,
		"stylist_id":     req.StylistID,
		"appointment_date": req.AppointmentDate,
		"appointment_time": req.AppointmentTime,
		"total_amount":   198.00,
		"status":         0, // 待确认
		"created_at":     time.Now().Format("2006-01-02 15:04:05"),
		message:          "预约创建成功，请等待确认",
	}

	c.JSON(http.StatusCreated, response.Response{
		Code:    0,
		Message: "预约申请成功",
		Data:    result,
	})
}

// GetAvailableTimeSlots 获取可用时间段
// GET /api/v1/stores/:id/stylists/:stylist_id/slots?date=2024-01-20
func (h *AppointmentHandler) GetAvailableTimeSlots(c *gin.Context) {
	storeID := c.Param("id")
	stylistID := c.Param("stylist_id")
	dateStr := c.Query("date")

	if dateStr == "" {
		response.BadRequest(c, "请提供查询日期(date)")
		return
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "日期格式错误")
		return
	}

	// TODO: 
	// 1. 查询理发师排班表（当天是否上班、工作时段）
	// 2. 查询已有预约（锁定已占用的时间段）
	// 3. 根据服务时长计算可用时段（如剪发1小时，染发2小时）

	_ = storeID
	_ = stylistID

	// Mock返回：09:00-21:00营业时间，每30分钟一个slot，排除已预约的
	slots := []map[string]interface{}{
		{"time": "09:00", "end_time": "10:00", "available": true},
		{"time": "10:00", "end_time": "11:00", "available": true},
		{"time": "11:00", "end_time": "12:00", "available": false, "reason": "已预约"},
		{"time": "14:00", "end_time": "15:00", "available": true},
		{"time": "15:00", "end_time": "16:00", "available": true},
		{"time": "16:00", "end_time": "17:00", "available": true},
		{"time": "18:00", "end_time": "19:00", "available": false, "reason": "休息"},
		{"time": "19:00", "end_time": "20:00", "available": true},
		{"time": "20:00", "end_time": "21:00", "available": true},
	}

	response.Success(c, map[string]interface{}{
		"date":  dateStr,
		"slots": slots,
	})
}

// CancelAppointment 取消预约
// POST /api/v1/appointments/:id/cancel
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	userID := c.GetUint64("user_id")
	aptID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	// TODO:
	// 1. 查询预约记录，验证归属权限
	// 2. 检查状态是否允许取消（待确认/已确认可取消，进行中不可取消）
	// 3. 更新预约状态为已取消
	// 4. 释放Redis锁定的时间段
	// 5. 发送取消通知给理发师/店长

	_ = userID
	_ = aptID

	response.SuccessWithMessage(c, "预约已成功取消", nil)
}

// MyAppointments 我的预约列表
// GET /api/v1/user/appointments?status=&page=1&per_page=20
func (h *AppointmentHandler) MyAppointments(c *gin.Context) {
	userID := c.GetUint64("user_id")
	status := c.Query("status") // 可选过滤: pending/confirmed/completed/cancelled

	// TODO: 调用Service查询当前用户的预约列表（分页）
	_ = userID
	_ = status

	appointments := []map[string]interface{}{
		{
			"id":               1001,
			"order_no":         "APT202401200001",
			"store_name":       "HairCut 精品沙龙（静安店）",
			"stylist_name":     "Kevin老师",
			"stylist_avatar":   "",
			"services":         "首席设计师剪发",
			"appointment_date": "2024-01-22",
			"appointment_time": "14:00-15:00",
			"total_amount":     198.00,
			"status":           1, // 已确认
			"status_text":      "已确认",
			"can_cancel":       true,
		},
	}

	response.Success(c, appointments)
}
