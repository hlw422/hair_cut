package mysql

import (
	"time"
)

// Order 订单表 - 交易订单核心表
type Order struct {
	BaseModel
	// 订单编号
	OrderNo string `json:"order_no" gorm:"type:varchar(64);uniqueIndex;not null;comment:订单号"`

	// 关键实体关联
	UserID       uint64  `json:"user_id" gorm:"not null;index;comment:用户ID"`
	StoreID      uint64  `json:"store_id" gorm:"not null;index;comment:门店ID"`
	StylistID    uint64  `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`
	AppointmentID *uint64 `json:"appointment_id,omitempty" gorm:"index;comment:关联预约ID(预约转订单)"`

	// 金额明细（单位：元）
	TotalAmount   float64 `json:"total_amount" gorm:"type:decimal(10,2);not null;comment:订单总额(元)"`
	DiscountAmount float64 `json:"discount_amount" gorm:"type:decimal(10,2);default:0.00;comment:会员折扣减免(元)"`
	CouponAmount  float64 `json:"coupon_amount" gorm:"type:decimal(10,2);default:0.00;comment:优惠券抵扣(元)"`
	PointsAmount  float64 `json:"points_amount" gorm:"type:decimal(10,2);default:0.00;comment:积分抵扣(元)"`
	PayAmount     float64 `json:"pay_amount" gorm:"type:decimal(10,2);not null;index;comment:实付金额(元)"`

	// 支付信息
	PayMethod    int8        `json:"pay_method" gorm:"type:tinyint;default:1;comment:支付方式: 1微信支付 2余额 3积分全额 4混合支付"`
	PayTime      *time.Time  `json:"pay_time,omitempty" gorm:"comment:支付时间"`
	TransactionID string     `json:"transaction_id,omitempty" gorm:"type:varchar(64);index;comment:微信/支付宝交易号"`

	// 状态机（核心字段）
	Status int8 `json:"status" gorm:"type:tinyint;default:0;index;comment:订单状态: 0待支付 1已支付 2服务中 3已完成 4已退款(部分) 5已退款(全部) 6取消"`

	// 退款信息
	RefundAmount   float64    `json:"refund_amount" gorm:"type:decimal(10,2);default:0.00;comment:退款金额(元)"`
	RefundReason   string     `json:"refund_reason,omitempty" gorm:"type:varchar(200);comment:退款原因"`
	RefundTime     *time.Time `json:"refund_time,omitempty" gorm:"comment:退款时间"`
	RefundTransactionID string `json:"refund_transaction_id,omitempty" gorm:"type:varchar(64);comment:退款交易号"`

	// 备注与来源
	Remark string `json:"remark" gorm:"type:varchar(500);comment:用户备注/订单备注"`
	Source int8   `json:"source" gorm:"type:tinyint;default:1;comment:来源: 1小程序预约 2到店消费 3电话预订 4后台创建"`

	// 关联关系
	User           *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Store          *Store            `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Stylist        *Stylist          `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
	Appointment    *Appointment      `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	OrderItems     []OrderItem       `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`         // 订单项列表
	PaymentRecords []PaymentRecord   `json:"payment_records,omitempty" gorm:"foreignKey:OrderID"`    // 支付记录
	Review         *Review           `json:"review,omitempty" gorm:"foreignKey:OrderID"`            // 评价记录
}

func (Order) TableName() string {
	return "orders"
}

// 订单状态常量
const (
	OrderStatusPendingPayment = 0 // 待支付
	OrderStatusPaid           = 1 // 已支付（待服务）
	OrderStatusInProgress     = 2 // 服务中
	OrderStatusCompleted      = 3 // 已完成
	OrderStatusPartialRefund  = 4 // 已部分退款
	OrderStatusFullRefund     = 5 // 已全额退款
	OrderStatusCancelled      = 6 // 已取消
)

// IsPaid 检查是否已支付
func (o *Order) IsPaid() bool {
	return o.Status >= OrderStatusPaid && o.Status <= OrderStatusCompleted
}

// CanRefund 检查是否可申请退款
func (o *Order) CanRefund() bool {
	return o.Status == OrderStatusPaid || o.Status == OrderStatusInProgress || o.Status == OrderStatusCompleted
}

// CanComment 检查是否可以评价（已完成且未评价）
func (o *Order) CanComment() bool {
	return o.Status == OrderStatusCompleted && o.Review == nil
}

// OrderItem 订单明细表 - 订单中的具体服务项目
type OrderItem struct {
	BaseModel
	OrderID     uint64  `json:"order_id" gorm:"not null;index;comment:订单ID"`
	ServiceItemID uint64 `json:"service_item_id" gorm:"not null;comment:服务项目ID"`
	ServiceName string  `json:"service_name" gorm:"type:varchar(100);not null;comment:服务名称(冗余，防止原项目删除)"`
	Quantity    int     `json:"quantity" gorm:"type:int;default:1;comment:数量"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(10,2);not null;comment:单价(元)"`
	TotalPrice  float64 `json:"total_price" gorm:"type:decimal(10,2);not null;comment:小计金额(元)"`

	// 关联
	Order       *Order       `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ServiceItem *ServiceItem `json:"service_item,omitempty" gorm:"foreignKey:ServiceItemID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

// PaymentRecord 支付记录表 - 支付流水（支持多次支付/部分退款）
type PaymentRecord struct {
	BaseModel
	OrderID       uint64    `json:"order_id" gorm:"not null;index;comment:订单号"`
	TransactionNo string    `json:"transaction_no" gorm:"type:varchar(64);uniqueIndex;not null;comment:支付流水号"`
	Type          int8      `json:"type" gorm:"type:tinyint;default:1;comment:类型: 1收入(支付) 2支出(退款)"`
	Method        int8      `json:"method" gorm:"type:tinyint;not null;comment:方式: 1微信 2余额 3积分"`
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:金额(元)"`
	Status        int8      `json:"status" gorm:"type:tinyint;default:0;index;comment:状态: 0待处理 1成功 2失败 3已关闭"`
	PayTime       *time.Time `json:"pay_time,omitempty" gorm:"comment:完成时间"`
	ThirdPartyNo  string    `json:"third_party_no" gorm:"type:varchar(64);comment:第三方交易号"`
	RawData       string    `json:"-" gorm:"type:text;comment:原始回调数据(JSON)"` // 不返回给前端

	// 关联
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

func (PaymentRecord) TableName() string {
	return "payment_records"
}
