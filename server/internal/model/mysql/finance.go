package mysql

import (
	"time"
)

// FinanceRecord 财务记录表 - 收入/支出/退款等流水
type FinanceRecord struct {
	BaseModel
	// 归属
	StoreID   *uint64 `json:"store_id,omitempty" gorm:"index;comment:门店ID(NULL=总部)"`
	StylistID *uint64 `json:"stylist_id,omitempty" gorm:"index;comment:理发师ID(如有)"`

	// 记录信息
	Type       int8    `json:"type" gorm:"type:tinyint;not null;index;comment:类型: 1营业收入 2退款支出 3采购支出 4员工薪资 5其他收入 6其他支出"`
	Category   string  `json:"category" gorm:"type:varchar(50);index;comment:分类(服务销售/产品销售/会员充值等)"`
	Amount     float64 `json:"amount" gorm:"type:decimal(12,2);not null;comment:金额(元)(正数=收入/负数=支出)"`
	Currency   string  `json:"currency" gorm:"type:varchar(10);default:CNY;comment:币种"`

	// 关联业务
	RelatedType string  `json:"related_type" gorm:"type:varchar(50);comment:关联业务类型(order/purchase/payroll等)"`
	RelatedID   *uint64 `json:"related_id,omitempty" gorm:"comment:关联业务ID"`
	OrderNo     string  `json:"order_no" gorm:"type:varchar(64);index;comment:关联单号"`

	// 时间与描述
	TransactionDate time.Time `json:"transaction_date" gorm:"type:date;not null;index;comment:交易日期"`
	Description     string    `json:"description" gorm:"type:text;comment:备注说明"`

	// 审批信息（如需要）
	Status      int8       `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0待审核 1已确认 2已驳回"`
	ApprovedBy  *uint64    `json:"approved_by,omitempty" gorm:"comment:审批人ID"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" gorm:"comment:审批时间"`
}

func (FinanceRecord) TableName() string {
	return "finance_records"
}

const (
	FinanceTypeRevenue        = 1 // 营业收入
	FinanceTypeRefund         = 2 // 退款支出
	FinanceTypePurchase       = 3 // 采购支出
	FinanceTypePayroll        = 4 // 员工薪资
	FinanceTypeOtherIncome    = 5 // 其他收入
	FinanceTypeOtherExpense   = 6 // 其他支出
)

// StylistCommission 理发师提成记录表 - 服务提成结算
type StylistCommission struct {
	BaseModel
	StylistID uint64    `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`
	OrderID   uint64    `json:"order_id" gorm:"not null;index;comment:订单ID"`
	OrderAmount float64 `json:"order_amount" gorm:"type:decimal(10,2);not null;comment:订单金额(元)"`
	Rate      float64   `json:"rate" gorm:"type:decimal(4,3);not null;comment:提成比例(如0.300表示30%)"`
	Commission float64  `json:"commission" gorm:"type:decimal(10,2);not null;comment:提成金额(元)"`
	Status    int8      `json:"status" gorm:"type:tinyint;default:0;index;comment:状态: 0待结算 1已结算 2已发放"`
	SettledAt *time.Time `json:"settled_at,omitempty" gorm:"comment:结算时间"`
	PaidAt    *time.Time `json:"paid_at,omitempty" gorm:"comment:发放时间"`
	PayMethod int8       `json:"pay_method" gorm:"type:tinyint;comment:发放方式: 1银行转账 2现金 3余额"`
	Remark    string    `json:"remark" gorm:"type:varchar(200);comment:备注"`

	// 关联关系
	Stylist *Stylist `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

func (StylistCommission) TableName() string {
	return "stylist_commissions"

const (
	CommissionStatusPending  = 0 // 待结算
	CommissionStatusSettled  = 1 // 已结算
	CommissionStatusPaid     = 2 // 已发放
)
}
