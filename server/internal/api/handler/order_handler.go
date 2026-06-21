package handler

import (
	"haircut-server/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderHandler 订单相关API处理器
type OrderHandler struct{}

// CreateOrder 从预约创建订单（或到店消费直接下单）
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req struct {
		// 预约转订单
		AppointmentID *uint64 `json:"appointment_id,omitempty"`
		
		// 或直接下单（到店消费）
		StoreID     uint64   `json:"store_id" binding:"required"`
		StylistID   uint64   `json:"stylist_id" binding:"required"`
		Items       []struct {
			ServiceItemID uint64  `json:"service_item_id" binding:"required"`
			Quantity      int     `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		
		// 支付与优惠
		CouponCode    *string `json:"coupon_code,omitempty"`    // 优惠券码
		UsePoints     int     `json:"use_points,omitempty"`    // 使用积分数量
		PayMethod     int8    `json:"pay_method" binding:"required"` // 支付方式: 1微信 2余额
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	// TODO: 核心业务流程（事务）
	// 1. 验证门店/理发师/服务项目有效性
	// 2. 计算订单金额明细：
	//    - 原价总额 = Σ(单价 × 数量)
	//    - 会员折扣减免 = 原价总额 × (1 - 会员折扣率)
	//    - 优惠券抵扣 = 查询coupon表计算抵扣金额
	//    - 积分抵扣 = use_points × 积分汇率(如100积分=1元)
	//    - 实付金额 = 折后金额 - 优惠券 - 积分抵扣
	// 3. 检查用户余额是否足够（如果选择余额支付）
	// 4. 创建订单记录(状态=待支付) + 订单项记录
	// 5. 如果使用优惠券，预扣减优惠券状态
	// 6. 如果是预约转订单，更新预约关联

	orderNo := generateOrderNo("ORD")
	
	result := map[string]interface{}{
		"order_no":        orderNo,
		"total_amount":    268.00,
		"discount_amount": 26.80, // 9折会员折扣
		"coupon_amount":   20.00, // 使用了满200减20券
		"points_amount":   0.00,  // 本次未使用积分
		"pay_amount":      221.20, // 最终实付
		"status":          0,      // 待支付
		"expire_time":     "2024-01-20 15:30:00", // 订单过期时间(30分钟内支付)
		// 微信支付参数（前端调起支付用）
		"payment_params": map[string]string{
			"timeStamp": "1705746600",
			"nonceStr":  "random_string",
			"packageValue: "prepay_id=wx201410272009395522657a670389285100",
			"signType":  "MD5",
			"paySign":    "signature_here",
		},
	}

	c.JSON(http.StatusCreated, response.Response{
		Code:    0,
		Message: "订单创建成功，请完成支付",
		Data:    result,
	})
}

// PayOrder 发起支付
// POST /api/v1/orders/:id/pay
func (h *OrderHandler) PayOrder(c *gin.Context) {
	userID := c.GetUint64("user_id")
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		PayMethod int8 `json:"pay_method" binding:"required"` // 1微信支付 2余额
	}
	c.ShouldBindJSON(&req)

	// TODO:
	// 1. 查询订单，验证归属权限和状态（必须是待支付状态）
	// 2. 根据PayMethod处理：
	//    a) 微信支付：调用微信统一下单API获取prepay_id，返回给前端调起支付
	//    b) 余额支付：检查用户余额是否充足，扣除余额，直接完成支付流程
	// 3. 创建PaymentRecord记录

	_ = userID
	_ = orderID

	if req.PayMethod == 1 {
		// 返回微信支付参数
		response.Success(c, map[string]interface{}{
			"payment_type": "wechat",
			"params": map[string]string{
				"appId":     "wx1234567890",
				"timeStamp": "1705746600",
			"nonceStr":  "abcxyz",
			"packageValue: "prepay_id=xxx",
			"signType":  "MD5",
			"paySign":   "signed_value",
		},
		})
	} else if req.PayMethod == 2 {
		// 余额支付成功（同步返回）
		response.SuccessWithMessage(c, "支付成功", map[string]interface{}{
			"order_id": orderID,
			"status":   1, // 已支付
		})
	}
}

// GetOrderDetail 订单详情
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_ = userID

	// TODO: 调用Service查询完整订单信息
	detail := map[string]interface{}{
		"id":              orderID,
		"order_no":        "ORD20240120000123",
		"status":          3, // 已完成
		"status_text":     "已完成",
		"store":           map[string]interface{}{"name": "HairCut 精品沙龙（静安店）"},
		"stylist":         map[string]interface{}{"name": "Kevin老师"},
		"items": []map[string]interface{}{
			{"name": "首席设计师剪发", "price": 198.0, "qty": 1, "subtotal": 198.0},
			{"name": "头皮护理SPA", "price": 70.0, "qty": 1, "subtotal": 70.0},
		},
		"amount_detail": map[string]float64{
			"total_amount":    268.00,
			"discount_amount": 26.80,
			"coupon_amount":   20.00,
			"points_amount":   0.00,
			"pay_amount":      221.20,
		},
		"pay_method_text": "微信支付",
		"pay_time":         "2024-01-20 14:35:22",
		"can_refund":      true,
		"can_comment":      true, // 已完成且未评价
		"created_at":       "2024-01-20 14:05:18",
	}

	response.Success(c, detail)
}

// ListOrders 我的订单列表
// GET /api/v1/user/orders?status=&page=1&per_page=10
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetUint64("user_id")
	status := c.Query("status") // all/pending/paid/completed/refunded

	// TODO: 分页查询用户的订单列表
	_ = userID
	_ = status

	orders := []map[string]interface{}{
		{
			"id":            1001,
			"order_no":      "ORD2024011900001",
			"store_name":    "HairCut 精品沙龙（静安店）",
			"services":      "剪发+染发套餐",
			"pay_amount":    388.00,
			"status":        3,
			"status_text":   "已完成",
			"created_at":    "2024-01-19 16:30:00",
		},
		{
			"id":            1002,
			"order_no":      "ORD2024011800002",
			"store_name":    "HairCut 造型中心（徐汇店）",
			"services":      "男士精剪",
			"pay_amount":    168.00,
			"status":        1,
			"status_text":   "待服务",
			"appointment_date": "2024-01-21 10:00",
			"created_at":    "2024-01-18 09:15:00",
		},
	}

	response.Success(c, orders)
}

// ApplyRefund 申请退款
// POST /api/v1/orders/:id/refund
func (h *OrderHandler) ApplyRefund(c *gin.Context) {
	userID := c.GetUint64("user_id")
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason" binding:"required,max=200"`
		Type    int8   `json:"refund_type"` // 1部分退款 2全额退款
	}
	c.ShouldBindJSON(&req)

	// TODO:
	// 1. 检查订单是否可退款（已完成/已支付的订单可申请）
	// 2. 创建退款工单
	// 3. 调用微信退款API（如果是微信支付）
	// 4. 更新订单退款状态
	// 5. 退还优惠券/积分（如有使用）

	_ = userID
	_ = orderID

	response.SuccessWithMessage(c, "退款申请已提交，预计1-3个工作日原路退回", nil)
}
